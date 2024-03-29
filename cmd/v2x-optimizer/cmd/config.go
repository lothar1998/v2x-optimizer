package cmd

import (
	"errors"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
	"github.com/lothar1998/v2x-optimizer/pkg/data/encoder"
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
		jsonFormat:  {"json", encoder.JSON{}},
		plainFormat: {"plain (CSV-like)", encoder.Plain{}},
		cplexFormat: {"CPLEX", encoder.CPLEX{}},
	}

	availableFileFormats = getAvailableFileFormats()
)

var (
	errCannotCreatePath  = errors.New("cannot create path")
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
