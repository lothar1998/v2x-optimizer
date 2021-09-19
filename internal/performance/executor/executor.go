package executor

import (
	"context"

	"github.com/lothar1998/v2x-optimizer/internal/identifiable"
)

// Executor represents an object that can be used to execute something, especially
// to compute some value and return it. It also should allow for error handling.
// Executor should also have a name that is returned by Name() function.
type Executor interface {
	identifiable.Identifiable
	Execute(ctx context.Context) (int, error)
}
