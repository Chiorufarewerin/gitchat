package function

import (
	"encoding/json"
	"errors"
	"io"
)

const numbers = "0123456789"

func isColorCorrect(color string) bool {
	if color == "" {
		return true
	}
	if len(color) != 7 || color[0] != '#' {
		return false
	}
	for colorNumber := range color[1:] {
		isCorrect := false
		for num := range numbers {
			if colorNumber == num {
				isCorrect = true
				break
			}
		}
		if !isCorrect {
			return false
		}
	}
	return true
}

func ReadAndValidateComment(r io.Reader) (*Comment, error) {
	requestData := &RequestData{}
	if err := json.NewDecoder(r).Decode(requestData); err != nil {
		return nil, errors.New("Incorrect body data")
	}

	if requestData.Version != "v1" {
		return nil, errors.New("Incorrect field: version")
	}

	if len(requestData.Author) < 1 || len(requestData.Author) > 50 {
		return nil, errors.New("Incorrect field: author")
	}

	if len(requestData.Text) < 1 || len(requestData.Text) > 5000 {
		return nil, errors.New("Incorrect field: text")
	}

	if !isColorCorrect(requestData.Color) {
		return nil, errors.New("Incorrect field: color")
	}

	comment := &Comment{
		ID:     "",
		Author: requestData.Author,
		Text:   requestData.Text,
		Color:  requestData.Color,
	}

	return comment, nil
}
