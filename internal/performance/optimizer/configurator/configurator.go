package configurator

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/spf13/cobra"
)

type BuildFunc func(*cobra.Command) (optimizer.IdentifiableCacheableOptimizer, error)

type Configurator interface {
	TypeName() string
	Builder() BuildFunc
	SetUpFlags(*cobra.Command)
}
