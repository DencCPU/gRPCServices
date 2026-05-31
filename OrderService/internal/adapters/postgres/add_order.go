package postgres

import (
	"fmt"
	"time"

	"context"

	postgresdto "github.com/DencCPU/gRPCServices/OrderService/internal/adapters/dto/postgres"
	ordererrors "github.com/DencCPU/gRPCServices/OrderService/internal/domain/error"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Добавление заказа в БД (таблица orers)
func (p *PostgresDB) AddOrderStorage(ctx context.Context, newOrder orderdomain.Order) (string, string, error) {

	dto := postgresdto.CreatOrderDTO(newOrder)

	//Transaction begin
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return "", "", fmt.Errorf("error starting transaction:%w", err)
	}

	//Add orderID to database
	refOrderId, orderID, err := p.AddOrderID(tx, ctx, newOrder)
	if err != nil {
		return "", "", fmt.Errorf("error adding OrderID:%w", err)
	}

	//Add user to database
	refUserID, err := p.AddUserID(tx, ctx, newOrder)
	if err != nil {
		return "", "", fmt.Errorf("error adding UserID:%w", err)
	}

	//Add market to database
	refMarketID, err := p.AddMarketID(tx, ctx, newOrder)
	if err != nil {
		return "", "", fmt.Errorf("error adding UserID:%w", err)
	}

	//Add reference to order dto
	dto.RefUserId = refUserID
	dto.RefMarketId = refMarketID
	dto.RefOrderId = refOrderId
	dto.Status = orderdomain.StatusCreated
	dto.CreatedAt = time.Now()

	//Добавление заказа
	var orderStatus string
	err = tx.QueryRow(ctx, `
	INSERT INTO orders(user_id,market_id,order_type,price,quantity,status,order_id,created_at)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	RETURNING status
	`,
		dto.RefUserId,
		dto.RefMarketId,
		dto.OrderType,
		dto.Price,
		dto.Quantity,
		dto.Status,
		dto.RefOrderId,
		dto.CreatedAt,
	).Scan(&orderStatus)
	if err != nil {
		return "", "", fmt.Errorf("Error adding order to database:%w", err)
	}

	orderInfo := orderdomain.OrderInfo{
		OrderType: newOrder.OrderType,
		UserId:    newOrder.UserId,
		OrderId:   orderID,
	}

	p.controlOrderChan <- orderInfo

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	if err = tx.Commit(ctx); err != nil {
		return "", "", err
	}

	return orderID, orderStatus, nil
}

// Добавленение рынка в БД (таблица markets)
func (p *PostgresDB) AddMarketID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order) (int, error) {
	var marketName string

	p.marketMu.RLock()
	marketName = p.marketCache[newOrder.MarketId].MarketName
	p.marketMu.RUnlock()

	//Инициализация DTO
	dto := postgresdto.CreateMarketDTO(newOrder.MarketId, marketName)
	dto.CreatedAt = time.Now()

	//Добавление маркета
	var id int
	err := tx.QueryRow(ctx, `
		INSERT INTO markets (name,market_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (market_id) DO UPDATE
		SET market_id = EXCLUDED.market_id
		RETURNING id
	`,
		dto.MarketName,
		dto.MarketId,
		dto.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, nil
}

// Add orderID to database (orders_id table)
func (p *PostgresDB) AddOrderID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order) (int, string, error) {

	var foundMarket bool //fing market flag

	//Check market cache
	for _, m := range p.marketCache {
		p.marketMu.RLock()

		if m.MarketId == newOrder.MarketId && m.UserAccess == newOrder.UserRole {
			foundMarket = true
			p.marketMu.RUnlock()
			break
		}
		p.marketMu.RUnlock()

	}

	if foundMarket != true {
		return 0, "", ordererrors.Avalible_markets
	}

	//New Order_id
	var orderID string
	orderID = uuid.New().String()

	//DTO
	dto := postgresdto.CreateOrders_idDTO(orderID)
	dto.CreatedAt = time.Now()

	//Add record to database
	var id int
	err := tx.QueryRow(ctx, ` 
	INSERT INTO orders_id(order_id,created_at) 
	VALUES ($1,$2)
	RETURNING id
	`,
		dto.OrderId,
		dto.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, "", err
	}

	return id, orderID, nil
}

// Add userID to database (users table)
func (p *PostgresDB) AddUserID(tx pgx.Tx, ctx context.Context, newOrder orderdomain.Order) (int, error) {
	//DTO
	dto := postgresdto.CreateUserDTO(newOrder.UserId)
	dto.CreatedAt = time.Now()

	var id int
	err := tx.QueryRow(ctx, `
		INSERT INTO users (user_id, created_at)
		VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE
		SET user_id = EXCLUDED.user_id
		RETURNING id
	`,
		dto.UserId,
		dto.CreatedAt,
	).Scan(&id)

	if err != nil {
		return 0, err
	}

	return id, nil
}
