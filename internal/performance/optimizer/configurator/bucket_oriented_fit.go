package configurator

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/bucketorientedfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/helper"
	"github.com/spf13/cobra"
)

const (
	bucketOrientedFitParameterBucketReorderFunctionID = "bof_bucket_reorder_fun"
	bucketOrientedFitParameterItemsReorderFunctionID  = "bof_items_order"
	bucketOrientedFitName                             = "BucketOrientedFit"
)

type BucketOrientedFitWrapper struct {
	Name                string `id_name:""`
	BucketReorderFuncID int    `id_include:"true"`
	ItemReorderFuncID   int    `id_include:"true"`
	bucketorientedfit.BucketOrientedFit
}

type BucketOrientedFitConfigurator struct{}

func (b BucketOrientedFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (optimizer.PerformanceSubjectOptimizer, error) {
		bucketReorderID, err := command.Flags().GetUint(bucketOrientedFitParameterBucketReorderFunctionID)
		if err != nil {
			return nil, err
		}

		if bucketReorderID > 4 {
			return nil, errors.New("unsupported bucket reorder function")
		}

		itemReorderID, err := command.Flags().GetUint(bucketOrientedFitParameterItemsReorderFunctionID)
		if err != nil {
			return nil, err
		}

		if itemReorderID > 1 {
			return nil, errors.New("unsupported item reorder function")
		}

		bof := BucketOrientedFitWrapper{
			Name:                bucketOrientedFitName,
			BucketReorderFuncID: int(bucketReorderID),
			ItemReorderFuncID:   int(itemReorderID),
			BucketOrientedFit: bucketorientedfit.BucketOrientedFit{
				ReorderBucketsByItemsFunc: bucketOrientedFitToBucketReorderFunc(bucketReorderID),
				ItemOrderComparatorFunc:   bucketOrientedFitToItemOrderComparatorFunc(itemReorderID),
			},
		}

		return optimizer.NewPerformanceSubjectAdapter(bof, true), nil
	}
}

func (b BucketOrientedFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(bucketOrientedFitParameterBucketReorderFunctionID, "", 0,
		"BucketOrientedFit buckets reorder function (defines order in which buckets are used):\n"+
			"\t0 - no op (order defined by input data)\n"+
			"\t1 - sort buckets in ascending order by the total sum of possible items' sizes in the bucket\n"+
			"\t2 - sort buckets in descending order by the total sum of possible items' sizes in the bucket\n"+
			"\t3 - sort buckets in ascending order by the total sum of possible items' sizes divided by bucket size\n"+
			"\t4 - sort buckets in descending order by the total sum of possible items' sizes divided by bucket size\n"+
			"\t(default 0)")
	command.Flags().UintP(bucketOrientedFitParameterItemsReorderFunctionID, "", 0,
		"BucketOrientedFit items reorder function (defines order in which items are added to bucket):\n"+
			"\t0 - ascending order of items by their size\n"+
			"\t1 - descending order of items by their size\n"+
			"\t(default 0)")
}

func (b BucketOrientedFitConfigurator) TypeName() string {
	return bucketOrientedFitName
}

func bucketOrientedFitToBucketReorderFunc(intValue uint) helper.ReorderBucketsByItemsFunc {
	switch intValue {
	case 0:
		return helper.NoOpReorderByItems
	case 1:
		return helper.AscendingTotalSizeOfItemsInBucketReorder
	case 2:
		return helper.DescendingTotalSizeOfItemsInBucketReorder
	case 3:
		return helper.AscendingRelativeSizeReorder
	case 4:
		return helper.DescendingRelativeSizeReorder
	default:
		return nil
	}
}

func bucketOrientedFitToItemOrderComparatorFunc(intValue uint) bucketorientedfit.ItemOrderComparatorFunc {
	switch intValue {
	case 0:
		return bucketorientedfit.AscendingItemSize
	case 1:
		return bucketorientedfit.DescendingItemSize
	default:
		return nil
	}
}
