package optimizerfactory

import (
	identifiable "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const bestFitParameterFunctionID = "bf_fit"

type BestFit struct{}

func (n BestFit) Builder() BuildFunc {
	return func(command *cobra.Command) (identifiable.IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(bestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}
		return &identifiable.IdentifiableWrapper{Optimizer: &optimizer.BestFit{FitnessFuncID: int(k)}}, nil
	}
}

func (n BestFit) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bestFitParameterFunctionID, "", 0,
		"BestFit fitness function parameter:\n"+
			"0 - classic fitness function\n"+
			"1 - take into account bucket size\n"+
			"2 - take into account left space in bucket\n"+
			"(default 0)")
}

func (n BestFit) Identifier() string {
	return "BestFit"
}
