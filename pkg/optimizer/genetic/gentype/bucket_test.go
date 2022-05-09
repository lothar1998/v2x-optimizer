package gentype

import (
	"errors"
	"reflect"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBucket(t *testing.T) {
	t.Parallel()

	bucketID := 1
	capacity := 3

	bucket := NewBucket(bucketID, capacity)

	assert.Equal(t, bucketID, bucket.id)
	assert.Equal(t, capacity, bucket.capacity)
	assert.Equal(t, capacity, bucket.freeSpace)
	assert.NotNil(t, bucket.itemMapping)
}

func TestBucket_AddItem(t *testing.T) {
	t.Parallel()

	t.Run("should item and update internals", func(t *testing.T) {
		t.Parallel()

		bucket := NewBucket(1, 5)

		err := bucket.AddItem(&Item{id: 10, size: 4})

		assert.NoError(t, err)
		assert.Equal(t, 5, bucket.capacity)
		assert.Equal(t, 1, bucket.freeSpace)
		assert.Contains(t, bucket.itemMapping, 10)
	})

	t.Run("should return error since item is already in bucket", func(t *testing.T) {
		t.Parallel()

		bucket := NewBucket(1, 5)
		err := bucket.AddItem(&Item{id: 10, size: 1})
		require.NoError(t, err)

		err = bucket.AddItem(&Item{id: 10, size: 4})

		assert.EqualError(t, errors.Unwrap(err), ErrItemAlreadyInBucket.Error())
		assert.Equal(t, 5, bucket.capacity)
		assert.Equal(t, 4, bucket.freeSpace)
		assert.Len(t, bucket.itemMapping, 1)
		assert.Contains(t, bucket.itemMapping, 10)
	})

	t.Run("should return error since item exceeds free space in bucket", func(t *testing.T) {
		t.Parallel()

		bucket := NewBucket(1, 5)
		err := bucket.AddItem(&Item{id: 10, size: 1})
		require.NoError(t, err)

		err = bucket.AddItem(&Item{id: 20, size: 5})

		assert.EqualError(t, errors.Unwrap(err), ErrBucketCapacityLimit.Error())
		assert.Equal(t, 5, bucket.capacity)
		assert.Equal(t, 4, bucket.freeSpace)
		assert.Len(t, bucket.itemMapping, 1)
		assert.Contains(t, bucket.itemMapping, 10)
	})
}

func TestBucket_Capacity(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 10)
	capacity := bucket.Capacity()
	assert.Equal(t, 10, capacity)
}

func TestBucket_DeepCopy(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 19)
	err := bucket.AddItem(&Item{id: 10, size: 3})
	require.NoError(t, err)

	bucketCopy := bucket.DeepCopy()

	assert.Equal(t, bucket, bucketCopy)
	assert.Equal(t, bucket.itemMapping, bucketCopy.itemMapping)
	assertReflectNotEqual(t, bucket.itemMapping, bucketCopy.itemMapping)
}

func TestBucket_FreeSpace(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 10)
	err := bucket.AddItem(&Item{10, 4})
	require.NoError(t, err)

	freeSpace := bucket.FreeSpace()

	assert.Equal(t, 6, freeSpace)
}

func TestBucket_ID(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 10)
	bucketID := bucket.ID()
	assert.Equal(t, 1, bucketID)
}

func TestBucket_IsEmpty(t *testing.T) {
	t.Parallel()

	t.Run("should return that bucket is empty", func(t *testing.T) {
		t.Parallel()

		bucket := NewBucket(1, 10)
		assert.True(t, bucket.IsEmpty())
	})

	t.Run("should return that bucket is not empty", func(t *testing.T) {
		t.Parallel()

		bucket := NewBucket(1, 10)
		err := bucket.AddItem(&Item{id: 10, size: 2})
		require.NoError(t, err)

		assert.False(t, bucket.IsEmpty())
	})
}

func TestBucket_Map(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 5)
	err := bucket.AddItem(&Item{id: 10, size: 1})
	require.NoError(t, err)

	assert.Equal(t, bucket.itemMapping, bucket.Map())
	assertReflectEqual(t, bucket.itemMapping, bucket.Map())
}

func TestBucket_copyMapping(t *testing.T) {
	t.Parallel()

	bucket := NewBucket(1, 19)
	err := bucket.AddItem(&Item{id: 10, size: 3})
	require.NoError(t, err)

	mapping := bucket.itemMapping
	mappingCopy := bucket.copyMapping()

	assert.Equal(t, mapping, mappingCopy)
	assertReflectNotEqual(t, mapping, mappingCopy)
}

func TestNewBucketFactory(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{}
	bucketFactory := NewBucketFactory(inputData)
	assert.Same(t, inputData, bucketFactory.data)
}

func TestBucketFactory_CreateBucket(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{1},
		R: [][]int{
			{2},
			{3},
		},
	}
	bucketFactory := NewBucketFactory(inputData)

	bucket := bucketFactory.CreateBucket(0)

	assert.Equal(t, 0, bucket.id)
	assert.Equal(t, 1, bucket.capacity)
	assert.Equal(t, 1, bucket.freeSpace)
	assert.NotNil(t, bucket.itemMapping)
}

func TestBucketFactory_MaxID(t *testing.T) {
	t.Parallel()

	inputData := &data.Data{
		MRB: []int{1, 2, 3, 4, 5, 6, 7},
	}
	bucketFactory := NewBucketFactory(inputData)

	assert.Equal(t, 6, bucketFactory.MaxID())
}

func assertReflectEqual(t *testing.T, v1, v2 any) {
	assert.Equal(t, reflect.ValueOf(v1).Pointer(), reflect.ValueOf(v2).Pointer())
}

func assertReflectNotEqual(t *testing.T, v1, v2 any) {
	assert.NotEqual(t, reflect.ValueOf(v1).Pointer(), reflect.ValueOf(v2).Pointer())
}
