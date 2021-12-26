package encoder

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// CPLEXEncoder facilitates encoding Data to CPLEX data format.
type CPLEXEncoder struct{}

// Encode allows for encoding Data to CPLEX data format.
func (e CPLEXEncoder) Encode(input *data.Data, writer io.Writer) error {
	lengths := fmt.Sprintf("V = %d;\nN = %d;\n", len(input.R), len(input.MRB))
	_, err := writer.Write([]byte(lengths))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("MRB = " + toIntArray(input.MRB) + ";\n"))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("R = [\n"))
	if err != nil {
		return err
	}

	for _, elems := range input.R {
		_, err = writer.Write([]byte(toIntArray(elems) + "\n"))
		if err != nil {
			return err
		}
	}

	_, err = writer.Write([]byte("];\n"))
	if err != nil {
		return err
	}

	return nil
}

// Decode allows for decoding CPLEX data format to Data.
// It returns an error if the sizes of R and MRB are not equal to size variables [V, N].
// It is possible to decode data with additional variables defined. In such a case Decode skips these values.
func (e CPLEXEncoder) Decode(reader io.Reader) (*data.Data, error) {
	var output data.Data
	bufferedReader := bufio.NewReader(reader)
	var n, v int

	for {
		line, err := bufferedReader.ReadString(';')
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		line = line[:len(line)-1]
		line = strings.TrimLeft(line, " \t\n")

		switch {
		case strings.HasPrefix(line, "N"):
			integer, err := parseInt(line)
			if err != nil {
				return nil, err
			}
			n = integer

		case strings.HasPrefix(line, "V"):
			integer, err := parseInt(line)
			if err != nil {
				return nil, err
			}
			v = integer

		case strings.HasPrefix(line, "R"):
			arrays := parseArrayOfArrays(line)
			var r [][]int
			for _, array := range arrays {
				intArray, err := parseIntArray(array)
				if err != nil {
					return nil, err
				}
				r = append(r, intArray)
			}
			output.R = r

		case strings.HasPrefix(line, "MRB"):
			variable := findValue(line)
			mrb, err := parseIntArray(variable)
			if err != nil {
				return nil, err
			}
			output.MRB = mrb

		default:
		}
	}

	if err := verifyDecodedData(&output, n, v); err != nil {
		return nil, err
	}

	return &output, nil
}

func toIntArray(elems []int) string {
	switch len(elems) {
	case 0:
		return "[]"
	case 1:
		return "[" + strconv.Itoa(elems[0]) + "]"
	}

	var sb strings.Builder

	sb.WriteRune('[')
	sb.WriteString(strconv.Itoa(elems[0]))

	for _, s := range elems[1:] {
		sb.WriteRune(' ')
		sb.WriteString(strconv.Itoa(s))
	}

	sb.WriteRune(']')

	return sb.String()
}

func parseIntArray(array string) (result []int, err error) {
	array = array[1 : len(array)-1]

	leftIndex := 0
	length := 0
	isPreviousSpace := true

	for i, c := range array {
		switch {
		case unicode.IsDigit(c):
			length++
			isPreviousSpace = false
		case unicode.IsSpace(c) && !isPreviousSpace:
			parsedInt, err := strconv.ParseInt(array[leftIndex:leftIndex+length], 10, 32)
			if err != nil {
				return nil, err
			}

			result = append(result, int(parsedInt))

			length = 0
			leftIndex = i + 1
			isPreviousSpace = true
		case unicode.IsSpace(c) && isPreviousSpace:
			leftIndex = i + 1
		}
	}

	if leftIndex < len(array) {
		parsedInt, err := strconv.ParseInt(array[leftIndex:leftIndex+length], 10, 32)
		if err != nil {
			return nil, err
		}

		result = append(result, int(parsedInt))
	}

	return result, nil
}

func parseArrayOfArrays(str string) []string {
	var arrays []string

	str = findValue(str)
	str = str[1 : len(str)-1]

	startIndex := 0
	for i, c := range str {
		if c == '[' {
			startIndex = i
		} else if c == ']' {
			arrays = append(arrays, str[startIndex:i+1])
		}
	}

	return arrays
}

func parseInt(str string) (int, error) {
	value := findValue(str)
	parsedInt, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(parsedInt), nil
}

func findValue(str string) string {
	index := strings.Index(str, "=")
	str = str[index+1:]
	str = strings.Trim(str, " \n\t")
	str = strings.ReplaceAll(str, "\t", "")
	return strings.ReplaceAll(str, "\n", "")
}

func verifyDecodedData(output *data.Data, n, v int) error {
	if output.MRB == nil || output.R == nil {
		return data.ErrMalformedData
	}

	if len(output.MRB) != n && len(output.R) != v {
		return data.ErrMalformedData
	}

	for _, e := range output.R {
		if len(e) != n {
			return data.ErrMalformedData
		}
	}

	return nil
}
