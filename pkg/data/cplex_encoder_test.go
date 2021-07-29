package data

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

func TestCPLEXEncoder_Encode_Decode_Compatibility(t *testing.T) {
	t.Parallel()

	mbr := []int{1, 2, 3, 4, 5}
	r := [][]int{
		{11, 12, 13, 14, 15},
		{21, 22, 23, 24, 25},
		{31, 32, 33, 34, 35},
		{41, 42, 43, 44, 45},
	}
	data := &Data{MBR: mbr, R: r}

	encoder := CPLEXEncoder{}

	var buffer bytes.Buffer

	err := encoder.Encode(data, &buffer)

	assert.NoError(t, err)

	decodedData, err := encoder.Decode(&buffer)
	assert.NoError(t, err)

	assert.Equal(t, data, decodedData)
}

func TestCPLEXEncoder_Decode(t *testing.T) {
	t.Parallel()

	t.Run("should decode data from CPLEX format", func(t *testing.T) {
		t.Parallel()

		cplexStr := "V = 4;\n" +
			"N = 5;\n" +
			"MBR = [1 2 3 4 5];\n" +
			"R = [\n" +
			"[11 12 13 14 15]\n" +
			"[21 22 23 24 25]\n" +
			"[31 32 33 34 35]\n" +
			"[41 42 43 44 45]\n" +
			"];\n"

		mbr := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		data := &Data{MBR: mbr, R: r}

		decodeData, err := CPLEXEncoder{}.Decode(strings.NewReader(cplexStr))

		assert.NoError(t, err)
		assert.Equal(t, data, decodeData)
	})

	t.Run("should skip unknown variables", func(t *testing.T) {
		t.Parallel()

		cplexStr := "V = 4;\n" +
			"N = 5;\n" +
			"C = \"abcd\";\n" +
			"MBR = [1 2 3 4 5];\n" +
			"R = [\n" +
			"[11 12 13 14 15]\n" +
			"[21 22 23 24 25]\n" +
			"[31 32 33 34 35]\n" +
			"[41 42 43 44 45]\n" +
			"];\n"

		mbr := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		data := &Data{MBR: mbr, R: r}

		decodeData, err := CPLEXEncoder{}.Decode(strings.NewReader(cplexStr))

		assert.NoError(t, err)
		assert.Equal(t, data, decodeData)
	})

	t.Run("should not decode incorrect data", func(t *testing.T) {
		t.Parallel()

		cplexStr := "V = 4;\n" +
			"N = 5;\n" +
			"MBR = [1 2 3 4 5];\n" +
			"R = [\n" +
			"[11 12 13 14 15]\n" +
			"[21 22 24 25]\n" +
			"];\n"

		decodeData, err := CPLEXEncoder{}.Decode(strings.NewReader(cplexStr))

		assert.ErrorIs(t, err, ErrMalformedData)
		assert.Zero(t, decodeData)
	})

	t.Run("should not decode if one of size variables is skipped", func(t *testing.T) {
		t.Parallel()

		cplexStr := "V = 4;\n" +
			"MBR = [1 2 3 4 5];\n" +
			"R = [\n" +
			"[11 12 13 14 15]\n" +
			"[21 22 23 24 25]\n" +
			"];\n"

		decodeData, err := CPLEXEncoder{}.Decode(strings.NewReader(cplexStr))

		assert.ErrorIs(t, err, ErrMalformedData)
		assert.Zero(t, decodeData)
	})

	t.Run("should return error if data is malformed", func(t *testing.T) {
		t.Parallel()

		jsonString := `{
			  "MBR": [1, 2, 3, 4, 5],
			  "R": [
				[11, 22, 33, 44, 55],
				[11, 22, 33, 44, 55],
				[11, 22, 33, 44, 55],
				[11, 22, 33, 44, 55],
				[11, 22, 33, 44, 55]
			  ]
             }`

		decodeData, err := CPLEXEncoder{}.Decode(strings.NewReader(jsonString))

		assert.ErrorIs(t, err, ErrMalformedData)
		assert.Zero(t, decodeData)
	})
}

func TestCPLEXEncoder_Encode(t *testing.T) {
	t.Parallel()

	t.Run("should encode data to CPLEX format", func(t *testing.T) {
		t.Parallel()

		mbr := []int{1, 2, 3, 4, 5}
		r := [][]int{
			{11, 12, 13, 14, 15},
			{21, 22, 23, 24, 25},
			{31, 32, 33, 34, 35},
			{41, 42, 43, 44, 45},
		}
		data := &Data{MBR: mbr, R: r}

		expectedEncodedString := "V = 4;\n" +
			"N = 5;\n" +
			"MBR = [1 2 3 4 5];\n" +
			"R = [\n" +
			"[11 12 13 14 15]\n" +
			"[21 22 23 24 25]\n" +
			"[31 32 33 34 35]\n" +
			"[41 42 43 44 45]\n" +
			"];\n"

		var buffer bytes.Buffer

		err := CPLEXEncoder{}.Encode(data, &buffer)

		assert.NoError(t, err)
		assert.Equal(t, expectedEncodedString, buffer.String())
	})
}

func Test_toIntArray(t *testing.T) {
	t.Parallel()

	t.Run("should convert slice to int list", func(t *testing.T) {
		t.Parallel()

		slice := []int{1, 2, 3, 4, 5}
		expectedList := "[1 2 3 4 5]"

		list := toIntArray(slice)

		assert.Equal(t, expectedList, list)
	})

	t.Run("should convert empty slice to empty int list", func(t *testing.T) {
		t.Parallel()

		var slice []int
		expectedList := "[]"

		list := toIntArray(slice)

		assert.Equal(t, expectedList, list)
	})

	t.Run("should convert single item slice to single item int list", func(t *testing.T) {
		t.Parallel()

		slice := []int{1}
		expectedList := "[1]"

		list := toIntArray(slice)

		assert.Equal(t, expectedList, list)
	})
}

