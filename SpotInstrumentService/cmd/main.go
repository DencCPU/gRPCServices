package main

import (
	spotAPI "Academy/gRPCServices/Protobuf/gen/spot"
	"Academy/gRPCServices/SpotInstrumentService/internal/adapters/memory"
	redisadapter "Academy/gRPCServices/SpotInstrumentService/internal/adapters/redis"
	spothandlers "Academy/gRPCServices/SpotInstrumentService/internal/controllers/grpc_handlers"

	"Academy/gRPCServices/SpotInstrumentService/pkg/opentelimetry"
	grpcserver "Academy/gRPCServices/SpotInstrumentService/pkg/spotserver"
	"context"
	"time"

	"Academy/gRPCServices/SpotInstrumentService/internal/usecase"

	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup

	//Инициализация хранилища
	storage, err := memory.NewStorage()
	if err != nil {
		log.Fatal("ошибка инициализации хранилища:", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		//Создание контекста для времени управления состояниями рынков
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		l := storage.AccessControl(ctx) //Запуск управления рынками
		log.Println("управление рынками:", l)
	}()

	//Инициализая нового сервиса
	service := usecase.NewSpotInstrument(storage)

	//Инициализация обработчиков
	handlers := spothandlers.NewHandlers(service)

	//Создание контекста для redis
	rctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//Инициализация redis
	redis, err := redisadapter.NewRedis(rctx)
	if err != nil {
		log.Fatal("ошибка инициализации redis:", err)
	}

	//Инициализация трейсера
	tp, err := opentelimetry.NewTrace(context.Background(), "spotService")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	//Инициализация нового grpc сервера
	grpcServer, err := grpcserver.New(redis)

	//Регистрация grpc методов в сервере
	spotAPI.RegisterSpotInstrumentServiceServer(grpcServer, handlers)

	//Запуск работы сервера
	fmt.Println("Сервер запущен на порту 8080...")
	err = grpcServer.Serve(grpcServer.Listener)
	if err != nil {
		log.Fatal("ошибка работы сервера:", err)
	}
	wg.Wait()
}
