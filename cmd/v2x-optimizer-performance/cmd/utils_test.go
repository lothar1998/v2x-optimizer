package cmd

import (
	"github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path"
	"strings"
	"testing"
)

func Test_outputToFile(t *testing.T) {
	t.Parallel()

	results := map[string]*resultForPath{
		"/path1/subpath1/subsubpath1": {
			AverageRelativeError: 0.5,
			PathToRelativeErrors: map[string]*calculator.ErrorInfo{
				"/example1": {
					CustomResult:  5,
					CPLEXResult:   2,
					AbsoluteError: 1,
					RelativeError: 0.5,
				},
			},
		},
		"path2/subpath2/": {
			AverageRelativeError: 1.5,
			PathToRelativeErrors: map[string]*calculator.ErrorInfo{
				"/example2": {
					CustomResult:  10,
					CPLEXResult:   4,
					AbsoluteError: 6,
					RelativeError: 1.5,
				},
			},
		},
	}

	t.Run("should save result to CSV file", func(t *testing.T) {
		t.Parallel()

		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-csv-output-*")
		assert.NoError(t, err)

		expectedFile1 := path.Join(dir, "path1_subpath1_subsubpath1.csv")
		expectedFile2 := path.Join(dir, "path2_subpath2.csv")

		err = outputToCSVFile(results, dir)

		assert.NoError(t, err)
		assert.FileExists(t, expectedFile1)
		assert.FileExists(t, expectedFile2)

		file1, err := ioutil.ReadFile(expectedFile1)
		assert.NoError(t, err)
		lines1 := strings.Split(string(file1), "\n")
		assert.Equal(t, strings.Join(csvHeaders, ","), lines1[0])
		assert.Equal(t, "/example1,5,2,1,0.500,0.500", lines1[1])

		file2, err := ioutil.ReadFile(expectedFile2)
		assert.NoError(t, err)
		lines2 := strings.Split(string(file2), "\n")
		assert.Equal(t, strings.Join(csvHeaders, ","), lines2[0])
		assert.Equal(t, "/example2,10,4,6,1.500,1.500", lines2[1])
	})
}

func Test_toSeparatedValues(t *testing.T) {
	t.Parallel()

	t.Run("should transform resultForPath into separated values", func(t *testing.T) {
		t.Parallel()

		results := &resultForPath{
			AverageRelativeError: 2.25,
			PathToRelativeErrors: map[string]*calculator.ErrorInfo{
				"path1": {
					CustomResult:  5,
					CPLEXResult:   2,
					AbsoluteError: 1,
					RelativeError: 0.5,
				},
				"path2": {
					CustomResult:  10,
					CPLEXResult:   2,
					AbsoluteError: 8,
					RelativeError: 4,
				},
			}}

		separatedValues := toSeparatedValues(results)

		assert.Equal(t, []string{"path1", "5", "2", "1", "0.500", "2.250"}, separatedValues[0])
		assert.Equal(t, []string{"path2", "10", "2", "8", "4.000", "2.250"}, separatedValues[1])
	})
}