func Test_parseIntArray(t *testing.T) {
	t.Parallel()

	t.Run("should decode string array to int array", func(t *testing.T) {
		t.Parallel()

		expectedList := []int{2, 3, 5, 1, 5, 1, 22, 12123, 12}

		str1 := "[2 3 5 1 5 1 22 12123 12]"
		str2 := "[2 3 5 1 5  1 22  12123 12]"
		str3 := "[   2 3 5 1 5  1 22  12123   12   ]"

		list1, err1 := parseIntArray(str1)
		list2, err2 := parseIntArray(str2)
		list3, err3 := parseIntArray(str3)

		assert.NoError(t, err1)
		assert.Equal(t, expectedList, list1)

		assert.NoError(t, err2)
		assert.Equal(t, expectedList, list2)

		assert.NoError(t, err3)
		assert.Equal(t, expectedList, list3)
	})

	t.Run("should decode one item array to int array", func(t *testing.T) {
		t.Parallel()

		expectedList := []int{22}

		str1 := "[22]"
		str2 := "[    22   ]"

		list1, err1 := parseIntArray(str1)
		list2, err2 := parseIntArray(str2)

		assert.NoError(t, err1)
		assert.Equal(t, expectedList, list1)

		assert.NoError(t, err2)
		assert.Equal(t, expectedList, list2)
	})

	t.Run("should decode empty array to empty int array", func(t *testing.T) {
		t.Parallel()

		var expectedList []int

		str1 := "[]"
		str2 := "[       ]"

		list1, err1 := parseIntArray(str1)
		list2, err2 := parseIntArray(str2)

		assert.NoError(t, err1)
		assert.Equal(t, expectedList, list1)

		assert.NoError(t, err2)
		assert.Equal(t, expectedList, list2)
	})
}

func Test_parseArrayOfArrays(t *testing.T) {
	t.Parallel()

	t.Run("should parse array of arrays", func(t *testing.T) {
		t.Parallel()

		expectedArrays := []string{"[2 3 5 1 5 1 22 12123 12]", "[8 7 5 3 8 1 99 21 37]"}

		str := "R = [\n[2 3 5 1 5 1 22 12123 12]\n[8 7 5 3 8 1 99 21 37]\n]\n"

		assert.Equal(t, expectedArrays, parseArrayOfArrays(str))
	})

	t.Run("should parse array of arrays with whitespace chars", func(t *testing.T) {
		t.Parallel()

		expectedArrays := []string{"[2 3   5    1 5 1 22 12123 12]", "[8 7 5 3  8 1 99 21 37]"}

		str := "\t\n R \t \n =  \n  [ \t \n[2 3\n   5 \t   1 5 1 22 12123 12] \n   [8 7 5 3\n  8 1 99 21 37]   \n]  \n"

		assert.Equal(t, expectedArrays, parseArrayOfArrays(str))
	})
}

func Test_parseInt(t *testing.T) {
	t.Parallel()

	t.Run("should parse int variable", func(t *testing.T) {
		t.Parallel()

		str := "N = 32"

		variable, err := parseInt(str)

		assert.NoError(t, err)
		assert.Equal(t, 32, variable)
	})

	t.Run("should return error if variable is not an integer", func(t *testing.T) {
		t.Parallel()

		var expectedError *strconv.NumError

		str := "N = 32.3"

		variable, err := parseInt(str)

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, variable)
	})
}

func Test_findValue(t *testing.T) {
	t.Parallel()

	t.Run("should parse single value", func(t *testing.T) {
		t.Parallel()

		str := "N = 32"

		assert.Equal(t, "32", findValue(str))
	})

	t.Run("should parse single value with whitespace chars", func(t *testing.T) {
		t.Parallel()

		str := " \t N   \n    = \n \t  32   \t  \n"

		assert.Equal(t, "32", findValue(str))
	})

	t.Run("should parse array value", func(t *testing.T) {
		t.Parallel()

		str := "MBR = [2 3 5 1 5 1 22 12123 12]"

		assert.Equal(t, "[2 3 5 1 5 1 22 12123 12]", findValue(str))
	})

	t.Run("should parse array value with whitespace chars", func(t *testing.T) {
		t.Parallel()

		str := " \t MBR   \n    = \t  [2 3 5 1 5\n 1   22\t 12123 12]    \t   \n   "

		assert.Equal(t, "[2 3 5 1 5 1   22 12123 12]", findValue(str))
	})

	t.Run("should parse multiline array", func(t *testing.T) {
		t.Parallel()

		str := "R = [\n[2 3 5 1 5 1 22 12123 12]\n[8 7 5 3 8 1 99 21 37]\n]\n"

		assert.Equal(t, "[[2 3 5 1 5 1 22 12123 12][8 7 5 3 8 1 99 21 37]]", findValue(str))
	})

	t.Run("should parse multiline array with whitespace chars", func(t *testing.T) {
		t.Parallel()

		str := "\t  \n  R \t \n =  \n  [ \t \n[2 3\n   5 \t   1 5 1 22 12123 12] \n [8 7 5 3\n  8 1 99 21 37]   \n]  \n"

		assert.Equal(t, "[  [2 3   5    1 5 1 22 12123 12]  [8 7 5 3  8 1 99 21 37]   ]", findValue(str))
	})
}
