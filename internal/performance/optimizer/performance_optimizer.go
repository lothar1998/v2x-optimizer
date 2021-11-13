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
	identifiableAdapter
	isCacheEligible bool
}

func NewPerformanceAdapter(optimizer optimizer.Optimizer, isCacheEligible bool) PerformanceOptimizer {
	return &performanceAdapter{
		identifiableAdapter: identifiableAdapter{optimizer},
		isCacheEligible:     isCacheEligible,
	}
}

func (p *performanceAdapter) Identifier() string {
	return p.identifiableAdapter.Identifier()
}

func (p *performanceAdapter) CacheEligible() bool {
	return p.isCacheEligible
}
