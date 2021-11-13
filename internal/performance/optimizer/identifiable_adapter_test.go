package optimizer

import (
	"context"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"github.com/stretchr/testify/assert"
)

type valueReceiver struct {
	Name string `id_name:""`
	A    string `id_include:"true"`
	B    int
	C    bool `id_include:"true"`
}

func (o valueReceiver) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type pointerReceiver struct {
	Name string `id_name:""`
	A    string `id_include:"true"`
	B    int
}

func (o *pointerReceiver) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type empty struct {
	Name string `id_name:""`
	A    bool
}

func (e empty) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type overriddenNames struct {
	Name string `id_name:""`
	A    bool   `id_include:"true" id_rename:"A_renamed"`
	B    int    `id_include:"true"`
}

func (o overriddenNames) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type defaultStructName struct {
	Name string `id_name:""`
	A    bool   `id_include:"true"`
}

func (d defaultStructName) Optimize(_ context.Context, _ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

type defaultStructNameWithoutNameFieldDeclared struct {
	A bool `id_include:"true"`
}

func (d defaultStructNameWithoutNameFieldDeclared) Optimize(_ context.Context,
	_ *data.Data) (*optimizer.Result, error) {
	return nil, nil
}

func TestOptimizer_Identifier(t *testing.T) {
	t.Parallel()

	t.Run("should make an identifier from struct name and tagged values of struct", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{Optimizer: valueReceiver{
			Name: "name1",
			A:    "a",
			B:    2,
			C:    true,
		}}

		assert.Equal(t, "name1,A:a,C:true", o.Identifier())
	})

	t.Run("should make an identifier from struct name and tagged values of pointer to struct", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{Optimizer: &pointerReceiver{Name: "name2", A: "a", B: 12}}

		assert.Equal(t, "name2,A:a", o.Identifier())
	})

	t.Run("should make an identifier consisted of only struct name", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{empty{Name: "name3"}}

		assert.Equal(t, "name3", o.Identifier())
	})

	t.Run("should override name of field", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{overriddenNames{Name: "name4", A: false, B: 32}}

		assert.Equal(t, "name4,A_renamed:false,B:32", o.Identifier())
	})

	t.Run("should use default name if name field is empty", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{defaultStructName{A: false}}

		assert.Equal(t, "defaultStructName,A:false", o.Identifier())
	})

	t.Run("should use default name if name field is not provided", func(t *testing.T) {
		t.Parallel()

		o := identifiableAdapter{defaultStructNameWithoutNameFieldDeclared{A: false}}

		assert.Equal(t, "defaultStructNameWithoutNameFieldDeclared,A:false", o.Identifier())
	})
}
