package data

import (
	"bufio"
	"io"
	"strconv"
	"strings"
)

// DefaultDelimiter defines default delimiter expected in encoded data.
var DefaultDelimiter = ','

// PlainEncoder facilitates encoding/decoding Data into CSV-like format. For example:
// 1,2,3,4,5
// 6,7,8,9,2
// 9,4,2,1,0
// where the first line is MBR values and further lines are R values.
type PlainEncoder struct{}

// Encode allows for encoding data to writer into CSV-like format.
// It returns ErrMalformedData if the lengths of R slices are not equal to MBR slice length.
// It is possible to encode Data that consists only of MBR values.
func (e PlainEncoder) Encode(data *Data, writer io.Writer) error {
	if len(data.MBR) == 0 {
		return ErrMalformedData
	}

	_, err := writer.Write([]byte(joinInts(data.MBR, DefaultDelimiter) + "\n"))
	if err != nil {
		return err
	}

	for _, line := range data.R {
		_, err := writer.Write([]byte(joinInts(line, DefaultDelimiter) + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}

// Decode allows for decoding data from reader in CSV-like format.
// It returns an error if the lengths of R slices are not equal to MBR slice length. It is possible to decode
// data that consists only of MBR values. Input data line should be ended with newline '\n' character,
// however, it is possible to do not use '\n' in the last line.
func (e PlainEncoder) Decode(reader io.Reader) (*Data, error) {
	var data Data

	lineReader := bufio.NewReader(reader)

	line, err := lineReader.ReadString('\n')

	if err == io.EOF && len(line) > 0 {
		err := setMBR(&data, line)
		if err != nil {
			return nil, err
		}
		return &data, nil
	} else if err == io.EOF {
		return &data, nil
	} else if err != nil {
		return nil, err
	}

	err = setMBR(&data, line[:len(line)-1])
	if err != nil {
		return nil, err
	}

	mbrLen := len(data.MBR)

	for {
		line, err := lineReader.ReadString('\n')

		if err == io.EOF && len(line) > 0 {
			err := appendR(&data, line, mbrLen)
			if err != nil {
				return nil, err
			}
			return &data, nil
		} else if err == io.EOF {
			return &data, nil
		} else if err != nil {
			return nil, err
		}

		err = appendR(&data, line[:len(line)-1], mbrLen)
		if err != nil {
			return nil, err
		}
	}
}

func joinInts(elems []int, sep rune) string {
	switch len(elems) {
	case 0:
		return ""
	case 1:
		return strconv.Itoa(elems[0])
	}

	var b strings.Builder

	b.WriteString(strconv.Itoa(elems[0]))
	for _, s := range elems[1:] {
		b.WriteRune(sep)
		b.WriteString(strconv.Itoa(s))
	}
	return b.String()
}

func setMBR(data *Data, line string) error {
	mbr, err := splitIntString(line, DefaultDelimiter)
	if err != nil {
		return err
	}

	data.MBR = mbr

	return nil
}

func appendR(data *Data, line string, mbrLen int) error {
	rLine, err := splitIntString(line, DefaultDelimiter)
	if err != nil {
		return err
	}

	if len(rLine) != mbrLen {
		return ErrMalformedData
	}

	data.R = append(data.R, rLine)

	return nil
}

func splitIntString(str string, sep rune) (result []int, err error) {
	if len(str) == 0 {
		return
	}

	startIndex := 0

	for i := 0; i < len(str); i++ {
		if rune(str[i]) != sep {
			continue
		}

		parsedInt, err := strconv.ParseInt(str[startIndex:i], 10, 32)
		if err != nil {
			return nil, err
		}

		result = append(result, int(parsedInt))
		startIndex = i + 1
	}

	if startIndex >= len(str) {
		return
	}

	parsedInt, err := strconv.ParseInt(str[startIndex:], 10, 32)
	if err != nil {
		return nil, err
	}

	result = append(result, int(parsedInt))

	return
}
