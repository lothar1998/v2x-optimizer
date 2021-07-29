package data

import (
	"errors"
	"github.com/lothar1998/resource-optimization-in-v2x-networks/pkg/data"
)

const (
	emptyStringFlag = ""
	emptyIntFlag    = 0
)

type formatEncoderInfo struct {
	FormatDisplayName string
	Encoder           data.EncoderDecoder
}

var (
	formatsToEncodersInfo = map[string]formatEncoderInfo{
		"json":  {"json", data.JSONEncoder{}},
		"plain": {"plain (CSV-like)", data.PlainEncoder{}},
		"cplex": {"CPLEX", data.CPLEXEncoder{}},
	}
)

var (
	errCannotOpenFile   = errors.New("cannot open file")
	errCannotParseData  = errors.New("cannot parse data")
	errCannotEncodeData = errors.New("cannot encode data")
)
