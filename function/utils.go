package function

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
)

func GenerateUniqueId() string {
	return uuid.New().String()
}

func MakeJsonField(name string, value interface{}) (string, error) {
	marshaled, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("\"%s\": %s", name, marshaled), nil
}
