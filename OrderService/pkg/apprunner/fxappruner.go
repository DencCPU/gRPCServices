package apprunner

import (
	"context"
	"fmt"
	"net/http"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	"github.com/DencCPU/gRPCServices/OrderService/internal/adapters/notify"
	"github.com/DencCPU/gRPCServices/OrderService/internal/adapters/postgres"
	spotservice "github.com/DencCPU/gRPCServices/OrderService/internal/adapters/spot_service"
	orderhandlers "github.com/DencCPU/gRPCServices/OrderService/internal/controllers/grpc_handlers"
	"github.com/DencCPU/gRPCServices/OrderService/internal/usecase"
	"github.com/DencCPU/gRPCServices/OrderService/pkg/orderserver"
	order "github.com/DencCPU/gRPCServices/Protobuf/gen/order_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryorderservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_order_service"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"github.com/DencCPU/gRPCServices/Shared/opentelemetry"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Подключение логера
func LoggerModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func() (*zap.Logger, error) {
				logger, err := logger.NewLogger()
				if err != nil {
					return nil, fmt.Errorf("ошибка инициализации логгера:%w", err)
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

// Получение нового конфига
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
			func(loader *config.ConfigLoader, logger *zap.Logger) (*orderconfig.Config, error) {
				cfg, err := config.NewConfig[orderconfig.Config](loader)
				if err != nil {
					logger.Error("ошибка получения нового конфига:",
						zap.Error(err),
					)
					return nil, err
				}
				return cfg, nil
			},
		),
	)
}

// Подключение к Postgres-хранилищу
func PostgresModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *orderconfig.Config) (*postgres.PostgresDB, error) {
				storage, err := postgres.NewDB(context.Background(), logger, cfg.Postgres)
				if err != nil {
					logger.Error("ошибка инициализации хранилища:",
						zap.Error(err),
					)
					return nil, err
				}
				return storage, nil
			},
		),
	)
}

// Подключение клиента SpotService
func SpotClientModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config, logger *zap.Logger) (*spotservice.Client, error) {
				spotClient, err := spotservice.NewClient(cfg.BreakerSetting, logger)
				if err != nil {
					logger.Error("ошибка инициализации хранилища:",
						zap.Error(err),
					)
					return nil, err
				}
				return spotClient, nil
			},
		),
	)
}

// Подключение сервиса нотификаций
func NotifyModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func() *notify.StatusStorage {
				return notify.NewStatStorage()
			},
		),
	)
}

// Подключение трейсера
func TracerModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *orderconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				//Инициализация трейсера
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{}))

				trace, err := opentelemetry.NewTrace(context.Background(), "OrderSrevice", cfg.Jaeger.Host, cfg.Jaeger.Port)
				if err != nil {
					logger.Error("ошибка запуска трейсера:",
						zap.Error(err),
					)
					return nil, nil, err
				}
				tracer := trace.Tracer("OrderService")
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

// Подключение метрик
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger) (*sdkmetric.MeterProvider, error) {
				provider, err := opentelemetry.NewMetricPrometeus(context.Background(), "OrderService")
				if err != nil {
					logger.Error("ошибка инициализации метрик:",
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
							logger.Info("Сервер для метрик запущен на порту 9465...")
							err := http.ListenAndServe(":9465", nil)
							if err != nil {
								logger.Error("Ошибка работы сервера с метриками:",
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
	return fx.Options(
		fx.Provide(
			func(storage *postgres.PostgresDB, spotClient *spotservice.Client, notify *notify.StatusStorage, logger *zap.Logger, trace trace.Tracer) *usecase.OrderService {
				return usecase.NewOrderServ(storage, spotClient, notify, logger, trace)
			},
		),
	)
}

// Подключение обработчиков
func HandlersModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(service *usecase.OrderService) *orderhandlers.Handlers {
				return orderhandlers.NewHandlers(service)
			},
		),
	)
}

// Создание и запуск gRPC сервера
func GrpcModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *orderconfig.Config, logger *zap.Logger) (*orderserver.Server, error) {
				server, err := orderserver.New(cfg.Server, logger)
				if err != nil {
					logger.Error("ошибка инициализации grpc-сервера:",
						zap.Error(err),
					)
					return nil, err
				}
				return server, nil
			},
		),

		//Запуск сервера
		fx.Invoke(
			func(lc fx.Lifecycle, logger *zap.Logger, server *orderserver.Server, handlers *orderhandlers.Handlers) {
				order.RegisterOrderServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							logger.Info("Сервер работает на порту 8081...")
							err := server.Serve(server.Listener)
							if err != nil {
								logger.Error("ошибка инициализации grpc-сервера:",
									zap.Error(err),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						logger.Info("Остановка работы сервера")
						done := make(chan struct{})
						go func() {
							server.GracefulStop()
							close(done)
						}()
						select {
						case <-done:
							logger.Info("Сервер остановлен")
							return nil
						case <-ctx.Done():
							logger.Info("Сервер остановлен по тайамауту контекста")
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
		ConfigModul(),
		PostgresModul(),
		SpotClientModul(),
		NotifyModul(),
		TracerModul(),
		MetricModul(),
		ServiceModule(),
		HandlersModule(),
		GrpcModule(),
	)
	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("ошибка инициализации графа зависимостей:%w", err)
	}
	return app, nil
}
