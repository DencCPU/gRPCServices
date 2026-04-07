package apprunner

import (
	"context"
	"fmt"
	"net/http"

	user "github.com/DencCPU/gRPCServices/Protobuf/gen/user_service"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryuserservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_user_service"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"github.com/DencCPU/gRPCServices/Shared/opentelemetry"
	userconfig "github.com/DencCPU/gRPCServices/UserService/config"
	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/jwt"
	"github.com/DencCPU/gRPCServices/UserService/internal/adapters/postgres"
	userhandlers "github.com/DencCPU/gRPCServices/UserService/internal/controllers/grpc_handlers"
	"github.com/DencCPU/gRPCServices/UserService/internal/usecase"
	"github.com/DencCPU/gRPCServices/UserService/pkg/userserver"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Logger
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

// Config
func ConfigModul() fx.Option {
	return fx.Provide(
		func() *config.ConfigLoader {
			loader := config.NewConfigLoader(
				entryuserservice.GlobalPathToEnv,
				entryuserservice.EnvFile,
				entryuserservice.ConfigType,
				entryuserservice.PathToLocalEnv,
				entryuserservice.PathToConfig,
			)
			return loader
		},
		func(loader *config.ConfigLoader, logger *zap.Logger) (*userconfig.Config, error) {
			config, err := config.NewConfig[userconfig.Config](loader)
			if err != nil {
				logger.Error("error getting config",
					zap.Error(err),
				)
				return nil, err
			}
			return config, nil
		},
	)
}

// Postgres
func PostgresModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *userconfig.Config) (*postgres.PostgresDB, error) {
				storage, err := postgres.NewDB(context.Background(), logger, cfg.Postgres)
				if err != nil {
					logger.Error("storage initialization error:",
						zap.Error(err),
					)
					return nil, err
				}
				return storage, nil
			},
		),
	)
}

// Tracers
func TracingModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, config *userconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{},
				))
				trace, err := opentelemetry.NewTrace(context.Background(), "userService", config.Jaeger.Host, config.Jaeger.Port)
				if err != nil {
					logger.Error("tracer initialization error:",
						zap.Error(err))
					return nil, nil, err
				}
				tracer := trace.Tracer("UserService")
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

// Metrics
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger) (*sdkmetric.MeterProvider, error) {
				provider, err := opentelemetry.NewMetricPrometeus(context.Background(), "UserService")
				if err != nil {
					logger.Error("error initialization metric:",
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
							logger.Info("The server is running on port 9466...")
							err := http.ListenAndServe(":9466", nil)
							if err != nil {
								logger.Error("server error for metric:",
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

// JWT
func JwtModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *userconfig.Config) *jwt.JWT {
				return jwt.NewJWT(cfg.JWT)
			},
		),
	)
}

// Service
func ServiceModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(storage *postgres.PostgresDB, logger *zap.Logger, jwt *jwt.JWT, tracer trace.Tracer) *usecase.Service {
				return usecase.NewService(storage, logger, jwt, tracer)
			},
		),
	)
}

// Handlers
func HandlersModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(service *usecase.Service) *userhandlers.Handler {
				return userhandlers.NewHandler(service)
			},
		),
	)
}

// GRPC server
func GrpcModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *userconfig.Config, logger *zap.Logger) (*userserver.Server, error) {
				server, err := userserver.NewServer(cfg.Server, logger)
				if err != nil {
					logger.Error("error creating server:",
						zap.Error(err),
					)
					return nil, err
				}
				return server, nil
			},
		),
		fx.Invoke(
			func(lc fx.Lifecycle, server *userserver.Server, handlers *userhandlers.Handler, logger *zap.Logger) {
				user.RegisterUserServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							logger.Info("Server start on port 8083...")
							err := server.Serve(server.Listener)
							if err != nil {
								logger.Error("error working server:",
									zap.Error(err),
								)
							}
						}()
						return nil
					},
					OnStop: func(ctx context.Context) error {
						done := make(chan struct{})
						go func() {
							logger.Info("Stoping the server...")
							server.GracefulStop()
							close(done)
						}()
						select {
						case <-done:
							logger.Info("Server stop")
							return nil
						case <-ctx.Done():
							logger.Warn("stoping the server by timeout")
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
		TracingModule(),
		MetricModul(),
		JwtModule(),
		ServiceModule(),
		HandlersModule(),
		GrpcModule(),
	)
	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("dependency graph initialization error:%w", err)
	}
	return app, nil
}
