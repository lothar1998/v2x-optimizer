package wrapper

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
)

type Cacheable struct {
	optimizer.Optimizer
}

func (c Cacheable) Optimize(ctx context.Context, data *data.Data) (*optimizer.Result, error) {
	return c.Optimizer.Optimize(ctx, data)
}

func (c Cacheable) Name() string {
	return c.Optimizer.Name()
}

func (c Cacheable) IsCacheable() bool {
	return true
}
