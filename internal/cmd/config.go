package cmd

import (
	"errors"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizers"
)

type formatEncoderInfo struct {
	FormatDisplayName string
	Encoder           data.EncoderDecoder
}

const (
	jsonFormat  = "json"
	plainFormat = "plain"
	cplexFormat = "cplex"
)

var (
	formatsToEncodersInfo = map[string]formatEncoderInfo{
		jsonFormat:  {"json", data.JSONEncoder{}},
		plainFormat: {"plain (CSV-like)", data.PlainEncoder{}},
		cplexFormat: {"CPLEX", data.CPLEXEncoder{}},
	}

	namesToOptimizers = map[string]optimizers.Optimizer{}
)

var (
	errCannotOpenFile   = errors.New("cannot open file")
	errCannotParseData  = errors.New("cannot parse data")
	errCannotEncodeData = errors.New("cannot encode data")

	errIncorrectDataType = errors.New("incorrect provided data type")
)
