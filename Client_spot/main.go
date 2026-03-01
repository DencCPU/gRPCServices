package main

import (
	spotAPI "Academy/gRPCServices/Protobuf/gen/spot"
	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

func main() {
	// Контекст с таймаутом для подключения
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Подключение к серверу gRPC с OTel stats handler
	conn, err := grpc.NewClient(
		"localhost:8080",
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()), // <-- OTel клиентский handler
	)
	if err != nil {
		log.Fatal("ошибка подключения:", err)
	}
	defer conn.Close()

	// Создаём клиент
	client := spotAPI.NewSpotInstrumentServiceClient(conn)

	// Делаем запрос
	resp, err := client.ViewMarket(ctx, &spotAPI.ViewReq{})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Доступные рынки:", resp.EnableMarkets)
}
