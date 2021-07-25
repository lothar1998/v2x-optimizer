package data

import (
	"encoding/json"
	"os"
)

// Load allows for conveniently loading data in json format to Data structure
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

// Export facilitates exporting Data to json file in given filepath.
// Path should exist before exporting.
func Export(filepath string, data *Data) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(data)
	return err
}
