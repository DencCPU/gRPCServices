package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func (p *PostgresDB) ControlOrder(orderType string, user_id int64, orderID string) chan string {
	var wg sync.WaitGroup
	stateCh := make(chan string, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		var id int
		var status string

		err := p.QueryRow(context.Background(), `
	SELECT orders.id
	FROM orders
	JOIN users ON orders.user_id = users.id
	JOIN orders_id ON orders.order_id = orders_id.id
	WHERE users.user_id = $1
	AND orders_id.order_id = $2
	`, user_id, orderID).Scan(&id)
		if err != nil {
			fmt.Println("ошибка обработки заказа:", err)
			return
		}

		switch orderType {

		case "normal":
			time.Sleep(5 * time.Second)
			status = "in progress"
			err = p.UpdateStatus(id, status)
			if err != nil {
				fmt.Println("ошибка обработки заказа:", err)
				return
			}
			stateCh <- status
			time.Sleep(5 * time.Second)
			status = "completed"
			err = p.UpdateStatus(id, status)
			if err != nil {
				fmt.Println("ошибка обработки заказа:", err)
				return
			}
			stateCh <- status

		case "express":
			time.Sleep(2 * time.Second)
			status = "in progress"
			err = p.UpdateStatus(id, status)
			if err != nil {
				fmt.Println("ошибка обработки заказа:", err)
				return
			}
			stateCh <- status

			time.Sleep(2 * time.Second)
			status = "completed"
			err = p.UpdateStatus(id, status)
			if err != nil {
				fmt.Println("ошибка обработки заказа:", err)
				return
			}
			stateCh <- status
		}
	}()
	go func() {
		wg.Wait()
		close(stateCh)
	}()
	return stateCh
}
