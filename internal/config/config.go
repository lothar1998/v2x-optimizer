package config

import (
	wrapper2 "github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	wrapper "github.com/lothar1998/v2x-optimizer/internal/performance/wrapper"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

// CPLEXOptimizerName is a name of cplex optimizer. It needs to be defined here
// because CPLEX doesn't have optimizer.Wrapper implementation.
const CPLEXOptimizerName = "cplex"

// RegisteredOptimizers is a list of all possible optimizers.
var RegisteredOptimizers = []wrapper2.Wrapper{
	wrapper.Cacheable{Optimizer: optimizer.FirstFit{}},
	wrapper.Cacheable{Optimizer: optimizer.NextFit{}},
}
