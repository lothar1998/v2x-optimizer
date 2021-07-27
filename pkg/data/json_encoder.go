package data

import (
	"encoding/json"
	"io"
)

// JSONEncoder provide methods for encoding an decoding Data structure into/from json.
type JSONEncoder struct{}

// Encode facilitates encoding Data to json.
func (e JSONEncoder) Encode(data *Data, w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(data)
}

// Decode allows for conveniently decoding data from json format to Data structure.
func (e JSONEncoder) Decode(r io.Reader) (*Data, error) {
	var data Data

	decoder := json.NewDecoder(r)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
