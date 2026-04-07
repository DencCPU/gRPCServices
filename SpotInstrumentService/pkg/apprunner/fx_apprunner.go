package apprunner

import (
	"context"
	"fmt"
	"net/http"
	"time"

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
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Подкючение логгера
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

// Подключение In-memory хранилища
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
			func(lc fx.Lifecycle, logger *zap.Logger, storage *memory.Storage) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							//Создание контекста для времени управления состояниями рынков
							ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
							defer cancel()
							msg := storage.AccessControl(ctx) //Запуск управления рынками
							logger.Info(msg)
						}()
						return nil
					},
				})
			},
		),
	)
}

// Получение конфига
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
		func(loader *config.ConfigLoader, logger *zap.Logger) (*spotconfig.Config, error) {
			config, err := config.NewConfig[spotconfig.Config](loader)
			if err != nil {
				logger.Error("error getting config:",
					zap.Error(err),
				)
				return nil, err
			}
			return config, nil
		},
	)
}

// Подключение redis-кэш
func RedisModule() fx.Option {
	return fx.Provide(
		func(config *spotconfig.Config, logger *zap.Logger) (*redisadapter.RedisDB, error) {
			rctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
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

// Подключение трейсера
func TracingModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, config *spotconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{},
				))
				trace, err := opentelemetry.NewTrace(context.Background(), "spotService", config.Jaeger.Host, config.Jaeger.Port)
				if err != nil {
					logger.Error("tracer initialization error:",
						zap.Error(err))
					return nil, nil, err
				}
				tracer := trace.Tracer("SpotService")
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

// Подключение метрик
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger) (*sdkmetric.MeterProvider, error) {
				provider, err := opentelemetry.NewMetricPrometeus(context.Background(), "SpotService")
				if err != nil {
					logger.Error("metrics initialization error:",
						zap.Error(err))
					return nil, err
				}

				return provider, err
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, provider *sdkmetric.MeterProvider) {
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							http.Handle("/metrics", promhttp.Handler())
							logger.Info("The metrics server has been launched on port 9464...")
							err := http.ListenAndServe(":9464", nil)
							if err != nil {
								logger.Error("server error for metrics:",
									zap.Error(err),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {

						return provider.Shutdown(ctx)
					},
				})
			},
		),
	)
}

// Подключение сервиса обработки
func ServiceModule() fx.Option {
	return fx.Provide(
		func(storage *memory.Storage, logger *zap.Logger, trace trace.Tracer) *usecase.SpotService {
			return usecase.NewSpotInstrument(storage, logger, trace)
		},
	)
}

// Подключение обработчиков
func HandlersModule() fx.Option {
	return fx.Provide(
		func(service *usecase.SpotService) *spothandlers.Handlers {
			return spothandlers.NewHandlers(service)
		},
	)
}

// Создание gRPC сервера
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
			func(lc fx.Lifecycle, server *grpcserver.Server, handlers *spothandlers.Handlers, logger *zap.Logger) {
				//Регистрация grpc методов в сервере
				spot.RegisterSpotInstrumentServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							logger.Info("Server start on port 8080...")
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
		LoggerModul(),
		StorageModul(),
		NewConfigModul(),
		RedisModule(),
		TracingModule(),
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
