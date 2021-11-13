package common

type Cacheable interface {
	CacheEligible() bool
}
