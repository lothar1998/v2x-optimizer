package data

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
	"unicode"
)

// CPLEXEncoder facilitates encoding Data to CPLEX data format.
type CPLEXEncoder struct{}

// Encode allows for encoding Data to CPLEX data format.
func (e CPLEXEncoder) Encode(data *Data, writer io.Writer) error {
	lengths := fmt.Sprintf("V = %d;\nN = %d;\n", len(data.R), len(data.MBR))
	_, err := writer.Write([]byte(lengths))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("MBR = " + toCPLEXArray(data.MBR) + ";\n"))
	if err != nil {
		return err
	}

	_, err = writer.Write([]byte("R = [\n"))
	if err != nil {
		return err
	}

	for _, elems := range data.R {
		_, err = writer.Write([]byte(toCPLEXArray(elems) + "\n"))
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
// It returns an error if the sizes of R and MBR are not equal to size variables [V, N].
// It is possible to decode data with additional variables defined. In such a case Decode skips these values.
func (e CPLEXEncoder) Decode(reader io.Reader) (*Data, error) {
	var data Data

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
		line = strings.Trim(line, " \t\n")

		switch {
		case strings.HasPrefix(line, "N"):
			variable, err := parseSingleIntVariable(line)
			if err != nil {
				return nil, err
			}
			n = variable

		case strings.HasPrefix(line, "V"):
			variable, err := parseSingleIntVariable(line)
			if err != nil {
				return nil, err
			}
			v = variable

		case strings.HasPrefix(line, "R"):
			arrays := parseArrayOfArrays(line)
			var r [][]int
			for _, array := range arrays {
				list, err := fromCPLEXArray(array)
				if err != nil {
					return nil, err
				}
				r = append(r, list)
			}
			data.R = r

		case strings.HasPrefix(line, "MBR"):
			variable := parseValue(line)
			mbr, err := fromCPLEXArray(variable)
			if err != nil {
				return nil, err
			}
			data.MBR = mbr

		default:
		}
	}

	if len(data.MBR) != n && len(data.R) != v {
		return nil, ErrMalformedData
	}

	for _, e := range data.R {
		if len(e) != n {
			return nil, ErrMalformedData
		}
	}

	return &data, nil
}

func toCPLEXArray(elems []int) string {
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

func fromCPLEXArray(list string) (result []int, err error) {
	list = list[1 : len(list)-1]
	list = strings.Trim(list, " ")

	leftIndex := 0
	length := 0
	isPreviousSpace := false

	for i, c := range list {
		if unicode.IsDigit(c) {
			length++
			isPreviousSpace = false
		} else if unicode.IsSpace(c) && !isPreviousSpace {
			parsedInt, err := strconv.ParseInt(list[leftIndex:leftIndex+length], 10, 32)
			if err != nil {
				return nil, err
			}

			result = append(result, int(parsedInt))

			length = 0
			leftIndex = i + 1
			isPreviousSpace = true
		} else if unicode.IsSpace(c) && isPreviousSpace {
			leftIndex = i + 1
		}
	}

	if leftIndex < len(list) {
		parsedInt, err := strconv.ParseInt(list[leftIndex:leftIndex+length], 10, 32)
		if err != nil {
			return nil, err
		}

		result = append(result, int(parsedInt))
	}

	return
}

func parseArrayOfArrays(str string) []string {
	var stringArrays []string

	str = parseValue(str)
	str = str[1 : len(str)-1]
	str = strings.Trim(str, " ")

	startIndex := 0
	for i, c := range str {
		if c == '[' {
			startIndex = i
		} else if c == ']' {
			stringArrays = append(stringArrays, str[startIndex:i+1])
		}
	}

	return stringArrays
}

func parseSingleIntVariable(str string) (int, error) {
	value := parseValue(str)
	parsedInt, err := strconv.ParseInt(value, 10, 32)
	if err != nil {
		return 0, err
	}
	return int(parsedInt), nil
}

func parseValue(str string) string {
	index := strings.Index(str, "=")
	str = str[index+1:]
	str = strings.Trim(str, " \n\t")
	str = strings.ReplaceAll(str, "\t", "")
	return strings.ReplaceAll(str, "\n", "")
}
