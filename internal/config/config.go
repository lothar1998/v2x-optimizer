package config

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

// NamesToOptimizers is an assignment of names to optimizers.
// Should consist of all possible implemented optimizers.
var NamesToOptimizers = map[string]optimizer.Optimizer{
	"first-fit": optimizer.FirstFit{},
}
