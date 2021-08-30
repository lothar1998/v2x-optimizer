package executor

import (
	"context"
	"github.com/lothar1998/v2x-optimizer/internal/namer"
)

type Executor interface {
	Execute(ctx context.Context) (int, error)
	namer.Namer
}
