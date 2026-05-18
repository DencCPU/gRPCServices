package postgres

import "context"

// Обновление статуса заказа
func (p *PostgresDB) UpdateStatus(ctx context.Context, id int, status string) error {
	_, err := p.db.Exec(ctx, `
	UPDATE orders SET status = $1
	WHERE id = $2
	`, status, id)
	if err != nil {
		return err
	}
	return nil
}
