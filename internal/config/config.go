package config

import (
	optimizerconfigurator "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Optimizer implementation.
const CPLEXOptimizerName = "CPLEX"

// RegisteredOptimizerConfigurators is a list of all possible configurators.
var RegisteredOptimizerConfigurators = []optimizerconfigurator.Configurator{
	optimizerconfigurator.NewParameterless(optimizer.FirstFit{}),
	optimizerconfigurator.NewParameterless(optimizer.NextFit{}),
	optimizerconfigurator.NewParameterless(optimizer.WorstFit{}),
	optimizerconfigurator.NextKFitConfigurator{},
	optimizerconfigurator.BestFitConfigurator{},
}
