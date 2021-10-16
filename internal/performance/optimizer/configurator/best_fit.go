package configurator

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/bestfit"
	"github.com/spf13/cobra"
)

const (
	bestFitParameterFunctionID = "bf_fit"
	bestFitName                = "BestFit"
)

type BestFitWrapper struct {
	Name          string `id_name:""`
	FitnessFuncID int    `id_include:"true"`
	bestfit.BestFit
}

type BestFitConfigurator struct{}

func (n BestFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (optimizer.IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(bestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}

		if k > 4 {
			return nil, errors.New("unsupported fitness function")
		}

		bf := &BestFitWrapper{
			Name:          bestFitName,
			FitnessFuncID: int(k),
			BestFit:       bestfit.BestFit{FitnessFunc: intToFitness(int(k))},
		}

		return &optimizer.IdentifiableAdapter{Optimizer: bf}, nil
	}
}

func (n BestFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bestFitParameterFunctionID, "", 0,
		"BestFit fitness function parameter:\n"+
			"0 - classic fitness function\n"+
			"1 - take into account bucket size\n"+
			"2 - take into account left space in bucket and prefer big items\n"+
			"3 - take into account left space in bucket and prefer small items\n"+
			"4 - take into account left space in bucket and prefer as little space left as possible"+
			" before and after item assignment\n"+
			"(default 0)")
}

func (n BestFitConfigurator) TypeName() string {
	return bestFitName
}

func intToFitness(intValue int) bestfit.FitnessFunc {
	switch intValue {
	case 0:
		return bestfit.FitnessClassic
	case 1:
		return bestfit.FitnessWithBucketSize
	case 2:
		return bestfit.FitnessWithBucketLeftSpacePreferringBigItems
	case 3:
		return bestfit.FitnessWithBucketLeftSpacePreferringSmallItems
	case 4:
		return bestfit.FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment
	default:
		return nil
	}
}
