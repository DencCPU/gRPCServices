package main

import (
	"Academy/gRPCServices/OrderService/pkg/orderclient"
	orderAPI "Academy/gRPCServices/Protobuf/gen/order"
	"Academy/gRPCServices/SpotInstrumentService/pkg/opentelimetry"
	"io"

	"context"
	"fmt"
	"log"
	"time"

	"go.opentelemetry.io/otel"
)

func main() {
	tp, err := opentelimetry.NewTrace(context.Background(), "Client")
	if err != nil {
		log.Fatal(err)
	}
	defer tp.Shutdown(context.Background())

	client, err := orderclient.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	var (
		userID    int64
		marketID  int64
		orderType string
		price     float64
		quantity  int64
	)

	fmt.Println("Доступные рынки:\n0.Yandex Market;\n1.OZON;\n2.Wildberris;\n3.Aliexpress.")
	fmt.Print("Введите user_id:")
	fmt.Scanln(&userID)
	fmt.Println("Введите market_id из предоставленного списка:")
	fmt.Scanln(&marketID)
	fmt.Println("Укажите тип заказа (normal,express):")
	fmt.Scanln(&orderType)
	fmt.Println("Укажите price:")
	fmt.Scanln(&price)
	fmt.Println("Укажите кол-во товара:")
	fmt.Scanln(&quantity)

	tr := otel.Tracer("main-client")

	ctx, span := tr.Start(context.Background(), "CreateOrder")
	defer span.End()

	ctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	createResp, err := client.CreateOrder(ctx, &orderAPI.CreateReq{UserId: userID, MarketId: marketID, OrderType: orderType, Price: price, Quantity: quantity})
	if err != nil {
		log.Fatal(err)
	}

	orderId := createResp.OrderId
	fmt.Printf("Заказ:%s. Статус:%s.\n", orderId, createResp.Status)

	streamResp, err := client.StreamOrderUpdate(ctx, &orderAPI.GetReq{OrderId: orderId, UserId: userID})
	for {
		if err != io.EOF {
			resp, err := streamResp.Recv()
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(resp.OrderStatus)
		} else {
			break
		}
	}

	time.Sleep(3 * time.Second)
	ctx, span = tr.Start(context.Background(), "GetOrder")
	defer span.End()
	getResp, err := client.GetOrderStatus(ctx, &orderAPI.GetReq{OrderId: orderId, UserId: userID})
	if err != nil {
		log.Fatal("Ошибка при запросе GetOrderStatus:", err)
	}
	fmt.Println("Статус заказа:", getResp.OrderStatus)
}
