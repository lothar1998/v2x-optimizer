package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateNormal(t *testing.T) {
	t.Parallel()

	v := 10
	n := 5

	data := GenerateNormal(v, n)

	assert.Len(t, data.MRB, n)
	assert.Len(t, data.R, v)
	assert.Len(t, data.R[0], n)
	assert.GreaterOrEqual(t, sum(data.MRB), len(data.MRB))
	for _, s := range data.R {
		assert.GreaterOrEqual(t, sum(s), len(s))
	}
}
