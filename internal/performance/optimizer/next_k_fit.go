package optimizer

import (
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/spf13/cobra"
)

const (
	nextKFitParameterK = "nkf_k"
	nextKFitName       = "NextKFit"
)

type NextKFitWrapper struct {
	Name string `id_name:""`
	K    int    `id_include:"true"`
	optimizer.NextKFit
}

type NextKFitConfigurator struct{}

func (n NextKFitConfigurator) Builder() BuildFunc {
	return func(command *cobra.Command) (IdentifiableOptimizer, error) {
		k, err := command.Flags().GetUint(nextKFitParameterK)
		if err != nil {
			return nil, err
		}

		nkf := &NextKFitWrapper{
			Name:     nextKFitName,
			K:        int(k),
			NextKFit: optimizer.NextKFit{K: int(k)},
		}
		return &IdentifiableAdapter{Optimizer: nkf}, nil
	}
}

func (n NextKFitConfigurator) SetUpFlags(command *cobra.Command) {
	command.Flags().UintP(nextKFitParameterK, "", 1, "NextKFit k parameter")
}

func (n NextKFitConfigurator) TypeName() string {
	return nextKFitName
}
