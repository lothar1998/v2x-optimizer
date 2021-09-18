package optimizerfactory

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/spf13/cobra"
)

type Parameterless struct {
	optimizer.Identifiable
}

func (p *Parameterless) Builder() BuildFunc {
	return func(_ *cobra.Command) (optimizer.Identifiable, error) {
		return p.Identifiable, nil
	}
}

func (p *Parameterless) SetUpFlags(_ *cobra.Command) {}

func (p *Parameterless) Name() string {
	return p.Identifier()
}
