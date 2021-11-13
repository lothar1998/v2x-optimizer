package configurator

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/nextkfit"
	"github.com/spf13/cobra"
)

const (
	nextKFitParameterK = "nkf_k"
	nextKFitName       = "NextKFit"
)

type NextKFitWrapper struct {
	Name string `id_name:""`
	K    int    `id_include:"true"`
	nextkfit.NextKFit
}

type NextKFitConfigurator struct{}

func (n NextKFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (optimizer.IdentifiableCacheableOptimizer, error) {
		k, err := command.Flags().GetUint(nextKFitParameterK)
		if err != nil {
			return nil, err
		}

		if k == 0 {
			return nil, errors.New("parameter K cannot be equal to 0")
		}

		nkf := &NextKFitWrapper{
			Name:     nextKFitName,
			K:        int(k),
			NextKFit: nextkfit.NextKFit{K: int(k)},
		}

		return &optimizer.CacheableAdapter{
			IsCacheEligible:       true,
			IdentifiableOptimizer: &optimizer.IdentifiableAdapter{Optimizer: nkf},
		}, nil
	}
}

func (n NextKFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(nextKFitParameterK, "", 1, "NextKFit k parameter")
}

func (n NextKFitConfigurator) TypeName() string {
	return nextKFitName
}
