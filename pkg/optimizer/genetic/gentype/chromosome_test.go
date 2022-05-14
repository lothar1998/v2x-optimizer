package gentype

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewChromosome(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(5)

	assert.Len(t, chromosome.buckets, 5)
	assert.NotNil(t, chromosome.idsToBuckets)
}

func TestChromosome_At(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(2)
	chromosome.buckets[1] = &Bucket{}

	assert.NotNil(t, chromosome.At(1))
}

func TestChromosome_Slice(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(3)
	chromosome.buckets[0] = &Bucket{id: 1}
	chromosome.buckets[1] = &Bucket{id: 2}
	chromosome.buckets[2] = &Bucket{id: 3}

	slice := chromosome.Slice(1, 3)
	assert.Equal(t, 2, slice[0].id)
	assert.Equal(t, 3, slice[1].id)
}

func TestChromosome_SetAt(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(2)

	chromosome.SetAt(1, &Bucket{id: 1})

	assert.Nil(t, chromosome.buckets[0])
	assert.NotNil(t, chromosome.buckets[1])
	assert.Equal(t, 1, chromosome.buckets[1].id)
	assert.Len(t, chromosome.idsToBuckets, 1)
	assert.Contains(t, chromosome.idsToBuckets, 1)
}

func TestChromosome_Append(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(0)
	require.Len(t, chromosome.buckets, 0)
	require.Len(t, chromosome.idsToBuckets, 0)

	chromosome.Append(&Bucket{})

	assert.Len(t, chromosome.buckets, 1)
	assert.Len(t, chromosome.idsToBuckets, 1)
}

func TestChromosome_ContainsBucket(t *testing.T) {
	t.Parallel()

	bucket := &Bucket{id: 5}

	chromosome := NewChromosome(1)
	chromosome.buckets[0] = bucket
	chromosome.idsToBuckets[bucket.id] = bucket

	assert.True(t, chromosome.ContainsBucket(bucket.id))
	assert.False(t, chromosome.ContainsBucket(123))
}

func TestChromosome_Len(t *testing.T) {
	t.Parallel()

	chromosome := NewChromosome(10)
	assert.Equal(t, 10, chromosome.Len())
}
