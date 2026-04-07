package postgres

import (
	"context"
	"fmt"

	userconfig "github.com/DencCPU/gRPCServices/UserService/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresDB struct {
	*pgxpool.Pool
	logger *zap.Logger
}

func NewDB(ctx context.Context, logger *zap.Logger, cfg userconfig.Postgres) (*PostgresDB, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Sslmode,
	)
	fmt.Println(dsn)
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres database is unavailable:%w", err)
	}

	// Проверка соединения
	conn, err := db.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire connection from pool:%w", err)
	}
	defer conn.Release()

	return &PostgresDB{db, logger}, nil
}
