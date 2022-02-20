package data

import (
	"encoding/json"
	"io"
	"net/http"
)

func FromJSON(data interface{}, r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(&data)
}

func ToJSON(rw http.ResponseWriter, data interface{}, status int, wrapper string) error {
	wrap := make(map[string]interface{})
	wrap[wrapper] = data

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(status)

	e := json.NewEncoder(rw)
	return e.Encode(wrap)
}

type GenericError struct {
	Message string `json:"message"`
}
