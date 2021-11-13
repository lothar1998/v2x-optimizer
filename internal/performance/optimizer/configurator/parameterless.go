package configurator

import (
	adapter "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

type Parameterless struct {
	adapter.IdentifiableOptimizer
}

func NewParameterless(optimizer optimizer.Optimizer) *Parameterless {
	return &Parameterless{&adapter.IdentifiableAdapter{Optimizer: optimizer}}
}

func (p *Parameterless) Builder() BuildFunc {
	return func(_ *cobra.Command) (adapter.IdentifiableCacheableOptimizer, error) {
		return &adapter.CacheableAdapter{
			IsCacheEligible:       true,
			IdentifiableOptimizer: p.IdentifiableOptimizer,
		}, nil
	}
}

func (p *Parameterless) SetUpFlags(_ *cobra.Command) {}

func (p *Parameterless) TypeName() string {
	return p.IdentifiableOptimizer.Identifier()
}
