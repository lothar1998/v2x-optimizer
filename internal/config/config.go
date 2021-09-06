package config

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Optimizer implementation
const CPLEXOptimizerName = "cplex"

// RegisteredOptimizers is a list of all possible optimizers.
var RegisteredOptimizers = []optimizer.Optimizer{optimizer.FirstFit{}}
