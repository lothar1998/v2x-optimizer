package genetic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnePointCrossover(t *testing.T) {
	t.Parallel()

	p1 := []int{1, 2, 3, 4}
	p2 := []int{5, 6, 7, 8}

	c1, c2 := OnePointCrossover(p1, p2)

	assertMultiPointCrossover(t, p1, p2, c1, c2, 1)
}

func Test_joinSlices(t *testing.T) {
	t.Parallel()

	t.Run("should join slices", func(t *testing.T) {
		t.Parallel()

		s1 := []int{1, 2, 3, 4}
		s2 := []int{5, 6, 7}
		s3 := []int{8, 9, 10, 11, 12}

		expectedSlice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12}

		slice := joinSlices(s1, s2, s3)

		assert.Equal(t, expectedSlice, slice)
	})

	t.Run("shouldn't refer to original slice", func(t *testing.T) {
		t.Parallel()

		s1 := []int{1, 2, 3, 4}
		s2 := []int{5, 6, 7}

		slice := joinSlices(s1, s2)

		slice[0] = 100

		assert.Equal(t, []int{1, 2, 3, 4}, s1)
	})
}

func assertMultiPointCrossover(t *testing.T, p1, p2, c1, c2 Chromosome, crossPointCount int) {
	var straightMatchSet bool
	var straightMatch bool

	count := 0
	for i := range p1 {
		if c1[i] != p1[i] {
			assert.Equal(t, c2[i], p1[i])
			assert.Equal(t, c1[i], p2[i])
			if straightMatch == true {
				if straightMatchSet {
					count++
				}
				straightMatchSet = true
				straightMatch = false
			}
		} else {
			assert.Equal(t, c2[i], p2[i])
			if straightMatch == false {
				if straightMatchSet {
					count++
				}
				straightMatchSet = true
				straightMatch = true
			}
		}
	}

	assert.Equal(t, crossPointCount, count)
}
