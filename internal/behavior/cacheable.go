package behavior

type Cacheable interface {
	CacheEligible() bool
}
