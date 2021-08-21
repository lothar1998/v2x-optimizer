package console

import (
	"bufio"
	"github.com/lothar1998/v2x-optimizer/pkg/optimizer"
	"io"
	"strconv"
	"strings"
)

const (
	nValue            = "N"
	vValue            = "V"
	rrhCount          = "RRH_COUNT"
	rrhEnable         = "RRH_ENABLE"
	vehicleAssignment = "VEHICLE_ASSIGNMENT"
)

// FromOutput parses result in console format.
// The optimizer.Result can be encoded into console format using ToOutput.
func FromOutput(output string) (*optimizer.Result, error) {
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
		case strings.HasPrefix(line, rrhEnable):
			values := toValues(value)
			result.RRHEnable = make([]bool, len(values))
			var currVal bool
			for i, v := range values {
				if v == "1" {
					currVal = true
				} else {
					currVal = false
				}
				result.RRHEnable[i] = currVal
			}
		case strings.HasPrefix(line, vehicleAssignment):
			values := toValues(value)
			result.VehiclesToRRHAssignment = make([]int, len(values))
			for i, v := range values {
				assignment, err := strconv.Atoi(v)
				if err != nil {
					return nil, err
				}
				result.VehiclesToRRHAssignment[i] = assignment
			}
		default:
		}
	}
}

// ToOutput encodes optimizer.Result into a user-readable string. It can be parsed using FromOutput.
func ToOutput(result *optimizer.Result) string {
	var sb strings.Builder

	sb.WriteString(nValue + " = " + strconv.Itoa(len(result.RRHEnable)) + "\n")
	sb.WriteString(vValue + " = " + strconv.Itoa(len(result.VehiclesToRRHAssignment)) + "\n")
	sb.WriteString(rrhCount + " = " + strconv.Itoa(result.RRHCount) + "\n")

	sb.WriteString(rrhEnable + " = [")
	if len(result.RRHEnable) > 0 {
		if result.RRHEnable[0] {
			sb.WriteRune('1')
		} else {
			sb.WriteRune('0')
		}

		for _, e := range result.RRHEnable[1:] {
			sb.WriteRune(' ')
			if e {
				sb.WriteRune('1')
			} else {
				sb.WriteRune('0')
			}
		}
	}
	sb.WriteString("]\n")

	sb.WriteString(vehicleAssignment + " = [")
	if len(result.VehiclesToRRHAssignment) > 0 {
		sb.WriteRune(rune(result.VehiclesToRRHAssignment[0] + '0'))

		for _, e := range result.VehiclesToRRHAssignment[1:] {
			sb.WriteRune(' ')
			sb.WriteRune(rune(e + '0'))
		}
	}
	sb.WriteString("]\n")

	return sb.String()
}

func toValues(value string) []string {
	value = value[1 : len(value)-1]
	return strings.Split(value, " ")
}
