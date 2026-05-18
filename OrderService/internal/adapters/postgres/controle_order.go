package postgres

import (
	"fmt"
	"time"

	ordererrors "github.com/DencCPU/gRPCServices/OrderService/internal/domain/error"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

func (p *PostgresDB) ControlOrder(errChan chan ordererrors.ErrStruct) {

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-p.ctx.Done():
			p.drainOrderChannel(errChan)
			return

		case orderInfo, ok := <-p.controlOrderChan:
			if !ok {
				errChan <- ordererrors.ErrStruct{Err: fmt.Errorf("the channel is close")}
				return
			}

			p.wg.Add(1)
			go func() {
				defer p.wg.Done()
				err := p.processOrder(orderInfo)
				if err != nil {
					errChan <- ordererrors.ErrStruct{
						OrderId: orderInfo.OrderId,
						Err:     err,
					}
				}
			}()

		case <-ticker.C:

		}
	}
}

func (p *PostgresDB) processOrder(orderInfo orderdomain.OrderInfo) error {
	var (
		id         int
		status     string
		statusChan = make(chan string, 1)
	)

	err := p.db.QueryRow(p.ctx, ` 
				SELECT orders.id
				FROM orders
				JOIN users ON orders.user_id = users.id
				JOIN orders_id ON orders.order_id = orders_id.id
				WHERE users.user_id = $1
				AND orders_id.order_id = $2
				`, orderInfo.UserId, orderInfo.OrderId).Scan(&id)
	if err != nil {
		return err
	}
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		p.notify.AddNewState(orderInfo.UserId, orderInfo.OrderId, statusChan)
	}()
	defer close(statusChan)

	switch orderInfo.OrderType {

	case orderdomain.ORDER_TYPE_NORMAL:
		status = orderdomain.StatusProcessing
		err := p.UpdateStatus(p.ctx, id, status)
		if err != nil {
			return err
		}
		statusChan <- status

		status = orderdomain.StatusComplited
		err = p.UpdateStatus(p.ctx, id, status)
		if err != nil {
			return err
		}
		statusChan <- status

	case orderdomain.ORDER_TYPE_EXPRESS:
		status = orderdomain.StatusProcessing
		err := p.UpdateStatus(p.ctx, id, status)
		if err != nil {
			return err
		}
		statusChan <- status

		status = orderdomain.StatusComplited
		err = p.UpdateStatus(p.ctx, id, status)
		statusChan <- status
	}
	return nil
}

func (p *PostgresDB) drainOrderChannel(errChan chan ordererrors.ErrStruct) {
	for {
		select {
		case orderInfo, ok := <-p.controlOrderChan:
			if !ok {
				errChan <- ordererrors.ErrStruct{Err: fmt.Errorf("the channel is close")}
				return
			}
			err := p.processOrder(orderInfo)
			if err != nil {
				errChan <- ordererrors.ErrStruct{
					OrderId: orderInfo.OrderId,
					Err:     err,
				}
				return
			}
		default:
			close(p.controlOrderChan)
			return
		}
	}

}

func (p *PostgresDB) StopControlOrder() {
	p.cancel()
	p.wg.Wait()
}
