package data

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/fs"
	"io/ioutil"
	"path"
	"testing"
)

const testFilepath = "../../test/data"

func TestLoad(t *testing.T) {
	t.Parallel()

	t.Run("should load data from filepath", func(t *testing.T) {
		t.Parallel()

		expectedMBR := []int{1, 2, 3, 4, 5}
		expectedR := [][]int{
			{11, 22, 33, 44, 55},
			{11, 22, 33, 44, 55},
			{11, 22, 33, 44, 55},
			{11, 22, 33, 44, 55},
			{11, 22, 33, 44, 55},
		}
		expectedData := &Data{MBR: expectedMBR, R: expectedR}

		data, err := Load(path.Join(testFilepath, "example.json"))

		assert.NoError(t, err)
		assert.Equal(t, expectedData, data)
	})

	t.Run("should return error due to not existing file", func(t *testing.T) {
		t.Parallel()

		var expectedError *fs.PathError

		data, err := Load("/wrong/filepath")

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, data)
	})

	t.Run("should return error due to malformed config file", func(t *testing.T) {
		t.Parallel()

		var expectedError *json.SyntaxError

		data, err := Load(path.Join(testFilepath, "example_malformed.json"))

		assert.ErrorAs(t, err, &expectedError)
		assert.Zero(t, data)
	})
}

func TestExport(t *testing.T) {
	t.Parallel()

	mbr := []int{1, 2, 3}
	r := [][]int{
		{11, 22, 33},
		{11, 22, 33},
		{11, 22, 33},
		{11, 22, 33},
		{11, 22, 33},
	}
	data := &Data{MBR: mbr, R: r}

	t.Run("should export data to file", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "resource-optimization-in-v2x-networks-loader-test-*")
		assert.NoError(t, err)

		filename := "exported_data.json"
		filepath := path.Join(dir, filename)

		err = Export(filepath, data)

		assert.NoError(t, err)
		assert.FileExists(t, filepath)
	})

	t.Run("should return error due to non-existing directory", func(t *testing.T) {
		t.Parallel()

		var expectedError *fs.PathError

		err := Export("/tmp/non-existing-dir/datafile.json", data)

		assert.ErrorAs(t, err, &expectedError)
	})
}
