package postgres

import (
	"time"

	"context"

	postgresdto "github.com/DencCPU/gRPCServices/OrderService/internal/adapters/dto/postgres"
	ordererrors "github.com/DencCPU/gRPCServices/OrderService/internal/domain/error"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// Добавление заказа в БД (таблица orers)
func (p *PostgresDB) AddOrderStorage(ctx context.Context, newOrder orderdomain.Order, markets []orderdomain.Market) (string, string, error) {

	dto := postgresdto.CreatOrderDTO(newOrder)

	//Начало транзакции
	tx, err := p.Begin(ctx)
	if err != nil {
		p.logger.Error("error starting transaction:",
			zap.Error(err),
		)
		return "", "", err
	}

	//Добавление ID заказа
	refOrderId, orderID, err := p.AddOrderID(tx, ctx, newOrder, markets)
	if err != nil {
		p.logger.Error("error adding OrderID:",
			zap.Error(err),
		)
		tx.Rollback(ctx)
		return "", "", err
	}
	//Добавление пользователя
	refUserID, err := p.AddUserID(tx, ctx, newOrder)
	if err != nil {
		p.logger.Error("error adding UserID:",
			zap.Error(err),
		)
		tx.Rollback(ctx)
		return "", "", err
	}

	//Добавление рынка
	refMarketID, err := p.AddMarketID(tx, ctx, newOrder, markets)
	if err != nil {
		p.logger.Error("error adding MarketID:",
			zap.Error(err),
		)
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
		p.logger.Error("Error adding order to database:",
			zap.Error(err),
		)
		return "", "", err
	}

	if err := tx.Commit(ctx); err != nil {
		return "", "", err
	}

	return orderID, order_status, nil
}

// Добавленение рынка в БД (таблица markets)
func (p *PostgresDB) AddMarketID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order, markets []orderdomain.Market) (int, error) {
	var marketName string

	for _, m := range markets {
		if m.ID == newOrder.Market_id {
			marketName = m.Name
		}
	}

	//Инициализация DTO
	dto := postgresdto.CreateMarketDTO(newOrder.Market_id, marketName)
	dto.Created_at = time.Now()

	//Добавление маркета
	var id int
	err := tx.QueryRow(ctx, `
		INSERT INTO markets (name,market_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (market_id) DO UPDATE
		SET market_id = EXCLUDED.market_id
		RETURNING id
	`, dto.Market_name, dto.Market_id, dto.Created_at).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// Добавление OrderID в БД (таблица orders_id)
func (p *PostgresDB) AddOrderID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order, markets []orderdomain.Market) (int, string, error) {

	var foundMarket bool //Флаг, показывающий найден нужный рынок или нет

	//Проверка наличия нужного рынка
	for _, m := range markets {
		if m.ID == newOrder.Market_id {
			foundMarket = true
			break
		}
	}
	if foundMarket != true {
		return 0, "", ordererrors.Avalible_markets
	}

	//Присвоение нового id заказа
	var orderID string
	orderID = uuid.New().String()

	//Ининциализация DTO
	dto := postgresdto.CreateOrders_idDTO(orderID)
	dto.Created_at = time.Now()

	//Запись в БД
	var id int
	err := tx.QueryRow(ctx, ` 
	INSERT INTO orders_id(order_id,created_at) 
	VALUES ($1,$2)
	RETURNING id
	`, dto.Order_id, dto.Created_at).Scan(&id)
	if err != nil {
		return 0, "", err
	}

	return id, orderID, nil
}

// Добавление UserID в БД (таблица users)
func (p *PostgresDB) AddUserID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order) (int, error) {
	//Инициализация DTO
	dto := postgresdto.CreateUserDTO(newOrder.User_id)
	dto.Created_at = time.Now()

	//Поиск пользователя с id
	var id int
	err := tx.QueryRow(ctx, `
		INSERT INTO users (user_id, created_at)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET user_id = EXCLUDED.user_id
		RETURNING id
	`, dto.User_id, dto.Created_at).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
