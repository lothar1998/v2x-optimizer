package configurator

import (
	adapter "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

type Parameterless struct {
	adapter.PerformanceOptimizer
}

func NewParameterless(optimizer optimizer.Optimizer) *Parameterless {
	return &Parameterless{adapter.NewPerformanceAdapter(optimizer, true)}
}

func (p *Parameterless) Builder() BuildFunc {
	return func(_ *cobra.Command) (adapter.PerformanceOptimizer, error) {
		return p.PerformanceOptimizer, nil
	}
}

func (p *Parameterless) SetUpFlags(_ *cobra.Command) {}

func (p *Parameterless) TypeName() string {
	return p.PerformanceOptimizer.Identifier()
}
