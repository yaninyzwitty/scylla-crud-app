package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
)

func NewResponseToJson(w http.ResponseWriter, status int, data interface{}) error {
	response, err := json.Marshal(data)
	if err != nil {
		return errors.New("error while marshalling response")
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
	return nil

}
