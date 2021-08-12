package cmd

import "github.com/lothar1998/v2x-optimizer/cmd/v2x-optimizer-performance/calculator"

type resultForPath struct {
	AverageApproxError float64
	PathToApproxErrors map[string]*calculator.ApproxErrorInfo
}

func newEmptyResultForPath() *resultForPath {
	return &resultForPath{PathToApproxErrors: make(map[string]*calculator.ApproxErrorInfo)}
}
