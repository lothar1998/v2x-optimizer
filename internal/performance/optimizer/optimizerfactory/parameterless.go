package optimizerfactory

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/spf13/cobra"
)

type Parameterless struct {
	optimizer.IdentifiableOptimizer
}

func (p *Parameterless) Builder() BuildFunc {
	return func(_ *cobra.Command) (optimizer.IdentifiableOptimizer, error) {
		return p.IdentifiableOptimizer, nil
	}
}

func (p *Parameterless) SetUpFlags(_ *cobra.Command) {}

func (p *Parameterless) Identifier() string {
	return p.IdentifiableOptimizer.Identifier()
}
