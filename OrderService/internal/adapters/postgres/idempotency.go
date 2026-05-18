package postgres

import (
	"time"
)

func (p *PostgresDB) IdempotencyCheck(idepotencyKey string) bool {
	p.cacheMu.RLock()
	defer p.cacheMu.RUnlock()
	if _, exist := p.idempotecyCache[idepotencyKey]; !exist {
		p.idempotecyCache[idepotencyKey] = time.Now().Add(p.cacheTTL)
		return true
	}
	return false
}

func (p *PostgresDB) CheckCacheTTL() {

	p.wg.Add(1)
	go func() {
		defer p.wg.Done()

		ticker := time.NewTicker(p.cacheTTL)
		defer ticker.Stop()

		for {
			select {
			case <-p.ctx.Done():
				return
			case <-ticker.C:
				for key, ttl := range p.idempotecyCache {
					p.cacheMu.Lock()
					if time.Now().After(ttl) {
						delete(p.idempotecyCache, key)
					}
					p.cacheMu.Unlock()
				}
			}
		}
	}()
}
