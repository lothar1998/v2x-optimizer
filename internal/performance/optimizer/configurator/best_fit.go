package configurator

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/wrapper"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const bestFitParameterFunctionID = "bf_fit"

type BestFit struct {
	FitnessFuncID int `id_include:"true"`
	optimizer.BestFit
}

type BestFitConfigurator struct{}

func (n BestFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (wrapper.IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(bestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}

		bf := &BestFit{
			FitnessFuncID: int(k),
			BestFit:       optimizer.BestFit{FitnessFunc: intToFitness(int(k))},
		}

		return &wrapper.Identifiable{Optimizer: bf}, nil
	}
}

func (n BestFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bestFitParameterFunctionID, "", 0,
		"BestFitConfigurator fitness function parameter:\n"+
			"0 - classic fitness function\n"+
			"1 - take into account bucket size\n"+
			"2 - take into account left space in bucket\n"+
			"(default 0)")
}

func (n BestFitConfigurator) TypeName() string {
	return "BestFit"
}

func intToFitness(intValue int) optimizer.BestFitFitnessFunc {
	switch intValue {
	case 0:
		return optimizer.BestFitFitnessClassic
	case 1:
		return optimizer.BestFitFitnessWithBucketSize
	case 2:
		return optimizer.BestFitFitnessWithBucketLeftSpace
	default:
		return nil
	}
}
