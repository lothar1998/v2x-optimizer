package optimizer

import (
	"github.com/lothar1998/v2x-optimizer/internal/behavior"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type PerformanceSubjectOptimizer interface {
	behavior.Identifiable
	behavior.Cacheable
	optimizer.Optimizer
}

type performanceSubjectAdapter struct {
	*identifiableOptimizerAdapter
	isCacheEligible bool
}

func NewPerformanceSubjectAdapter(optimizer optimizer.Optimizer, isCacheEligible bool) PerformanceSubjectOptimizer {
	return &performanceSubjectAdapter{
		identifiableOptimizerAdapter: &identifiableOptimizerAdapter{optimizer},
		isCacheEligible:              isCacheEligible,
	}
}

func (p *performanceSubjectAdapter) Identifier() string {
	return p.identifiableOptimizerAdapter.Identifier()
}

func (p *performanceSubjectAdapter) CacheEligible() bool {
	return p.isCacheEligible
}
