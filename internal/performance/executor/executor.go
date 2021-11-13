package executor

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/identifiable"
	"github.com/lothar1998/v2x-optimizer/internal/performance/optimizer"
)

type Executor interface {
	identifiable.Identifiable
	optimizer.Cacheable
	Execute(ctx context.Context) (int, error)
}
