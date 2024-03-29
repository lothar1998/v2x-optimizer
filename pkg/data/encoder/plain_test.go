package encoder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/stretchr/testify/assert"
)

func TestPlainEncoder_Encode_Decode_Compatibility(t *testing.T) {
	t.Parallel()

	mrb := []int{1, 2, 3, 4, 5}
	r := [][]int{
		{11, 12, 13, 14, 15},
		{21, 22, 23, 24, 25},
		{31, 32, 33, 34, 35},
		{41, 42, 43, 44, 45},
	}
	expectedData := &data.Data{MRB: mrb, R: r}

	encoder := Plain{}

	var buffer bytes.Buffer

	err := encoder.Encode(expectedData, &buffer)
	assert.NoError(t, err)

	decodedData, err := encoder.Decode(&buffer)
	assert.NoError(t, err)

	assert.Equal(t, expectedData, decodedData)
}

func TestPlainEncoder_Encode(t *testing.T) {
	t.Parallel()

	t.Run("should encode data", func(t *testing.T) {
		t.Parallel()

		expectedEncodedData := "1,2,3,4,5\n" +
			"11,12,13,14,15\n" +
			"21,22,23,24,25\n" +
			"31,32,33,34,35\n" +
			"41,42,43,44,45\n"

		mrb := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		expectedData := &data.Data{MRB: mrb, R: r}

		var buffer bytes.Buffer

		err := Plain{}.Encode(expectedData, &buffer)

		assert.NoError(t, err)
		assert.Equal(t, expectedEncodedData, buffer.String())
	})

	t.Run("should encode only MRB", func(t *testing.T) {
		t.Parallel()

		expectedData := &data.Data{MRB: []int{1, 2, 3}}

		var buffer bytes.Buffer

		err := Plain{}.Encode(expectedData, &buffer)

		assert.NoError(t, err)
		assert.Equal(t, "1,2,3\n", buffer.String())
	})

	t.Run("should not encode only R", func(t *testing.T) {
		t.Parallel()

		expectedData := &data.Data{R: [][]int{{1, 2, 3}}}

		var buffer bytes.Buffer

		err := Plain{}.Encode(expectedData, &buffer)

		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Equal(t, 0, buffer.Len())
	})
}

func TestPlainEncoder_Decode(t *testing.T) {
	t.Parallel()

	t.Run("should decode data", func(t *testing.T) {
		t.Parallel()

		dataString := "1,2,3,4,5\n" +
			"11,12,13,14,15\n" +
			"21,22,23,24,25\n" +
			"31,32,33,34,35\n" +
			"41,42,43,44,45\n"

		mrb := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		expectedData := &data.Data{MRB: mrb, R: r}

		decodedData, err := Plain{}.Decode(strings.NewReader(dataString))

		assert.NoError(t, err)
		assert.Equal(t, expectedData, decodedData)
	})

	t.Run("should decode data without new line at the end", func(t *testing.T) {
		t.Parallel()

		dataString := "1,2,3,4,5\n" +
			"11,12,13,14,15\n" +
			"21,22,23,24,25\n" +
			"31,32,33,34,35\n" +
			"41,42,43,44,45"

		mrb := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		expectedData := &data.Data{MRB: mrb, R: r}

		decodedData, err := Plain{}.Decode(strings.NewReader(dataString))

		assert.NoError(t, err)
		assert.Equal(t, expectedData, decodedData)
	})

	t.Run("should not decode incorrect data", func(t *testing.T) {
		t.Parallel()

		dataString := "1,2,3,4,5\n" +
			"11,12,13,14,15\n" +
			"21,22,23,24,25,27\n"

		decodedData, err := Plain{}.Decode(strings.NewReader(dataString))

		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Zero(t, decodedData)
	})

	t.Run("should decode only one line that represents MRB", func(t *testing.T) {
		t.Parallel()

		dataString := "1,2,3,4,5\n"

		expectedData := &data.Data{MRB: []int{1, 2, 3, 4, 5}}

		decodedData, err := Plain{}.Decode(strings.NewReader(dataString))

		assert.NoError(t, err)
		assert.Equal(t, expectedData, decodedData)
	})

	t.Run("should decode only one line that represents MRB without new line at the end", func(t *testing.T) {
		t.Parallel()

		dataString := "1,2,3,4,5"

		expectedData := &data.Data{MRB: []int{1, 2, 3, 4, 5}}

		decodedData, err := Plain{}.Decode(strings.NewReader(dataString))

		assert.NoError(t, err)
		assert.Equal(t, expectedData, decodedData)
	})

	t.Run("should decode empty string to empty data", func(t *testing.T) {
		t.Parallel()

		decodedData, err := Plain{}.Decode(strings.NewReader(""))

		assert.NoError(t, err)
		assert.Equal(t, &data.Data{}, decodedData)
	})

	t.Run("should decode empty string to empty data with new line at the end", func(t *testing.T) {
		t.Parallel()

		decodedData, err := Plain{}.Decode(strings.NewReader("\n"))

		assert.NoError(t, err)
		assert.Equal(t, &data.Data{}, decodedData)
	})

	t.Run("should not decode malformed data", func(t *testing.T) {
		t.Parallel()

		decodedData, err := Plain{}.Decode(strings.NewReader("11,a,2\n"))

		assert.Error(t, err)
		assert.Zero(t, decodedData)
	})
}

func Test_joinInts(t *testing.T) {
	t.Parallel()

	t.Run("should join ints into string delimited with given delimiter", func(t *testing.T) {
		t.Parallel()

		slice := []int{1, 2, 3, 4, 5}
		expectedString := "1,2,3,4,5"

		result := joinInts(slice, ',')

		assert.Equal(t, expectedString, result)
	})

	t.Run("should return empty string for empty slice", func(t *testing.T) {
		t.Parallel()

		assert.Empty(t, joinInts([]int{}, ','))
	})

	t.Run("should return single element for one element slice", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, "5", joinInts([]int{5}, ','))
	})
}

func Test_splitIntString(t *testing.T) {
	t.Parallel()

	t.Run("should split string into int slice using given delimiter", func(t *testing.T) {
		t.Parallel()

		str := "1,2,3,4,5"
		expectedSlice := []int{1, 2, 3, 4, 5}

		slice, err := splitIntString(str, ',')

		assert.NoError(t, err)
		assert.Equal(t, expectedSlice, slice)
	})

	t.Run("should split string into int slice using given delimiter - string ended with delimiter", func(t *testing.T) {
		t.Parallel()

		str := "1,2,3,4,5,"
		expectedSlice := []int{1, 2, 3, 4, 5}

		slice, err := splitIntString(str, ',')

		assert.NoError(t, err)
		assert.Equal(t, expectedSlice, slice)
	})

	t.Run("should return empty slice if string is empty", func(t *testing.T) {
		t.Parallel()

		slice, err := splitIntString("", ',')

		assert.NoError(t, err)
		assert.Empty(t, slice)
	})

	t.Run("should return error if string is malformed", func(t *testing.T) {
		t.Parallel()

		slice, err := splitIntString("1,2,a,b,3", ',')

		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Zero(t, slice)
	})

	t.Run("should return error if string has different delimiter than expected one", func(t *testing.T) {
		t.Parallel()

		slice, err := splitIntString("1,2,3,4,5", '.')

		assert.ErrorIs(t, err, data.ErrMalformedData)
		assert.Zero(t, slice)
	})
}
