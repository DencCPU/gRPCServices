package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	orderconfig "github.com/DencCPU/gRPCServices/OrderService/config"
	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Notify interface {
	AddNewState(userId string, orderId string, statCh chan string)
}

type PostgresDB struct {
	db               *pgxpool.Pool
	notify           Notify
	controlOrderChan chan orderdomain.OrderInfo
	idempotecyCache  map[string]time.Time
	cacheTTL         time.Duration
	cacheMu          sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

func NewDB(ctx context.Context, cfg orderconfig.Postgres, notify Notify) (*PostgresDB, error) {

	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Name,
		cfg.Sslmode,
	)
	dataBase, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("postgres database is unavailable:%w", err)
	}

	// Connection check
	conn, err := dataBase.Acquire(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot acquire connection from pool:%w", err)
	}
	defer conn.Release()

	dbCtx, dbCancel := context.WithCancel(ctx)
	postgres := PostgresDB{}
	postgres.cacheTTL = cfg.IdempotencyCacheTTL
	fmt.Println("TTL:", postgres.cacheTTL)
	return &PostgresDB{
		db:               dataBase,
		notify:           notify,
		controlOrderChan: make(chan orderdomain.OrderInfo, cfg.ControlChanSize),
		idempotecyCache:  map[string]time.Time{},
		cacheTTL:         cfg.IdempotencyCacheTTL,
		cacheMu:          sync.RWMutex{},
		ctx:              dbCtx,
		cancel:           dbCancel,
		wg:               sync.WaitGroup{}}, nil
}

func (p *PostgresDB) GetPgxPool() *pgxpool.Pool {
	return p.db
}
