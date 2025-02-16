package jsonwriter

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Errors string `json:"errors"`
}

func WriteJSON(data any, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(data)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(resp)
	if err != nil {
		return err
	}

	return nil
}

func WriteJSONError(err error, w http.ResponseWriter, status int) error {
	w.Header().Set("Content-Type", "application/json")

	res, err := json.Marshal(&ErrorResponse{
		Errors: err.Error(),
	})
	if err != nil {
		return err
	}

	w.WriteHeader(status)

	_, err = w.Write(res)
	if err != nil {
		return err
	}

	return nil
}

func ReadReqJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1024 * 1024 // 1mb
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&data)
	if err != nil {
		return err
	}

	return nil
}

func WriteSuccess(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
