package wrapper

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type NotCacheable struct {
	optimizer.Optimizer
}

func (c NotCacheable) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	return c.Optimizer.Optimize(ctx, data)
}

func (c NotCacheable) Name() string {
	return c.Optimizer.Name()
}

func (c NotCacheable) IsCacheable() bool {
	return false
}
