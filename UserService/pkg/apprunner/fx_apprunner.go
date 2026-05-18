package apprunner

import (
	"context"
	"fmt"

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
		func(loader *config.ConfigLoader) (*userconfig.Config, error) {
			config, err := config.NewConfig[userconfig.Config](loader)
			if err != nil {
				return nil, fmt.Errorf("error getting config:%w", err)
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
		fx.Invoke(
			func(lc fx.Lifecycle, db *postgres.PostgresDB) {
				lc.Append(fx.Hook{
					OnStop: func(ctx context.Context) error {
						db.Close()
						return nil
					},
				})
			},
		),
	)
}

// Tracers
func TracingModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, config *userconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				fmt.Println("TRACING:", config.OtelCollector.Host+":"+config.OtelCollector.Port, config.OtelCollector.MetricInterval)
				otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
					propagation.TraceContext{},
					propagation.Baggage{},
				))
				trace, err := opentelemetry.NewTracerProviderGrpc(context.Background(), "userService", config.OtelCollector.Host, config.OtelCollector.Port, config.OtelCollector.TracePercentage)
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
			func(logger *zap.Logger, cfg *userconfig.Config) (*sdkmetric.MeterProvider, error) {

				provider, err := opentelemetry.NewMetricProviderGrpc(
					context.Background(),
					"User/service",
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
		func(cfg *userconfig.Config) (*zap.Logger, error) {
			logger, err := logger.NewLogger()
			if err != nil {
				return nil, fmt.Errorf("zap logger initialition error:%w", err)
			}

			loggerProvider, err := opentelemetry.NewLoggerProviderGrpc(
				context.Background(),
				"User/service",
				cfg.OtelCollector.Host,
				cfg.OtelCollector.Port)

			if err != nil {
				logger.Error("logger provider initialization error",
					zap.Error(err),
				)
				return nil, err

			}
			otlpCore := otelzap.NewCore("User/service", otelzap.WithLoggerProvider(loggerProvider))
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
			func(lc fx.Lifecycle, server *userserver.Server, handlers *userhandlers.Handler, logger *zap.Logger, cfg *userconfig.Config) {
				user.RegisterUserServiceServer(server, handlers)

				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						go func() {
							info := fmt.Sprintf("Server start on port:%s", cfg.Server.Port)
							logger.Info(info)
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
		// LoggerModul(),
		OtelLogger(),
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
