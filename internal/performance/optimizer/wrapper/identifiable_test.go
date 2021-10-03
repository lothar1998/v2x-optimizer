package wrapper

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
)

type valueReceiver struct {
	A string `id_include:"true"`
	B int
	C bool `id_include:"true"`
}

func (o valueReceiver) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type pointerReceiver struct {
	A string `id_include:"true"`
	B int
}

func (o *pointerReceiver) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type empty struct {
	A bool
}

func (e empty) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type overriddenNames struct {
	A bool `id_include:"true" id_name:"A_renamed"`
	B int  `id_include:"true"`
}

func (o overriddenNames) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

func TestOptimizer_Identifier(t *testing.T) {
	t.Parallel()

	t.Run("should make an identifier from tagged values of struct", func(t *testing.T) {
		t.Parallel()

		o := Identifiable{Optimizer: valueReceiver{
			A: "a",
			B: 2,
			C: true,
		}}

		assert.Equal(t,
			"valueReceiver,A:a,C:true",
			o.Identifier())
	})

	t.Run("should make an identifier from exported values of pointer to struct", func(t *testing.T) {
		t.Parallel()

		o := Identifiable{Optimizer: &pointerReceiver{A: "a", B: 12}}

		assert.Equal(t, "pointerReceiver,A:a", o.Identifier())
	})

	t.Run("should make an identifier consisted of only struct name", func(t *testing.T) {
		t.Parallel()

		o := Identifiable{empty{}}

		assert.Equal(t, "empty", o.Identifier())
	})

	t.Run("should override name of field", func(t *testing.T) {
		t.Parallel()

		o := Identifiable{overriddenNames{A: false, B: 32}}

		assert.Equal(t, "overriddenNames,A_renamed:false,B:32", o.Identifier())
	})
}
