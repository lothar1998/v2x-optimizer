package config

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

// CPLEXOptimizerName is a name of cplex executor and will be provided as cplex optimizer name.
const CPLEXOptimizerName = "cplex"

// RegisteredOptimizers is a list of all possible optimizers.
var RegisteredOptimizers = []optimizer.Optimizer{optimizer.FirstFit{}}
