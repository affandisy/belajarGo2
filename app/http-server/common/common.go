package common

import (
	"encoding/json"
	"net/http"
)

func ErrorInvalidJJSON(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "invalid JSON request", "data": map[string]interface{}{}})
}

func ErrorValidation(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": "invalid parameter request", "data": map[string]interface{}{}})
}

func ErrorInternal(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(http.StatusInternalServerError), "data": map[string]interface{}{}})
}

func ErrorDataConflict(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusConflict)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(http.StatusConflict), "data": map[string]interface{}{}})
}

func ErrorDataNotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(http.StatusNotFound), "data": map[string]interface{}{}})
}

func ValidResponse(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_ = json.NewEncoder(w).Encode(map[string]interface{}{"message": http.StatusText(httpStatus), "data": data})
}
