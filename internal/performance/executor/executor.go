package executor

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/behavior"
)

type Executor interface {
	behavior.Identifiable
	behavior.Cacheable
	Execute(ctx context.Context) (int, error)
}
