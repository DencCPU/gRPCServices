package apprunner

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	apiconfig "github.com/DencCPU/gRPCServices/APIGetway/config"
	orderclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/order_client"
	spotclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/spot_client"
	userclient "github.com/DencCPU/gRPCServices/APIGetway/internal/adapters/user_client"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/controller/gin"
	"github.com/DencCPU/gRPCServices/APIGetway/internal/usecase"
	"github.com/DencCPU/gRPCServices/Shared/breaker"
	"github.com/DencCPU/gRPCServices/Shared/config"
	entryapiservice "github.com/DencCPU/gRPCServices/Shared/enter_points/entry_api_service"
	"github.com/DencCPU/gRPCServices/Shared/logger"
	"github.com/DencCPU/gRPCServices/Shared/opentelemetry"
	"github.com/sony/gobreaker"
	"go.opentelemetry.io/contrib/bridges/otelzap"
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
					entryapiservice.GlobalPathToEnv,
					entryapiservice.EnvFile,
					entryapiservice.ConfigType,
					entryapiservice.PathToLocalEnv,
					entryapiservice.PathToConfig,
				)
			},
			func(loader *config.ConfigLoader) (*apiconfig.Config, error) {
				cfg, err := config.NewConfig[apiconfig.Config](loader)
				if err != nil {
					return nil, fmt.Errorf("error getting new config:%w", err)
				}
				return cfg, nil
			},
		),
	)
}

// Add Breaker module
func BreakerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *apiconfig.Config, logger *zap.Logger) *gobreaker.CircuitBreaker {
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
			func(logger *zap.Logger, breaker *gobreaker.CircuitBreaker) (*spotclient.Client, error) {
				spotClient, err := spotclient.NewClient(breaker)
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

// Add Order client
func OrderClientModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, breaker *gobreaker.CircuitBreaker) (*orderclient.Client, error) {
				orderClient, err := orderclient.NewClient(breaker)
				if err != nil {
					logger.Error("error initialization order service client:",
						zap.Error(err),
					)
					return nil, err
				}
				return orderClient, nil
			},
		),
	)
}

// Add User client
func UserClientModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(cfg *apiconfig.Config, logger *zap.Logger, breaker *gobreaker.CircuitBreaker) (*userclient.Client, error) {
				userClient, err := userclient.NewClient(breaker)
				if err != nil {
					logger.Error("error initialization user service client:",
						zap.Error(err),
					)
					return nil, err
				}
				return userClient, nil
			},
		),
	)
}

// Add tracer
func TracerModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *apiconfig.Config) (*sdktrace.TracerProvider, trace.Tracer, error) {
				ctx := context.Background()

				tracerProvider, err := opentelemetry.NewTracerProviderHttp(
					ctx,
					"API/getway",
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
				trace := tracerProvider.Tracer("API/getway")
				return tracerProvider, trace, nil
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

// Add metric
func MetricModul() fx.Option {
	return fx.Options(
		fx.Provide(
			func(logger *zap.Logger, cfg *apiconfig.Config) (*sdkmetric.MeterProvider, error) {

				provider, err := opentelemetry.NewMetricProviderHttp(
					context.Background(),
					"API/getway",
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
		func(cfg *apiconfig.Config) (*zap.Logger, error) {
			logger, err := logger.NewLogger()
			if err != nil {
				return nil, fmt.Errorf("zap logger initialition error:%w", err)
			}

			loggerProvider, err := opentelemetry.NewLoggerProviderHttp(
				context.Background(),
				"API/getway",
				cfg.OtelCollector.Host,
				cfg.OtelCollector.Port)

			if err != nil {
				logger.Error("logger provider initialization error",
					zap.Error(err),
				)
				return nil, err

			}
			otlpCore := otelzap.NewCore("API/getway", otelzap.WithLoggerProvider(loggerProvider))
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
			func(userClient *userclient.Client, orderClient *orderclient.Client, spotClient *spotclient.Client, logger *zap.Logger, tracer trace.Tracer) *usecase.Service {
				return usecase.NewService(userClient, orderClient, spotClient, logger, tracer)
			},
		),
	)
}

// Add handlers
func HandlersModule() fx.Option {
	return fx.Options(
		fx.Provide(
			func(service *usecase.Service) gin.GinAPI {
				return gin.NewGinAPI(service)
			},
		),
	)
}

// Add server working
func ServerModule() fx.Option {
	return fx.Options(
		fx.Invoke(
			func(lc fx.Lifecycle, api gin.GinAPI, logger *zap.Logger, cfg *apiconfig.Config) {
				host := cfg.Server.Host
				port := strconv.Itoa(cfg.Server.Port)
				srv := &http.Server{
					Addr:    host + ":" + port,
					Handler: api.Router(),
				}
				lc.Append(fx.Hook{
					OnStart: func(ctx context.Context) error {
						info := fmt.Sprintf("Server is running on port:%d", cfg.Server.Port)
						logger.Info(info)

						go func() {
							err := srv.ListenAndServe()
							if err != nil {
								logger.Error("server error:",
									zap.Error(err),
								)
							}
						}()

						return nil
					},
					OnStop: func(ctx context.Context) error {

						shutdownCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
						defer cancel()

						return srv.Shutdown(shutdownCtx)
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
		BreakerModule(),
		SpotClientModul(),
		OrderClientModul(),
		UserClientModule(),
		TracerModule(),
		MetricModul(),
		ServiceModule(),
		HandlersModule(),
		ServerModule(),
	)
	if err := app.Err(); err != nil {
		return nil, fmt.Errorf("dependency graph initialization error:%w", err)
	}
	return app, nil
}
