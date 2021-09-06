package cmd

import (
	"bytes"
	"strings"
	"testing"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/errors"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"github.com/stretchr/testify/assert"
)

// func Test_outputToFile(t *testing.T) {
//	t.Parallel()
//
//	results := map[string]*pathsToErrors{
//		"/path1/subpath1/subsubpath1": {
//			AverageRelativeError: 0.5,
//			PathToErrors: map[string]*calculator.ErrorInfo{
//				"/example1": {
//					CustomResult:  5,
//					CPLEXResult:   2,
//					AbsoluteError: 1,
//					RelativeError: 0.5,
//				},
//			},
//		},
//		"path2/subpath2/": {
//			AverageRelativeError: 1.5,
//			PathToErrors: map[string]*calculator.ErrorInfo{
//				"/example2": {
//					CustomResult:  10,
//					CPLEXResult:   4,
//					AbsoluteError: 6,
//					RelativeError: 1.5,
//				},
//			},
//		},
//	}
//
//	t.Run("should save result to CSV file", func(t *testing.T) {
//		t.Parallel()
//
//		dir, err := ioutil.TempDir("", "v2x-optimizer-performance-csv-output-*")
//		assert.NoError(t, err)
//
//		expectedFile1 := path.Join(dir, "path1_subpath1_subsubpath1.csv")
//		expectedFile2 := path.Join(dir, "path2_subpath2.csv")
//
//		err = outputToCSVFile(results, dir)
//
//		assert.NoError(t, err)
//		assert.FileExists(t, expectedFile1)
//		assert.FileExists(t, expectedFile2)
//
//		file1, err := ioutil.ReadFile(expectedFile1)
//		assert.NoError(t, err)
//		lines1 := strings.Split(string(file1), "\n")
//		assert.Equal(t, strings.Join(csvHeaders, ","), lines1[0])
//		assert.Equal(t, "/example1,5,2,1,0.500,0.500", lines1[1])
//
//		file2, err := ioutil.ReadFile(expectedFile2)
//		assert.NoError(t, err)
//		lines2 := strings.Split(string(file2), "\n")
//		assert.Equal(t, strings.Join(csvHeaders, ","), lines2[0])
//		assert.Equal(t, "/example2,10,4,6,1.500,1.500", lines2[1])
//	})
//}
//
// func Test_toSeparatedValues(t *testing.T) {
//	t.Parallel()
//
//	t.Run("should transform pathsToErrors into separated values", func(t *testing.T) {
//		t.Parallel()
//
//		results := &pathsToErrors{
//			AverageRelativeError: 2.25,
//			PathToErrors: map[string]*calculator.ErrorInfo{
//				"path1": {
//					CustomResult:  5,
//					CPLEXResult:   2,
//					AbsoluteError: 1,
//					RelativeError: 0.5,
//				},
//				"path2": {
//					CustomResult:  10,
//					CPLEXResult:   2,
//					AbsoluteError: 8,
//					RelativeError: 4,
//				},
//			}}
//
//		separatedValues := toSeparatedValues(results)
//
//		assert.Len(t, separatedValues, 2)
//		assert.Contains(t, separatedValues, []string{"path1", "5", "2", "1", "0.500", "2.250"})
//		assert.Contains(t, separatedValues, []string{"path2", "10", "2", "8", "4.000", "2.250"})
//	})
//}

func Test_toErrors(t *testing.T) {
	t.Parallel()

	t.Run("should transform results into errors", func(t *testing.T) {
		t.Parallel()

		results := runner.PathsToResults{
			"/path/1": runner.FilesToResults{
				"file1": runner.OptimizersToResults{config.CPLEXOptimizerName: 2, "opt1": 3, "opt4": 13},
				"file2": runner.OptimizersToResults{config.CPLEXOptimizerName: 4, "opt1": 5, "opt4": 12},
				"file3": runner.OptimizersToResults{config.CPLEXOptimizerName: 5, "opt1": 12, "opt4": 32},
			},
			"/path/2": runner.FilesToResults{
				"file4": runner.OptimizersToResults{config.CPLEXOptimizerName: 12, "opt1": 23, "opt4": 53},
				"file5": runner.OptimizersToResults{config.CPLEXOptimizerName: 14, "opt1": 35, "opt4": 22},
			},
		}

		errs := toErrors(results)

		assert.Len(t, errs, 2)
		assert.Contains(t, errs, "/path/1")
		assert.Contains(t, errs, "/path/2")

		assert.Len(t, errs["/path/1"], 3)
		assert.Contains(t, errs["/path/1"], "file1")
		assert.Contains(t, errs["/path/1"], "file2")
		assert.Contains(t, errs["/path/1"], "file3")

		assert.Len(t, errs["/path/2"], 2)
		assert.Contains(t, errs["/path/2"], "file4")
		assert.Contains(t, errs["/path/2"], "file5")

		assert.Len(t, errs["/path/1"]["file1"], 2)
		assert.Len(t, errs["/path/1"]["file2"], 2)
		assert.Len(t, errs["/path/1"]["file3"], 2)

		assert.Len(t, errs["/path/2"]["file4"], 2)
		assert.Len(t, errs["/path/2"]["file5"], 2)
	})
}

