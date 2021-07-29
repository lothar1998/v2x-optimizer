package data

import (
	"errors"
	"io"
)

// Data is a structure that represents data required to the computation of the solution.
type Data struct {
	MRB []int
	R   [][]int
}

// EncoderDecoder represents object that is able to encode and decode Data structure.
type EncoderDecoder interface {
	Encode(*Data, io.Writer) error
	Decode(io.Reader) (*Data, error)
}

// ErrMalformedData is returned if some data in encoding/decoding are incorrect.
var ErrMalformedData = errors.New("malformed data")
