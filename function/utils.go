package function

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Chiorufarewerin/gitchat/internal/environment"
	"github.com/google/uuid"
)

func GenerateUniqueId() string {
	return uuid.New().String()
}

func MakeJsonField(name string, value interface{}) string {
	marshaled, err := json.Marshal(value)
	if err != nil {
		marshaled = []byte("null")
	}
	return fmt.Sprintf("\"%s\": %s", name, marshaled)
}

func GetCurrentDateUTCString() string {
	return time.Now().UTC().Format(environment.DateFormat)
}
