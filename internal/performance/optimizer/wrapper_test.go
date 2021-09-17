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
	panic("implement me")
}

func (o OptimizerWithParams) Name() string {
	panic("implement me")
}

type Empty struct{}

func (e Empty) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	panic("implement me")
}

func (e Empty) Name() string {
	panic("implement me")
}

func TestOptimizer_MapKey(t *testing.T) {
	t.Parallel()

	t.Run("should make a key from exported values of struct", func(t *testing.T) {
		t.Parallel()

		o := Wrapper{Optimizer: OptimizerWithParams{
			A: "a",
			B: 2,
			C: true,
			D: [][]int{{12}, {213, 133}},
			E: map[string]map[int]bool{"abc": {2: true}},
			f: "ASD",
		}}

		assert.Equal(t,
			"OptimizerWithParams,A:a,B:2,C:true,D:[[12] [213 133]],E:map[abc:map[2:true]]",
			o.MapKey())
	})

	t.Run("should make key consisted of only struct name", func(t *testing.T) {
		o := Wrapper{Empty{}}

		assert.Equal(t, "Empty", o.MapKey())
	})
}
