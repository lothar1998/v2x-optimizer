package data

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDecode(t *testing.T) {
	t.Parallel()

	mbr := []int{1, 2, 3, 4, 5}
	r := [][]int{
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
	}
	data := &Data{MBR: mbr, R: r}

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

	decodedData, err := JSONEncoder{}.Decode(strings.NewReader(jsonString))

	assert.NoError(t, err)
	assert.Equal(t, data, decodedData)
}

func TestEncode(t *testing.T) {
	t.Parallel()

	mbr := []int{1, 2, 3, 4, 5}
	r := [][]int{
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
		{11, 22, 33, 44, 55},
	}
	data := &Data{MBR: mbr, R: r}

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

	var buffer bytes.Buffer

	err := JSONEncoder{}.Encode(data, &buffer)

	assert.NoError(t, err)
	assert.Equal(t, trimWhiteSigns(jsonString), trimWhiteSigns(buffer.String()))
}

func trimWhiteSigns(str string) string {
	spaceTrimmed := strings.ReplaceAll(str, " ", "")
	newLineTrimmed := strings.ReplaceAll(spaceTrimmed, "\n", "")
	tabTrimmed := strings.ReplaceAll(newLineTrimmed, "\t", "")
	return tabTrimmed
}
