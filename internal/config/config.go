package config

import "github.com/lothar1998/v2x-optimizer/pkg/optimizer"

var NamesToOptimizers = map[string]optimizer.Optimizer{
	"first-fit": optimizer.FirstFit{},
}
