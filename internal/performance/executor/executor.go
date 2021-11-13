package executor

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/common"
)

type Executor interface {
	common.Identifiable
	common.Cacheable
	Execute(ctx context.Context) (int, error)
}
