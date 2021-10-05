package optimizer

import (
	"github.com/spf13/cobra"
)

type BuildFunc func(*cobra.Command) (IdentifiableOptimizer, error)

type Configurator interface {
	TypeName() string
	Builder() BuildFunc
	SetUpFlags(*cobra.Command)
}
