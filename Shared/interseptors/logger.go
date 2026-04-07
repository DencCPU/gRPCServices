package interseptors

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Логирование запросов
func LoggerInterseptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {

		//Изъятие ID-запроса из контекста
		requestID, _ := ctx.Value(requestIDKey).(string)

		//Логирование входящего запроса
		logger.Info("Service request",
			zap.String("Request number", requestID),
			zap.String("Called method", info.FullMethod),
		)

		//Вызов следующего обработчика
		resp, err := handler(ctx, req)
		return resp, err
	}
}
