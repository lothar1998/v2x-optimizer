package configurator

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/bucketpoolbestfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
	"github.com/spf13/cobra"
)

const (
	bucketPoolBestFitParameterFunctionID              = "bpbf_fit"
	bucketPoolBestFitParameterBucketReorderFunctionID = "bpbf_bucket_reorder_fun"
	bucketPoolBestFitParameterInitPoolSize            = "bpbf_init_pool_size"
	bucketPoolBestFitName                             = "BucketPoolBestFit"
)

type BucketPoolBestFitWrapper struct {
	Name                string `id_name:""`
	FitnessFuncID       int    `id_include:"true"`
	BucketReorderFuncID int    `id_include:"true"`
	InitPoolSize        int    `id_include:"true"`
	bucketpoolbestfit.BucketPoolBestFit
}

type BucketPoolBestFitConfigurator struct{}

func (b BucketPoolBestFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (optimizer.PerformanceSubjectOptimizer, error) {
		fitnessID, err := command.Flags().GetUint(bucketPoolBestFitParameterFunctionID)
		if err != nil {
			return nil, err
		}

		if fitnessID > 6 {
			return nil, errors.New("unsupported fitness function")
		}

		bucketReorderID, err := command.Flags().GetUint(bucketPoolBestFitParameterBucketReorderFunctionID)
		if err != nil {
			return nil, err
		}

		if bucketReorderID > 3 {
			return nil, errors.New("unsupported bucket reorder function")
		}

		initPoolSize, err := command.Flags().GetUint(bucketPoolBestFitParameterInitPoolSize)
		if err != nil {
			return nil, err
		}

		bpbf := BucketPoolBestFitWrapper{
			Name:                bucketPoolBestFitName,
			FitnessFuncID:       int(fitnessID),
			BucketReorderFuncID: int(bucketReorderID),
			InitPoolSize:        int(initPoolSize),
			BucketPoolBestFit: bucketpoolbestfit.BucketPoolBestFit{
				InitPoolSize:       int(initPoolSize),
				ReorderBucketsFunc: bucketPoolBestFitToBucketReorderFunc(bucketReorderID),
				FitnessFunc:        intToFitness(fitnessID),
			},
		}

		return optimizer.NewPerformanceSubjectAdapter(bpbf, bucketReorderID != 3), nil
	}
}

func (b BucketPoolBestFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bucketPoolBestFitParameterFunctionID, "", 0,
		"BucketPoolBestFit fitness function:\n"+
			"0 - classic fitness function\n"+
			"1 - take into account bucket size\n"+
			"2 - take into account left space in bucket and prefer big items\n"+
			"3 - take into account left space in bucket and prefer small items\n"+
			"4 - take into account left space in bucket and prefer small items and punish perfectly fitted items\n"+
			"5 - take into account left space in bucket and prefer as little space left as possible"+
			" before and after item assignment\n"+
			"6 - take into account left space in bucket and prefer as little space left as possible"+
			" before and after item assignment and punish perfectly fitted items\n"+
			"(default 0)")
	command.Flags().UintP(bucketPoolBestFitParameterBucketReorderFunctionID, "", 0,
		"BucketPoolBestFit bucket reorder function (defines order in which items are added to bucket pool):\n"+
			"0 - no op (order defined by input data)\n"+
			"1 - sort buckets in ascending order by size\n"+
			"2 - sort buckets in descending order by size\n"+
			"3 - random order\n"+
			"(default 0)")
	command.Flags().UintP(bucketPoolBestFitParameterInitPoolSize, "", 1,
		"BucketPoolBestFit init bucket pool size (default 1)")
}

func (b BucketPoolBestFitConfigurator) TypeName() string {
	return bucketPoolBestFitName
}

func bucketPoolBestFitToBucketReorderFunc(intValue uint) helper.ReorderBucketsFunc {
	switch intValue {
	case 0:
		return helper.NoOpReorder
	case 1:
		return helper.AscendingBucketSizeReorder
	case 2:
		return helper.DescendingBucketSizeReorder
	case 3:
		return helper.RandomReorder
	default:
		return nil
	}
}
