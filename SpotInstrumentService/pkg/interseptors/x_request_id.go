package interseptors

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
	xctx := context.Background()
	if ok {
		values := md.Get(requestIDKey)
		var id string
		if len(values) != 0 {
			id = values[0]
		} else {
			id = uuid.New().String()
		}
		xctx = context.WithValue(ctx, requestIDKey, id)
	} else {
		xctx = context.WithValue(ctx, requestIDKey, uuid.New().String())
	}
	return handler(xctx, req)
}
