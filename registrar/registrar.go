package registrar

import (
	"encoding/json"
	"golf/membership"
	"os"
)

func ReadJSON(fileName string) []*membership.User {
	datas := []*membership.User{}
	file, _ := os.ReadFile(fileName)
	json.Unmarshal(file, &datas)
	return datas
}

func WriteJSON(fileName string, users []*membership.User) error {
	output, err := json.MarshalIndent(users, "", "    ")
	if err == nil {
		os.WriteFile(fileName, output, 0644)
	} else {
		return err
	}
	return nil
}
