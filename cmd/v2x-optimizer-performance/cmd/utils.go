package cmd

import (
	"fmt"
	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/errors"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
	"os"
	"strings"
	"text/tabwriter"
)

var csvHeaders = []string{"path", "custom_result", "cplex_result", "absolute_error", "relative_error", "average_relative_error"}

// PathsToErrors represents mapping between paths and errors mapping.
type PathsToErrors map[string]FilesToErrors

// FilesToErrors represents mapping between files and errors mapping.
type FilesToErrors map[string]OptimizersToErrors

// OptimizersToErrors represents mapping between optimizers' names and error infos.
type OptimizersToErrors map[string]errors.Info

// PathsToAvgErrors represents mapping between paths and average errors mapping.
type PathsToAvgErrors map[string]OptimizersToAvgErrors

// OptimizersToAvgErrors represents mapping between optimizers' names and average errors.
type OptimizersToAvgErrors map[string]AvgErrors

// AvgErrors consists of two average error - relative and absolute.
type AvgErrors struct {
	AvgRelativeError float64
	AvgAbsolutError  float64
}

func toErrors(results runner.PathsToResults) PathsToErrors {
	pathsToErrors := make(PathsToErrors)

	for path, filesToResults := range results {
		pathsToErrors[path] = make(FilesToErrors)

		for file, optimizersToResults := range filesToResults {
			pathsToErrors[path][file] = make(OptimizersToErrors)

			cplexValue := optimizersToResults[config.CPLEXOptimizerName]

			for opt, value := range optimizersToResults {
				if opt != config.CPLEXOptimizerName {
					pathsToErrors[path][file][opt] = *errors.Calculate(cplexValue, value)
				}
			}
		}
	}

	return pathsToErrors
}

func toAverageErrors(pathsToErrors PathsToErrors) PathsToAvgErrors {
	pathsToAvgErrors := make(PathsToAvgErrors)

	for path, filesToErrors := range pathsToErrors {
		optimizersToTotalRelativeError := make(map[string]float64)
		optimizersToTotalAbsolutError := make(map[string]int)
		optimizersToTotalCount := make(map[string]int)

		for _, optimizersToErrors := range filesToErrors {
			for opt, errorInfo := range optimizersToErrors {
				if _, ok := optimizersToTotalRelativeError[opt]; !ok {
					optimizersToTotalRelativeError[opt] = errorInfo.RelativeError
					optimizersToTotalAbsolutError[opt] = errorInfo.AbsoluteError
					optimizersToTotalCount[opt] = 1
					continue
				}

				optimizersToTotalRelativeError[opt] += errorInfo.RelativeError
				optimizersToTotalAbsolutError[opt] += errorInfo.AbsoluteError
				optimizersToTotalCount[opt]++
			}
		}

		pathsToAvgErrors[path] = make(OptimizersToAvgErrors)

		for opt := range optimizersToTotalRelativeError {
			totalCount := float64(optimizersToTotalCount[opt])
			pathsToAvgErrors[path][opt] = AvgErrors{
				AvgRelativeError: optimizersToTotalRelativeError[opt] / totalCount,
				AvgAbsolutError:  float64(optimizersToTotalAbsolutError[opt]) / totalCount,
			}
		}
	}

	return pathsToAvgErrors
}

func outputToConsole(errs PathsToErrors, avgErrs PathsToAvgErrors, isVerbose bool) {
	w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)

	for path := range errs {
		_, _ = fmt.Fprintf(w, "Path: "+path)
		_, _ = fmt.Fprint(w, "\n\n")

		_, _ = fmt.Fprintf(w, "\t%s\t%s\t%s\n",
			"Optimizer", "Average relative error", "Average absolute error")

		for opt, avgValues := range avgErrs[path] {
			_, _ = fmt.Fprintf(w, "\t%s\t%.3f\t%.3f\n",
				opt, avgValues.AvgRelativeError, avgValues.AvgAbsolutError)
		}

		if isVerbose {
			_, _ = fmt.Fprint(w, "\n\n")

			for file, optimizersToErrors := range errs[path] {
				_, _ = fmt.Fprintln(w, "\tFile: "+file)
				_, _ = fmt.Fprintf(w, "\t\t%s\t%s\t%s\t%s\t%s\n",
					"Optimizer", "Value", "Optimal Value", "Relative Error", "Absolute Error")

				for opt, v := range optimizersToErrors {
					_, _ = fmt.Fprintf(w, "\t\t%s\t%d\t%d\t%.3f\t%d\n",
						opt, v.Value, v.ReferenceValue, v.RelativeError, v.AbsoluteError)
				}

				_, _ = fmt.Fprint(w, "\n")
			}
		}

		_, _ = fmt.Fprint(w, "\n")
		_, _ = fmt.Fprintln(w, strings.Repeat("-", 100))
	}

	_ = w.Flush()
}

//
//func outputToCSVFile(pathsToResults map[string]*pathsToErrors, outputFilepath string) error {
//	for rootPath, result := range pathsToResults {
//		err := os.MkdirAll(outputFilepath, 0755)
//		if err != nil {
//			return err
//		}
//
//		csvFilepath := path.Join(outputFilepath,
//			strings.Trim(strings.ReplaceAll(rootPath, "/", "_"), "_")+".csv")
//
//		file, err := os.OpenFile(csvFilepath, os.O_CREATE|os.O_WRONLY, 0644)
//		if err != nil {
//			return err
//		}
//
//		writer := csv.NewWriter(file)
//
//		err = writer.Write(csvHeaders)
//		if err != nil {
//			return err
//		}
//
//		err = writer.WriteAll(toSeparatedValues(result))
//		if err != nil {
//			return err
//		}
//
//		_ = file.Close()
//	}
//
//	return nil
//}
//
//func toSeparatedValues(resultForPath *pathsToErrors) [][]string {
//	result := make([][]string, len(resultForPath.PathToErrors))
//
//	var i int
//	for currentPath, info := range resultForPath.PathToErrors {
//		result[i] = make([]string, 6)
//
//		result[i][0] = currentPath
//		result[i][1] = strconv.Itoa(info.CustomResult)
//		result[i][2] = strconv.Itoa(info.CPLEXResult)
//		result[i][3] = strconv.Itoa(info.AbsoluteError)
//		result[i][4] = strconv.FormatFloat(info.RelativeError, 'f', 3, 64)
//		result[i][5] = strconv.FormatFloat(resultForPath.AverageRelativeError, 'f', 3, 64)
//
//		i++
//	}
//
//	return result
//}
