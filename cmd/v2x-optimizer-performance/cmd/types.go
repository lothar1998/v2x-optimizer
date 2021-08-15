package cmd

import "github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"

type resultForPath struct {
	AverageRelativeError float64
	PathToRelativeErrors map[string]*calculator.ErrorInfo
}

func newEmptyResultForPath() *resultForPath {
	return &resultForPath{PathToRelativeErrors: make(map[string]*calculator.ErrorInfo)}
}
