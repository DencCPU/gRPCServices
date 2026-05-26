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
	db                  *pgxpool.Pool
	notify              Notify
	controlOrderChan    chan orderdomain.OrderInfo
	idempotecyCache     map[string]time.Time
	idempotencyCacheTTL time.Duration
	idempotencyCacheMu  sync.RWMutex
	marketCache         map[string]orderdomain.Market
	marketMu            sync.RWMutex
	marketCacheTTL      time.Duration

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
	postgres.idempotencyCacheTTL = cfg.IdempotencyCacheTTL
	fmt.Println("TTL:", postgres.idempotencyCacheTTL)
	return &PostgresDB{
		db:                  dataBase,
		notify:              notify,
		controlOrderChan:    make(chan orderdomain.OrderInfo, cfg.ControlChanSize),
		idempotecyCache:     map[string]time.Time{},
		idempotencyCacheTTL: cfg.IdempotencyCacheTTL,
		idempotencyCacheMu:  sync.RWMutex{},
		marketCache:         make(map[string]orderdomain.Market),
		marketMu:            sync.RWMutex{},
		marketCacheTTL:      cfg.MarketCacheTTL,
		ctx:                 dbCtx,
		cancel:              dbCancel,
		wg:                  sync.WaitGroup{}}, nil
}

func (p *PostgresDB) GetPgxPool() *pgxpool.Pool {
	return p.db
}
