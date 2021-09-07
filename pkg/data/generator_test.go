package data

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	t.Parallel()

	v := 10
	n := 5

	data := Generate(v, n)

	assert.Len(t, data.MRB, n)
	assert.Len(t, data.R, v)
	assert.Len(t, data.R[0], n)
	assert.Greater(t, sum(data.MRB), 0)
	for _, s := range data.R {
		assert.Greater(t, sum(s), 0)
	}
}

func sum(slice []int) int {
	var result int

	for _, e := range slice {
		result += e
	}

	return result
}
