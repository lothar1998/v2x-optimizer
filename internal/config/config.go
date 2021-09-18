package config

import (
	identifiable "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/optimizerfactory"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Optimizer implementation.
const CPLEXOptimizerName = "CPLEX"

// RegisteredFactories is a list of all possible factories.
var RegisteredFactories = []optimizerfactory.Factory{
	&optimizerfactory.Parameterless{Identifiable: &identifiable.IdentifiableOptimizer{Optimizer: optimizer.FirstFit{}}},
	&optimizerfactory.Parameterless{Identifiable: &identifiable.IdentifiableOptimizer{Optimizer: optimizer.NextFit{}}},
}
