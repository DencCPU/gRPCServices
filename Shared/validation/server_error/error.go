package servererror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Функция проверяет вид ошибки(клиентская или серверная)
func ServerError(err error) bool {
	if err == nil {
		return true
	}
	status, ok := status.FromError(err)
	if !ok {
		return false //Неизыестна ошибка
	}
	switch status.Code() {
	case codes.Unavailable, // сервис недоступен
		codes.DeadlineExceeded,  // таймаут
		codes.Internal,          // внутренняя ошибка
		codes.Unknown,           // неизвестная ошибка
		codes.ResourceExhausted: // перегрузка
		return false
	default:
		return true
	}
}
