package cmd

import "github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"

type pathsToErrors struct {
	AverageRelativeError float64
	PathToErrors         map[string]*calculator.ErrorInfo
}

func newEmptyPathsToErrors() *pathsToErrors {
	return &pathsToErrors{PathToErrors: make(map[string]*calculator.ErrorInfo)}
}

type pathErrorInfoChannelPair struct {
	Path             string
	ErrorInfoChannel chan *calculator.ErrorInfo
}

type pathPathsToErrorsChannelPair struct {
	Path                 string
	PathsToErrorsChannel chan *pathsToErrors
}
