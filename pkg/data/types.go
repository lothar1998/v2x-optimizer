package data

import "io"

// Data is a structure that represents data required to the computation of the solution.
type Data struct {
	MBR []int
	R   [][]int
}

// EncoderDecoder represents object that is able to encode and decode Data structure.
type EncoderDecoder interface {
	Encode(*Data, io.Writer) error
	Decode(reader io.Reader) (*Data, error)
}
