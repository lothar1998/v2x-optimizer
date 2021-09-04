package executor

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/namer"
)

// Executor represents an object that can be used to execute something, especially
// to compute some value and return it. It also should allow for error handling.
// Executor should also have a name that is returned by Name() function.
type Executor interface {
	Execute(ctx context.Context) (int, error)
	namer.Namer
}
