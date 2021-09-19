package optimizerfactory

import (
	"github.com/lothar1998/v2x-optimizer/internal/identifiable"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/spf13/cobra"
)

type BuildFunc func(*cobra.Command) (optimizer.IdentifiableOptimizer, error)

type Factory interface {
	identifiable.Identifiable
	Builder() BuildFunc
	SetUpFlags(*cobra.Command)
}
