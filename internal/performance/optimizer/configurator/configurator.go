package configurator

import (
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer/wrapper"
	"github.com/spf13/cobra"
)

type BuildFunc func(*cobra.Command) (wrapper.IdentifiableOptimizer, error)

type Configurator interface {
	TypeName() string
	Builder() BuildFunc
	SetUpFlags(*cobra.Command)
}
