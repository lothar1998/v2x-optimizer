package utils

// TODO add support for VehiclesToRRHAssignment

import (
	"bufio"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"io"
	"strconv"
	"strings"
)

const (
	rrhCount = "RRH_COUNT"
	rrh      = "RRH"
)

func FromConsoleOutput(output string) (*optimizer.Result, error) {
	result := &optimizer.Result{}

	reader := bufio.NewReader(strings.NewReader(output))

	for {
		line, err := reader.ReadString('\n')

		if err == io.EOF {
			return result, nil
		} else if err != nil {
			return nil, err
		}

		line = line[:len(line)-1]
		value := line[strings.Index(line, "=")+1:]
		value = strings.TrimSpace(value)

		switch {
		case strings.HasPrefix(line, rrhCount):
			parsedInt, err := strconv.ParseInt(value, 10, 32)
			if err != nil {
				return nil, err
			}
			result.RRHCount = int(parsedInt)
		case strings.HasPrefix(line, rrh):
			value = value[1 : len(value)-1]
			values := strings.Split(value, " ")

			var currVal bool
			for _, v := range values {
				if v == "1" {
					currVal = true
				} else {
					currVal = false
				}
				result.RRH = append(result.RRH, currVal)
			}
		default:
		}
	}
}

func ToConsoleOutput(result *optimizer.Result) string {
	var sb strings.Builder

	sb.WriteString(rrhCount + " = " + strconv.Itoa(result.RRHCount) + "\n")
	sb.WriteString(rrh + " = [")

	if len(result.RRH) > 0 {
		if result.RRH[0] {
			sb.WriteRune('1')
		} else {
			sb.WriteRune('0')
		}

		for _, e := range result.RRH[1:] {
			sb.WriteRune(' ')
			if e {
				sb.WriteRune('1')
			} else {
				sb.WriteRune('0')
			}
		}
	}

	sb.WriteString("]\n")
	return sb.String()
}
