package cmd

import (
	"errors"
	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

type encoderInfo struct {
	FormatDisplayName string
	Encoder           data.EncoderDecoder
}

const (
	jsonFormat  = "json"
	plainFormat = "plain"
	cplexFormat = "cplex"
)

var (
	formatsToEncodersInfo = map[string]encoderInfo{
		jsonFormat:  {"json", data.JSONEncoder{}},
		plainFormat: {"plain (CSV-like)", data.PlainEncoder{}},
		cplexFormat: {"CPLEX", data.CPLEXEncoder{}},
	}

	availableFileFormats = getAvailableFileFormats()
)

var (
	errParseInt          = errors.New("cannot parse integer")
	errCannotOpenFile    = errors.New("cannot open file")
	errCannotParseData   = errors.New("cannot parse data")
	errCannotEncodeData  = errors.New("cannot encode data")
	errUnknownDataFormat = errors.New("unknown data format")
)

func getAvailableFileFormats() []string {
	formatList := make([]string, len(formatsToEncodersInfo))

	var i int
	for format := range formatsToEncodersInfo {
		formatList[i] = format
		i++
	}

	return formatList
}