package postgres

import (
	postgresdto "Academy/gRPCServices/OrderService/internal/adapters/dto/postgres"
	"Academy/gRPCServices/OrderService/internal/domain/order"
	"time"

	"context"
)

func (p *PostgresDB) AddOrderStorage(ctx context.Context, newOrder order.Order, marketsID []int64) (string, string, error) {

	dto := postgresdto.CreatOrderDTO(newOrder)

	//Начало транзакции
	tx, err := p.Begin(ctx)
	if err != nil {
		return "", "", err
	}

	//Добавление ID заказа
	refOrderId, orderID, err := p.AddOrderID(tx, ctx, newOrder, marketsID)
	if err != nil {
		tx.Rollback(ctx)
		return "", "", err
	}
	//Добавление пользователя
	refUserID, err := p.AddUserID(tx, ctx, newOrder)
	if err != nil {
		tx.Rollback(ctx)
		return "", "", err
	}

	//Добавление рынка
	refMarketID, err := p.AddMarketID(tx, ctx, newOrder)
	if err != nil {
		tx.Rollback(ctx)
		return "", "", err
	}

	//Добавление ссылок в users
	dto.Ref_User_Id = refUserID
	dto.Ref_Market_Id = refMarketID
	dto.Ref_Order_Id = refOrderId
	dto.Status = "created"
	dto.Created_at = time.Now()

	//Добавление заказа
	var order_status string
	err = tx.QueryRow(ctx, `
	INSERT INTO orders(user_id,market_id,order_type,price,quantity,status,order_id,created_at)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING status
	`, dto.Ref_User_Id, dto.Ref_Market_Id, dto.Order_type, dto.Price, dto.Quantity, dto.Status, dto.Ref_Order_Id, dto.Created_at).Scan(&order_status)
	if err != nil {
		tx.Rollback(ctx)
		return "", "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", err
	}

	return orderID, order_status, nil
}
