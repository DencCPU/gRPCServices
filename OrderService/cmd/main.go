package main

import (
	"Academy/gRPCServices/OrderService/internal/adapters/notify"
	"Academy/gRPCServices/OrderService/internal/adapters/postgres"
	spotservice "Academy/gRPCServices/OrderService/internal/adapters/spot_service"
	orderhandlers "Academy/gRPCServices/OrderService/internal/controllers/grpc_handlers"
	"Academy/gRPCServices/OrderService/internal/usecase"
	"Academy/gRPCServices/OrderService/pkg/orderserver"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
	"Academy/gRPCServices/SpotInstrumentService/pkg/opentelimetry"

	"context"
	"fmt"
	"log"
)

func main() {
	// storage := memory.NewStorage() //Инициализация хранилища in-memory

	ctx := context.Background()
	storage, err := postgres.NewDB(ctx) //Инициализаця хранилища postgres
	if err != nil {
		log.Fatal(err)
	}

	spotClient, err := spotservice.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	notify := notify.NewStatStorage()
	service := usecase.NewOrderServ(storage, spotClient, notify)

	handlers := orderhandlers.NewHandlers(service)

	tp, err := opentelimetry.NewTrace(ctx, "OrderSrevice")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(ctx)

	grpcServer, err := orderserver.New()
	if err != nil {
		log.Fatal(err)
	}

	orderAPI.RegisterOrderServiceServer(grpcServer, handlers)

	fmt.Println("Сервер работает на порту 8081...")
	err = grpcServer.Serve(grpcServer.Listener)
	if err != nil {
		log.Fatal("Ошибка работы сервера:", err)
	}
}
