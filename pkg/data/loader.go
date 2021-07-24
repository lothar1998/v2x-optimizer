package data

import (
	"encoding/json"
	"os"
)

// Load allows for conveniently loading
// data in json format to Data structure
func Load(filepath string) (*Data, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data Data

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return &data, nil
}
