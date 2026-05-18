package postgres

import (
	"context"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

func (p *PostgresDB) GetOrderState(ctx context.Context, key orderdomain.Key) (orderdomain.ReceivedOrderInfo, error) {
	var orderInfo orderdomain.ReceivedOrderInfo

	err := p.db.QueryRow(ctx, `
	SELECT 
	status,
	price,
	quantity,
	markets.name 
	FROM orders
	INNER JOIN users ON orders.user_id = users.id
	INNER JOIN orders_id ON orders.order_id = orders_id.id
	INNER JOIN markets ON orders.market_id = markets.id
	WHERE users.user_id = $1
	AND orders_id.order_id = $2
`, key.UserId, key.OrderId).
		Scan(
			&orderInfo.Status,
			&orderInfo.Price,
			&orderInfo.Quantity,
			&orderInfo.MarketName)
	if err != nil {
		return orderdomain.ReceivedOrderInfo{}, err
	}
	orderInfo.OrderId = key.OrderId
	return orderInfo, nil
}
