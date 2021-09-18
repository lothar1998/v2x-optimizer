package optimizerfactory

import (
	identifiable "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const nextKFitParameterK = "nkf_k"

type NextKFit struct{}

func (n NextKFit) Builder() BuildFunc {
	return func(command *cobra.Command) (identifiable.Identifiable, error) {
		k, err := command.Flags().GetUint(nextKFitParameterK)
		if err != nil {
			return nil, err
		}
		return &identifiable.IdentifiableOptimizer{Optimizer: &optimizer.NextKFit{K: int(k)}}, nil
	}
}

func (n NextKFit) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(nextKFitParameterK, "", 1, "k parameter")
}

func (n NextKFit) Name() string {
	return "NextKFit"
}
