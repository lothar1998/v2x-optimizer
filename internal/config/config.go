package config

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

// RegisteredOptimizers is a list of all possible optimizers.
var RegisteredOptimizers = []optimizer.Optimizer{optimizer.FirstFit{}}
