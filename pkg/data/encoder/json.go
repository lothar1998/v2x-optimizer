package encoder

import (
	"encoding/json"
	"io"

	"github.com/lothar1998/v2x-optimizer/pkg/data"
)

// JSONEncoder provide methods for encoding and decoding Data structure into/from json.
type JSONEncoder struct{}

// Encode facilitates encoding Data to json.
func (e JSONEncoder) Encode(input *data.Data, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(input)
}

// Decode allows for conveniently decoding data from json format to Data structure.
func (e JSONEncoder) Decode(r io.Reader) (*data.Data, error) {
	var output data.Data

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&output)
	if err != nil {
		return nil, data.ErrMalformedData
	}

	return &output, nil
}
