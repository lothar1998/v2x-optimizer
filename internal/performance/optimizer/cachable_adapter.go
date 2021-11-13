package optimizer

type Cacheable interface {
	CacheEligible() bool
}

type IdentifiableCacheableOptimizer interface {
	IdentifiableOptimizer
	Cacheable
}

type CacheableAdapter struct {
	IsCacheEligible bool
	IdentifiableOptimizer
}

func (c *CacheableAdapter) CacheEligible() bool {
	return c.IsCacheEligible
}
