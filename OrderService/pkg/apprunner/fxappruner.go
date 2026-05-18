package apprunner

import (
	"context"
	"fmt"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	"github.com/DencCPU/gRPCServices/OrderService/internal/adapters/notify"
	"github.com/DencCPU/gRPCServices/OrderService/internal/adapters/postgres"
	spotservice "github.com/DencCPU/gRPCServices/OrderService/internal/adapters/spot_service"
	orderhandlers "github.com/DencCPU/gRPCServices/OrderService/internal/controllers/grpc_handlers"
	ordererrors "github.com/DencCPU/gRPCServices/OrderService/internal/domain/error"
	"github.com/DencCPU/gRPCServices/OrderService/internal/usecase"
	"github.com/DencCPU/gRPCServices/OrderService/pkg/orderserver"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/DencCPU/gRPCServices/Shared/breaker"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryorderservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_order_service"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"github.com/DencCPU/gRPCServices/Shared/opentelemetry"
	"github.com/sony/gobreaker"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Add logger
func LoggerModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func() (*zap.Logger, error) {
				logger, err := logger.NewLogger()
				if err != nil {
					return nil, fmt.Errorf("logger initialization error:%w", err)
				}
				return logger, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger) {
				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							_ = logger.Sync()
							return nil
						},
					},
				)
			},
		),
	)
}

// Get new config
func ConfigModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *config.ConfigLoader {
				return config.NewConfigLoader(
					entryorderservice.GlobalPathToEnv,
					entryorderservice.EnvFile,
					entryorderservice.ConfigType,
					entryorderservice.PathToLocalEnv,
					entryorderservice.PathToConfig,
				)
			},
			func(loader *config.ConfigLoader) (*orderconfig.Config, error) {
				cfg, err := config.NewConfig[orderconfig.Config](loader)
				if err != nil {
					return nil, fmt.Errorf("error getting new config:%w", err)
				}
				return cfg, nil
			},
		),
	)
}

// Add notify service
func NotifyModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config) *notify.StatusStorage {
				return notify.NewStatStorage(cfg.Notify.TickerInterval)
			},
		),
	)
}

// Add Postgres
func PostgresModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *orderconfig.Config, notify *notify.StatusStorage) (*postgres.PostgresDB, error) {
				storage, err := postgres.NewDB(context.Background(), cfg.Postgres, notify)
				if err != nil {
					logger.Error("error initialization postgres database:",
						zap.Error(err),
					)
					return nil, err
				}
				return storage, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, dataBase *postgres.PostgresDB, logger *zap.Logger) {
				errChan := make(chan ordererrors.ErrStruct, 10)
				lc.Append(fx.Hook{

					OnStart: func(ctx context.Context) error {
						go dataBase.ControlOrder(errChan)
						dataBase.CheckCacheTTL()

						go func() {
							for el := range errChan {
								logger.Error("control order error:",
									zap.String("OrderId:", el.OrderId),
									zap.Error(el.Err),
								)
							}
						}()
						return nil
					},

					OnStop: func(ctx context.Context) error {
						dataBase.StopControlOrder()
						close(errChan)
						return nil
					},
				})
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, dataBase *postgres.PostgresDB) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						dataBase.GetPgxPool().Close()
						return nil
					},
				})
			},
		),
	)
}

// Add Breaker module
func BreakerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config, logger *zap.Logger) *gobreaker.CircuitBreaker {
				params := breaker.Params{
					Name:           cfg.BreakerSetting.Name,
					MaxRequest:     cfg.BreakerSetting.MaxRequests,
					Interval:       cfg.BreakerSetting.Interval,
					Timeout:        cfg.BreakerSetting.Timeout,
					MaxFailRequest: cfg.BreakerSetting.MaxFailRequest,
				}
				breaker := breaker.NewBreaker(logger, params)
				return breaker
			},
		),
	)
}

// Add SpotService client
func SpotClientModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config, logger *zap.Logger, breaker *gobreaker.CircuitBreaker) (*spotservice.Client, error) {
				spotClient, err := spotservice.NewClient(cfg.BreakerSetting, logger, breaker)
				if err != nil {
					logger.Error("error initialization spot service client:",
						zap.Error(err),
					)
					return nil, err
				}
				return spotClient, nil
			},
		),
	)
}

