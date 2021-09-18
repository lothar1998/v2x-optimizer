package optimizerfactory

import (
	"github.com/lothar1998/v2x-optimizer/internal/namer"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
	"github.com/spf13/cobra"
)

type BuildFunc func(*cobra.Command) (optimizer.Identifiable, error)

type Factory interface {
	Builder() BuildFunc
	SetUpFlags(*cobra.Command)
	namer.Namer
}
