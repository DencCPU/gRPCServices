package postgres

import (
	"time"

	orderdomain "github.com/DencCPU/gRPCServices/OrderService/internal/domain/order"
)

func (p *PostgresDB) UpdateMarketCache(markets []orderdomain.Market) {
	for _, m := range markets {

		p.marketMu.Lock()

		p.marketCache[m.MarketId] = orderdomain.Market{
			MarketId:   m.MarketId,
			MarketName: m.MarketName,
			UserAccess: m.UserAccess,
			TTL:        time.Now().Add(p.marketCacheTTL),
		}
		p.marketMu.Unlock()

	}

}

func (p *PostgresDB) RemoveIrrelevantMarkets() {
	for key, value := range p.marketCache {

		p.marketMu.Lock()
		if time.Now().After(value.TTL) {
			delete(p.marketCache, key)
		}
		p.marketMu.Unlock()

	}
}
