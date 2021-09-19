package optimizer

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
)

type optimizerWithParams struct {
	A string
	B int
	C bool
	D [][]int
	E map[string]map[int]bool
	f string
}

func (o optimizerWithParams) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type optimizerPointerReceiver struct {
	A string
	B int
	c bool
}

func (o *optimizerPointerReceiver) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type empty struct{}

func (e empty) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

func TestOptimizer_Identifier(t *testing.T) {
	t.Parallel()

	t.Run("should make an identifier from exported values of struct", func(t *testing.T) {
		t.Parallel()

		o := IdentifiableWrapper{Optimizer: optimizerWithParams{
			A: "a",
			B: 2,
			C: true,
			D: [][]int{{12}, {213, 133}},
			E: map[string]map[int]bool{"abc": {2: true}},
			f: "ASD",
		}}

		assert.Equal(t,
			"optimizerWithParams,A:a,B:2,C:true,D:[[12] [213 133]],E:map[abc:map[2:true]]",
			o.Identifier())
	})

	t.Run("should make an identifier from exported values of pointer to struct", func(t *testing.T) {
		t.Parallel()

		o := IdentifiableWrapper{Optimizer: &optimizerPointerReceiver{A: "a", B: 12, c: true}}

		assert.Equal(t, "optimizerPointerReceiver,A:a,B:12", o.Identifier())
	})

	t.Run("should make an identifier consisted of only struct name", func(t *testing.T) {
		t.Parallel()

		o := IdentifiableWrapper{empty{}}

		assert.Equal(t, "empty", o.Identifier())
	})
}
