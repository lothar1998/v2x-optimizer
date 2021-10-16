package config

import (
	optimizerConfigurator "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/configurator"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/almostworstfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/firstfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/nextfit"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer/worstfit"
)

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Optimizer implementation.
const CPLEXOptimizerName = "CPLEX"

// RegisteredOptimizerConfigurators is a list of all possible configurators.
var RegisteredOptimizerConfigurators = []optimizerConfigurator.Configurator{
	optimizerConfigurator.NewParameterless(firstfit.FirstFit{}),
	optimizerConfigurator.NewParameterless(nextfit.NextFit{}),
	optimizerConfigurator.NewParameterless(worstfit.WorstFit{}),
	optimizerConfigurator.NewParameterless(almostworstfit.AlmostWorstFit{}),
	optimizerConfigurator.NextKFitConfigurator{},
	optimizerConfigurator.BestFitConfigurator{},
}
