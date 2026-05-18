package interceptors

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func XRequestID(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		values := md.Get(string(requestID))
		var id string
		if len(values) != 0 {
			id = values[0]
		} else {
			id = uuid.New().String()
		}
		ctx = context.WithValue(ctx, string(requestID), id)
	} else {
		ctx = context.WithValue(ctx, string(requestID), uuid.New().String())
	}
	return handler(ctx, req)
}
