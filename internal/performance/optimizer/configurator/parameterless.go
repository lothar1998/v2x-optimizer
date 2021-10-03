package configurator

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/wrapper"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

type Parameterless struct {
	wrapper.IdentifiableOptimizer
}

func NewParameterless(optimizer optimizer.Optimizer) *Parameterless {
	return &Parameterless{&wrapper.Identifiable{Optimizer: optimizer}}
}

func (p *Parameterless) Builder() BuildFunc {
	return func(_ *cobra.Command) (wrapper.IdentifiableOptimizer, error) {
		return p.IdentifiableOptimizer, nil
	}
}

func (p *Parameterless) SetUpFlags(_ *cobra.Command) {}

func (p *Parameterless) TypeName() string {
	return p.IdentifiableOptimizer.Identifier()
}
