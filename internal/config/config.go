package config

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/configurator"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Optimizer implementation.
const CPLEXOptimizerName = "CPLEX"

// RegisteredOptimizerFactories is a list of all possible factories.
var RegisteredOptimizerFactories = []configurator.Configurator{
	configurator.NewParameterless(optimizer.FirstFit{}),
	configurator.NewParameterless(optimizer.NextFit{}),
	configurator.NewParameterless(optimizer.WorstFit{}),
	configurator.NextKFitConfigurator{},
	configurator.BestFitConfigurator{},
}
