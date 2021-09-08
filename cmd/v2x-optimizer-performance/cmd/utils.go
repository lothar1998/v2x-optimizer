package cmd

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/lothar1998/v2x-optimizer/internal/config"
	"github.com/lothar1998/v2x-optimizer/internal/performance/errors"
	"github.com/lothar1998/v2x-optimizer/internal/performance/runner"
)

type PathsToErrors map[string]FilesToErrors

type FilesToErrors map[string]OptimizersToErrors

type OptimizersToErrors map[string]errors.Info

type PathsToAvgErrors map[string]OptimizersToAvgErrors

type OptimizersToAvgErrors map[string]AvgErrors

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

	for path := range avgErrs {
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

func outputToCSVFile(errs PathsToErrors, avgErrs PathsToAvgErrors, outputFilepath string) error {
	for path, optimizersToAvgErrors := range avgErrs {
		err := os.MkdirAll(outputFilepath, 0755)
		if err != nil {
			return err
		}

		rootFilepath := filepath.Join(outputFilepath, pathToUnderscoreValue(path))

		csvFile, err := os.OpenFile(rootFilepath+".csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		err = writeAvgErrors(optimizersToAvgErrors, csvFile)
		_ = csvFile.Close()

		if err != nil {
			return err
		}

		csvFileDetails, err := os.OpenFile(rootFilepath+"_details.csv", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}

		err = writeErrors(errs[path], csvFileDetails)
		_ = csvFileDetails.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func pathToUnderscoreValue(path string) string {
	return strings.Trim(strings.Join(strings.Split(path, "/"), "_"), "_")
}

func writeAvgErrors(optimizersToAvgErrors OptimizersToAvgErrors, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"optimizer", "average absolute error", "average relative error"}

	if err := writer.Write(header); err != nil {
		return err
	}

	for optimizer, avgErr := range optimizersToAvgErrors {
		err := writer.Write([]string{
			optimizer,
			strconv.FormatFloat(avgErr.AvgAbsolutError, 'f', 3, 64),
			strconv.FormatFloat(avgErr.AvgRelativeError, 'f', 3, 64),
		})

		if err != nil {
			return err
		}
	}

	return nil
}

func writeErrors(filesToErrors FilesToErrors, w io.Writer) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	header := []string{"filename", "optimizer", "value", "optimal value", "absolute error", "relative error"}

	if err := writer.Write(header); err != nil {
		return err
	}

	for filename, optimizersToErrors := range filesToErrors {
		for optimizer, errorInfo := range optimizersToErrors {
			err := writer.Write([]string{
				filename,
				optimizer,
				strconv.Itoa(errorInfo.Value),
				strconv.Itoa(errorInfo.ReferenceValue),
				strconv.Itoa(errorInfo.AbsoluteError),
				strconv.FormatFloat(errorInfo.RelativeError, 'f', 3, 64),
			})

			if err != nil {
				return err
			}
		}
	}

	return nil
}