// Add tracer
func JaegerTracerModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *orderconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				//tracer initialization
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{}))

				trace, err := opentelemetry.NewTracerProviderGrpc(
					context.Background(),
					"Order/service",
					cfg.OtelCollector.Host,
					cfg.OtelCollector.Port,
					cfg.OtelCollector.TracePercentage,
				)

				if err != nil {
					logger.Error("tracer startup error:",
						zap.Error(err),
					)
					return nil, nil, err
				}
				tracer := trace.Tracer("Order/service")
				return trace, tracer, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, trace *sdktrace.TracerProvider, logger *zap.Logger) {
				lc.Append(
					fx.Hook{
						OnStop: func(ctx context.Context) error {
							trace.Shutdown(ctx)
							return nil
						},
					},
				)
			},
		),
	)
}

// Add metrics
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *orderconfig.Config) (*sdkmetric.MeterProvider, error) {

				provider, err := opentelemetry.NewMetricProviderGrpc(
					context.Background(),
					"Order/service",
					cfg.OtelCollector.Host,
					cfg.OtelCollector.Port,
					cfg.OtelCollector.MetricInterval,
				)
				if err != nil {
					logger.Error("error initialization metric:",
						zap.Error(err))
					return nil, err
				}
				return provider, err
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, metric *sdkmetric.MeterProvider) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						metric.Shutdown(ctx)
						logger.Info("Metric stoped")
						return nil
					},
				})
			},
		),
	)
}

// Add otelLogger
func OtelLogger() fx.Option {
	return fx.Options(fx.Provide(
		func(cfg *orderconfig.Config) (*zap.Logger, error) {
			logger, err := logger.NewLogger()
			if err != nil {
				return nil, fmt.Errorf("zap logger initialition error:%w", err)
			}

			loggerProvider, err := opentelemetry.NewLoggerProviderGrpc(
				context.Background(),
				"Order/service",
				cfg.OtelCollector.Host,
				cfg.OtelCollector.Port)

			if err != nil {
				logger.Error("logger provider initialization error",
					zap.Error(err),
				)
				return nil, err

			}
			otlpCore := otelzap.NewCore("Order/service", otelzap.WithLoggerProvider(loggerProvider))
			teeCore := zapcore.NewTee(logger.Core(), otlpCore)
			otlpLogger := zap.New(teeCore, zap.WithCaller(true))

			return otlpLogger, nil
		},
	),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						_ = logger.Sync()
						return nil
					},
				})
			},
		),
	)
}

// Add processing service
func ServiceModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(storage *postgres.PostgresDB, spotClient *spotservice.Client, notify *notify.StatusStorage, logger *zap.Logger, trace trace.Tracer) *usecase.OrderService {
				return usecase.NewOrderServ(storage, spotClient, notify, logger, trace)
			},
		),
	)
}

// Add handlers
func HandlersModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(service *usecase.OrderService) *orderhandlers.Handlers {
				return orderhandlers.NewHandlers(service)
			},
		),
	)
}

// Add gRCP-server
func GrpcModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config, logger *zap.Logger) (*orderserver.Server, error) {
				server, err := orderserver.New(cfg.Server, logger)
				if err != nil {
					logger.Error("grpc-server initialization error:",
						zap.Error(err),
					)
					return nil, err
				}
				return server, nil
			},
		),

		//Start server
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, server *orderserver.Server, handlers *orderhandlers.Handlers, cfg *orderconfig.Config) {
				order.RegisterOrderServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							info := fmt.Sprintf("The server is running on port:%d", cfg.Server.Port)
							logger.Info(info)
							err := server.Serve(server.Listener)
							if err != nil {
								logger.Error("grpc-server error:",
									zap.Error(err),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Stoping the server")
						done := make(chan struct{})
						go func() {
							server.GracefulStop()
							close(done)
						}()
						select {
						case <-done:
							logger.Info("Server is stopped")
							return nil
						case <-ctx.Done():
							logger.Info("The server stopped due to context timeout.")
							server.Stop()
							return ctx.Err()
						}
					},
				})
			},
		),
	)
}
func FxAppRunner() (*fx.App, error) {
	app := fx.New(
		// LoggerModul(),
		OtelLogger(),
		ConfigModul(),
		PostgresModul(),
		BreakerModule(),
		SpotClientModul(),
		NotifyModul(),
		JaegerTracerModul(),
		MetricModul(),
		ServiceModule(),
		HandlersModule(),
		GrpcModule(),
	)
	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("dependency graph initialization error:%w", err)
	}
	return app, nil
}
