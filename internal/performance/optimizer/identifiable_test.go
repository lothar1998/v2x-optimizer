package optimizer

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
)

type OptimizerWithParams struct {
	A string
	B int
	C bool
	D [][]int
	E map[string]map[int]bool
	f string
}

func (o OptimizerWithParams) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

func (o OptimizerWithParams) Name() string {
	return ""
}

type Empty struct{}

func (e Empty) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

func (e Empty) Name() string {
	return ""
}

func TestOptimizer_Identifier(t *testing.T) {
	t.Parallel()

	t.Run("should make an identifier from exported values of struct", func(t *testing.T) {
		t.Parallel()

		o := IdentifiableOptimizer{Optimizer: OptimizerWithParams{
			A: "a",
			B: 2,
			C: true,
			D: [][]int{{12}, {213, 133}},
			E: map[string]map[int]bool{"abc": {2: true}},
			f: "ASD",
		}}

		assert.Equal(t,
			"OptimizerWithParams,A:a,B:2,C:true,D:[[12] [213 133]],E:map[abc:map[2:true]]",
			o.Identifier())
	})

	t.Run("should make an identifier consisted of only struct name", func(t *testing.T) {
		o := IdentifiableOptimizer{Empty{}}

		assert.Equal(t, "Empty", o.Identifier())
	})
}
