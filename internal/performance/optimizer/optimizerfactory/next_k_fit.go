package optimizerfactory

import (
	identifiable "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const nextKFitParameterK = "nkf_k"

type NextKFit struct{}

func (n NextKFit) Builder() BuildFunc {
	return func(command *cobra.Command) (identifiable.IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(nextKFitParameterK)
		if err != nil {
			return nil, err
		}
		return &identifiable.IdentifiableWrapper{Optimizer: &optimizer.NextKFit{K: int(k)}}, nil
	}
}

func (n NextKFit) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(nextKFitParameterK, "", 1, "NextKFit k parameter")
}

func (n NextKFit) Identifier() string {
	return "NextKFit"
}
