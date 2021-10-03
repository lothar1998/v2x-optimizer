package configurator

import (
	identifiable "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/wrapper"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const nextKFitParameterK = "nkf_k"

type NextKFit struct {
	K int `id_include:"true"`
	optimizer.NextKFit
}

type NextKFitConfigurator struct{}

func (n NextKFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (identifiable.IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(nextKFitParameterK)
		if err != nil {
			return nil, err
		}

		nkf := &NextKFit{
			K:        int(k),
			NextKFit: optimizer.NextKFit{K: int(k)},
		}
		return &identifiable.Identifiable{Optimizer: nkf}, nil
	}
}

func (n NextKFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(nextKFitParameterK, "", 1, "NextKFitConfigurator k parameter")
}

func (n NextKFitConfigurator) TypeName() string {
	return "NextKFit"
}
