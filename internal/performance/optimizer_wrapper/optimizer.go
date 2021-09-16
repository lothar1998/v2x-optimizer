package wrapper

import (
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type Optimizer interface {
	optimizer.Optimizer
	IsCacheable() bool
}
