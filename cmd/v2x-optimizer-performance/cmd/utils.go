package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"
)

var csvHeaders = []string{"path", "custom_result", "cplex_result", "absolute_error", "relative_error", "average_relative_error"}

func outputToConsole(pathsToResults map[string]*pathsToErrors) {
	for rootPath, result := range pathsToResults {
		titleString := "Root: " + rootPath
		fmt.Println(strings.Repeat("-", utf8.RuneCountInString(titleString)+10))
		fmt.Println(titleString)
		fmt.Printf("Average relative error: %.3f", result.AverageRelativeError)
		fmt.Println()
		for subPath, errorInfo := range result.PathToErrors {
			fmt.Printf("%s\t->\tCustomResult: %d\t\tCPLEXResult: %d\t\tAbsoluteError: %d\t\tRelativeError: %.3f\n",
				filepath.Base(subPath), errorInfo.CustomResult, errorInfo.CPLEXResult,
				errorInfo.AbsoluteError, errorInfo.RelativeError)
		}
	}
}

func outputToCSVFile(pathsToResults map[string]*pathsToErrors, outputFilepath string) error {
	for rootPath, result := range pathsToResults {
		err := os.MkdirAll(outputFilepath, 0755)
		if err != nil {
			return err
		}

		csvFilepath := path.Join(outputFilepath,
			strings.Trim(strings.ReplaceAll(rootPath, "/", "_"), "_")+".csv")

		file, err := os.OpenFile(csvFilepath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}

		writer := csv.NewWriter(file)

		err = writer.Write(csvHeaders)
		if err != nil {
			return err
		}

		err = writer.WriteAll(toSeparatedValues(result))
		if err != nil {
			return err
		}

		_ = file.Close()
	}

	return nil
}

func toSeparatedValues(resultForPath *pathsToErrors) [][]string {
	result := make([][]string, len(resultForPath.PathToErrors))

	var i int
	for currentPath, info := range resultForPath.PathToErrors {
		result[i] = make([]string, 6)

		result[i][0] = currentPath
		result[i][1] = strconv.Itoa(info.CustomResult)
		result[i][2] = strconv.Itoa(info.CPLEXResult)
		result[i][3] = strconv.Itoa(info.AbsoluteError)
		result[i][4] = strconv.FormatFloat(info.RelativeError, 'f', 3, 64)
		result[i][5] = strconv.FormatFloat(resultForPath.AverageRelativeError, 'f', 3, 64)

		i++
	}

	return result
}
