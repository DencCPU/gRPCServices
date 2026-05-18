package apprunner

import (
	"context"
	"fmt"

	spot "github.com/DencCPU/gRPCServices/Protobuf/gen/spot_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryspotservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_spot_service"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"github.com/DencCPU/gRPCServices/Shared/opentelemetry"
	spotconfig "github.com/DencCPU/gRPCServices/SpotInstrumentService/config"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/adapters/memory"
	redisadapter "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/adapters/redis"
	spothandlers "github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/controllers/grpc_handlers"
	"github.com/DencCPU/gRPCServices/SpotInstrumentService/internal/usecase"
	grpcserver "github.com/DencCPU/gRPCServices/SpotInstrumentService/pkg/spotserver"
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
					return nil, fmt.Errorf("logger initialition error:%w", err)
				}
				return logger, nil
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

// Add in-memory storage
func StorageModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger) (*memory.Storage, error) {

				storage, err := memory.NewStorage(logger)
				if err != nil {
					return nil, err
				}
				return storage, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, storage *memory.Storage) {
				marketsPath := "./SpotInstrumentService/config/market/markets.txt"
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						err := storage.AddMarkets(marketsPath)
						if err != nil {
							return err
						}
						return nil
					},
				})
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, storage *memory.Storage, cfg *spotconfig.Config) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {

							ctx, cancel := context.WithTimeout(context.Background(), cfg.Storage.Timeout)
							defer cancel()
							msg := storage.AccessControl(ctx)
							logger.Info(msg)
						}()
						return nil
					},
				})
			},
		),
	)
}

// Add config module
func NewConfigModul() fx.Option {
	return fx.Provide(
		func() *config.ConfigLoader {
			loader := config.NewConfigLoader(
				entryspotservice.GlobalPathToEnv,
				entryspotservice.EnvFile,
				entryspotservice.ConfigType,
				entryspotservice.PathToLocalEnv,
				entryspotservice.PathToConfig,
			)
			return loader
		},
		func(loader *config.ConfigLoader) (*spotconfig.Config, error) {
			config, err := config.NewConfig[spotconfig.Config](loader)
			if err != nil {
				return nil, fmt.Errorf("error getting config:%w", err)
			}
			return config, nil
		},
	)
}

// Add redis
func RedisModule() fx.Option {
	return fx.Provide(
		func(config *spotconfig.Config, logger *zap.Logger) (*redisadapter.RedisDB, error) {
			rctx, cancel := context.WithTimeout(context.Background(), config.Redis.Timeout)
			defer cancel()
			redis, err := redisadapter.NewRedis(rctx, config.Redis)
			if err != nil {
				logger.Error("redis initialization error:",
					zap.Error(err))
				return nil, err
			}
			return redis, nil
		},
	)
}

// Add tracer
func TracingModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, config *spotconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{},
				))
				trace, err := opentelemetry.NewTracerProviderGrpc(
					context.Background(),
					"Spot/service",
					config.OtelCollector.Host,
					config.OtelCollector.Port,
					config.OtelCollector.TracePercentage,
				)

				if err != nil {
					logger.Error("tracer initialization error:",
						zap.Error(err))
					return nil, nil, err
				}
				tracer := trace.Tracer("spot/service")
				return trace, tracer, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, trace *sdktrace.TracerProvider) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						return trace.Shutdown(ctx)
					},
				})
			},
		),
	)
}

// Add metrics
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *spotconfig.Config) (*sdkmetric.MeterProvider, error) {

				provider, err := opentelemetry.NewMetricProviderGrpc(
					context.Background(),
					"Spot/service",
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

// OtelLogger
func OtelLogger() fx.Option {
	return fx.Options(fx.Provide(
		func(cfg *spotconfig.Config) (*zap.Logger, error) {
			logger, err := logger.NewLogger()
			if err != nil {
				return nil, fmt.Errorf("zap logger initialition error:%w", err)
			}

			loggerProvider, err := opentelemetry.NewLoggerProviderGrpc(
				context.Background(),
				"Spot/service",
				cfg.OtelCollector.Host,
				cfg.OtelCollector.Port)

			if err != nil {
				logger.Error("logger provider initialization error",
					zap.Error(err),
				)
				return nil, err

			}
			otlpCore := otelzap.NewCore("Spot/service", otelzap.WithLoggerProvider(loggerProvider))
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
	return fx.Provide(
		func(storage *memory.Storage, logger *zap.Logger, trace trace.Tracer) *usecase.SpotService {
			return usecase.NewSpotInstrument(storage, logger, trace)
		},
	)
}

// Add handlers
func HandlersModule() fx.Option {
	return fx.Provide(
		func(service *usecase.SpotService) *spothandlers.Handlers {
			return spothandlers.NewHandlers(service)
		},
	)
}

// Add gRPC-server
func GrpcModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(redis *redisadapter.RedisDB, config *spotconfig.Config, logger *zap.Logger) (*grpcserver.Server, error) {
				grpcServer, err := grpcserver.New(redis, config.Server, logger)
				if err != nil {
					logger.Error("grpc server initialization error:",
						zap.Error(err),
					)
					return nil, err
				}
				return grpcServer, nil
			},
		),

		fx.Invoke(
			func(lc fx.Lifecycle, server *grpcserver.Server, handlers *spothandlers.Handlers, logger *zap.Logger, config *spotconfig.Config) {

				spot.RegisterSpotInstrumentServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							info := fmt.Sprintf("Server start on port:%s", config.Server.Port)
							logger.Info(info)
							err := server.Serve(server.Listener)
							if err != nil {
								logger.Error("error working server:",
									zap.Error(err))
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						done := make(chan struct{})
						go func() {
							logger.Info("stopping the server...")
							server.GracefulStop()
							close(done)
						}()
						select {
						case <-done:
							logger.Info("Server stopped")
							return nil
						case <-ctx.Done():
							logger.Info("Stopping the server by timeout")
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
		StorageModul(),
		NewConfigModul(),
		RedisModule(),
		TracingModule(),
		MetricModul(),
		OtelLogger(),
		ServiceModule(),
		HandlersModule(),
		GrpcModule(),
	)

	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("dependency graph initialization error:%w", err)
	}

	return app, nil
}
