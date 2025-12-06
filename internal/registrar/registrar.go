// Package registrar provides JSON persistence for golf course data.
// It handles reading and writing user and standard fee data to JSON files.
package registrar

import (
	"encoding/json"
	"os"

	"golf/internal/membership"
)

// ReadJSON loads users from a JSON file.
// Returns a slice of users and any error encountered during file reading
// or JSON unmarshaling. If the file doesn't exist or is invalid, an error
// is returned along with an empty slice.
func ReadJSON(fileName string) ([]*membership.User, error) {
	datas := []*membership.User{}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return datas, err
	}
	err = json.Unmarshal(file, &datas)
	if err != nil {
		return datas, err
	}
	return datas, nil
}

// WriteJSON saves users to a JSON file with pretty-printing.
// The output is formatted with 4-space indentation for readability.
// Returns an error if marshaling fails or the file cannot be written.
// The file is created with permissions 0644 (readable by all, writable by owner).
func WriteJSON(fileName string, users []*membership.User) error {
	output, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, output, 0644)
	if err != nil {
		return err
	}
	return nil
}

// ReadStandardFeesJSON loads standard fees from a JSON file.
// Returns a slice of standard fees and any error encountered during
// file reading or JSON unmarshaling. If the file doesn't exist or is
// invalid, an error is returned along with an empty slice.
func ReadStandardFeesJSON(fileName string) ([]*membership.StandardFee, error) {
	fees := []*membership.StandardFee{}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return fees, err
	}
	err = json.Unmarshal(file, &fees)
	if err != nil {
		return fees, err
	}
	return fees, nil
}

// WriteStandardFeesJSON saves standard fees to a JSON file with pretty-printing.
// The output is formatted with 4-space indentation for readability.
// Returns an error if marshaling fails or the file cannot be written.
// The file is created with permissions 0644 (readable by all, writable by owner).
func WriteStandardFeesJSON(fileName string, fees []*membership.StandardFee) error {
	output, err := json.MarshalIndent(fees, "", "    ")
	if err != nil {
		return err
	}
	err = os.WriteFile(fileName, output, 0644)
	if err != nil {
		return err
	}
	return nil
}
