package optimizer

import (
	"github.com/lothar1998/v2x-optimizer/internal/common"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type PerformanceOptimizer interface {
	common.Identifiable
	common.Cacheable
	optimizer.Optimizer
}

type performanceAdapter struct {
	IdentifiableAdapter
	IsCacheEligible bool
}

func NewPerformanceAdapter(optimizer optimizer.Optimizer, isCacheEligible bool) PerformanceOptimizer {
	return &performanceAdapter{
		IdentifiableAdapter: IdentifiableAdapter{optimizer},
		IsCacheEligible:     isCacheEligible,
	}
}

func (p *performanceAdapter) Identifier() string {
	return p.IdentifiableAdapter.Identifier()
}

func (p *performanceAdapter) CacheEligible() bool {
	return p.IsCacheEligible
}
