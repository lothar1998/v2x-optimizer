package optimizer

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const (
	bestFitParameterFunctionID = "bf_fit"
	bestFitName                = "BestFit"
)

type BestFitWrapper struct {
	Name          string `id_name:""`
	FitnessFuncID int    `id_include:"true"`
	optimizer.BestFit
}

type BestFitConfigurator struct{}

func (n BestFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(bestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}

		if k > 2 {
			return nil, errors.New("unsupported fitness function")
		}

		bf := &BestFitWrapper{
			Name:          bestFitName,
			FitnessFuncID: int(k),
			BestFit:       optimizer.BestFit{FitnessFunc: intToFitness(int(k))},
		}

		return &IdentifiableAdapter{Optimizer: bf}, nil
	}
}

func (n BestFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bestFitParameterFunctionID, "", 0,
		"BestFit fitness function parameter:\n"+
			"0 - classic fitness function\n"+
			"1 - take into account bucket size\n"+
			"2 - take into account left space in bucket\n"+
			"(default 0)")
}

func (n BestFitConfigurator) TypeName() string {
	return bestFitName
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
