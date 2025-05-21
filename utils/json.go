package utils

import (
	"encoding/json"
	"net/http"
)

func ReadJSON(r *http.Request, dst any ) error{
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&dst)
	if err != nil {
		return err
	}
	
	return nil
}

func WriteJSON(w http.ResponseWriter, statusCode int, data any) error {
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	jsonData = append(jsonData, '\n')
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(jsonData)
	return nil
}
