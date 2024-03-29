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
	return func(command *cobra.Command) (optimizer.PerformanceSubjectOptimizer, error) {
		fitnessID, err := command.Flags().GetUint(bestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}

		if fitnessID > 6 {
			return nil, errors.New("unsupported fitness function")
		}

		bf := &BestFitWrapper{
			Name:          bestFitName,
			FitnessFuncID: int(fitnessID),
			BestFit:       bestfit.BestFit{FitnessFunc: intToFitness(fitnessID)},
		}

		return optimizer.NewPerformanceSubjectAdapter(bf, true), nil
	}
}

func (n BestFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bestFitParameterFunctionID, "", 0,
		"BestFit fitness function:\n"+
			"\t0 - classic fitness function\n"+
			"\t1 - take into account bucket size\n"+
			"\t2 - take into account left space in bucket and prefer big items\n"+
			"\t3 - take into account left space in bucket and prefer small items\n"+
			"\t4 - take into account left space in bucket and prefer small items and punish perfectly fitted items\n"+
			"\t5 - take into account left space in bucket and prefer as little space left as possible"+
			" before and after item assignment\n"+
			"\t6 - take into account left space in bucket and prefer as little space left as possible"+
			" before and after item assignment and punish perfectly fitted items\n"+
			"\t(default 0)")
}

func (n BestFitConfigurator) TypeName() string {
	return bestFitName
}

func intToFitness(intValue uint) bestfit.FitnessFunc {
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
		return bestfit.FitnessWithBucketLeftSpacePreferringSmallItemsPunishPerfectlyFittedItems
	case 5:
		return bestfit.FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignment
	case 6:
		return bestfit.FitnessWithBucketLeftSpacePreferringLittleSpaceBeforeAndAfterAssignmentPunishPerfectlyFittedItems
	default:
		return nil
	}
}
