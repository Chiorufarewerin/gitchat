package function

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	_ "github.com/Chiorufarewerin/gitchat/internal/environment"
	"github.com/Chiorufarewerin/gitchat/internal/git"
)

func writeError(w http.ResponseWriter, err error, status int) {
	errorResponse := map[string]map[string]interface{}{
		"error": {
			"status":  status,
			"message": err.Error(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse)
}

func writeSuccess(w http.ResponseWriter, value interface{}) {
	successResponse := map[string]interface{}{
		"data": value,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(successResponse)
}

func Handler(w http.ResponseWriter, r *http.Request) {
	comment, err := ReadAndValidateComment(r.Body)
	if err != nil {
		writeError(w, err, 400)
		return
	}

	comment, err = AddComment(comment)
	if err != nil {
		log.Println(err)
		writeError(w, errors.New("System error"), 500)
		return
	}
	writeSuccess(w, comment)
}

func init() {
	git.InitializeGit()
}