func Test_toAverageErrors(t *testing.T) {
	t.Parallel()

	t.Run("should compute average errors per path", func(t *testing.T) {
		t.Parallel()

		pathsToErrors := PathsToErrors{
			"/path/1": FilesToErrors{
				"file1": OptimizersToErrors{
					"opt1": errors.Info{RelativeError: 3.0, AbsoluteError: 6},
					"opt2": errors.Info{RelativeError: 4.5, AbsoluteError: 8}},
				"file2": OptimizersToErrors{
					"opt1": errors.Info{RelativeError: 12.0, AbsoluteError: 1},
					"opt2": errors.Info{RelativeError: 5.5, AbsoluteError: 2}},
				"file3": OptimizersToErrors{
					"opt1": errors.Info{RelativeError: 3.0, AbsoluteError: 2},
					"opt2": errors.Info{RelativeError: 2.3, AbsoluteError: 2}},
			},
			"/path/2": FilesToErrors{
				"file4": OptimizersToErrors{
					"opt1": errors.Info{RelativeError: 3.49, AbsoluteError: 1},
					"opt2": errors.Info{RelativeError: 4.0, AbsoluteError: 5}},
				"file5": OptimizersToErrors{
					"opt1": errors.Info{RelativeError: 4.51, AbsoluteError: 4},
					"opt2": errors.Info{RelativeError: 4.0, AbsoluteError: 5}},
			},
		}

		averageErrors := toAverageErrors(pathsToErrors)

		assert.Len(t, averageErrors, 2)
		assert.Contains(t, averageErrors, "/path/1")
		assert.Contains(t, averageErrors, "/path/2")

		assert.Len(t, averageErrors["/path/1"], 2)
		assert.Len(t, averageErrors["/path/2"], 2)

		assert.Contains(t, averageErrors["/path/1"], "opt1")
		assert.Contains(t, averageErrors["/path/1"], "opt2")
		assert.Contains(t, averageErrors["/path/2"], "opt1")
		assert.Contains(t, averageErrors["/path/2"], "opt2")

		assertErrorWithinDelta(t, AvgErrors{6.0, 3}, averageErrors["/path/1"]["opt1"])
		assertErrorWithinDelta(t, AvgErrors{4.1, 4}, averageErrors["/path/1"]["opt2"])
		assertErrorWithinDelta(t, AvgErrors{4.0, 2.5}, averageErrors["/path/2"]["opt1"])
		assertErrorWithinDelta(t, AvgErrors{4.0, 5}, averageErrors["/path/2"]["opt2"])
	})
}

func Test_pathToUnderscoreValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			"should transform path to underscore filename",
			"example/path/to/file",
			"example_path_to_file",
		},
		{
			"should transform path to underscore filename with slash at the beginning",
			"/example/path/to/file",
			"example_path_to_file",
		},
		{
			"should transform path to underscore filename with slash at the end",
			"example/path/to/file/",
			"example_path_to_file",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := pathToUnderscoreValue(tt.path); got != tt.want {
				t.Errorf("pathToUnderscoreValue() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_writeAvgErrors(t *testing.T) {
	t.Parallel()

	expectedHeader := "optimizer,average absolute error,average relative error"

	optToAvgErr := OptimizersToAvgErrors{
		"opt1": AvgErrors{1, 3},
		"opt2": AvgErrors{3.3, 2.5},
	}

	var buffer bytes.Buffer

	err := writeAvgErrors(optToAvgErr, &buffer)
	assert.NoError(t, err)

	s := buffer.String()
	lines := strings.Split(s, "\n")

	assert.Len(t, lines, 4)
	assert.Equal(t, expectedHeader, lines[0])
	assert.Contains(t, lines, "opt1,3.000,1.000")
	assert.Contains(t, lines, "opt2,2.500,3.300")
}

func Test_writeErrors(t *testing.T) {
	t.Parallel()

	expectedHeader := "filename,optimizer,value,optimal value,absolute error,relative error"

	filesToErrs := FilesToErrors{
		"file1": OptimizersToErrors{
			"opt1": errors.Info{Value: 1, ReferenceValue: 2, AbsoluteError: 1, RelativeError: 0.5},
			"opt2": errors.Info{Value: 3, ReferenceValue: 6, AbsoluteError: 3, RelativeError: 0.5}},
		"file2": OptimizersToErrors{
			"opt1": errors.Info{Value: 5, ReferenceValue: 10, AbsoluteError: 5, RelativeError: 0.5},
			"opt2": errors.Info{Value: 12, ReferenceValue: 6, AbsoluteError: 6, RelativeError: 1},
		},
	}

	var buffer bytes.Buffer

	err := writeErrors(filesToErrs, &buffer)
	assert.NoError(t, err)

	s := buffer.String()
	lines := strings.Split(s, "\n")

	assert.Len(t, lines, 6)
	assert.Equal(t, expectedHeader, lines[0])
	assert.Contains(t, lines, "file1,opt1,1,2,1,0.500")
	assert.Contains(t, lines, "file1,opt2,3,6,3,0.500")
	assert.Contains(t, lines, "file2,opt1,5,10,5,0.500")
	assert.Contains(t, lines, "file2,opt2,12,6,6,1.000")
}

func assertErrorWithinDelta(t *testing.T, expected, given AvgErrors) {
	assert.InDelta(t, expected.AvgRelativeError, given.AvgRelativeError, 0.1)
	assert.InDelta(t, expected.AvgAbsolutError, given.AvgAbsolutError, 0.1)
}
