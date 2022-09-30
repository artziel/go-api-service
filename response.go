package rest

import (
	"encoding/json"
	"net/http"
)

/*
Encode a interface data in JSON format and write to the response writer
*/
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

/*
Encode a golang error in JSON format and write to the response writer
*/
func RespondWithJSONError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, map[string]string{"message": err.Error()})
}

/*
Encode a text message in JSON format and write to the response writer
*/
func RespondWithJSONMessage(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"message": message})
}

/*
Encode a interface data in JSON format and write to the response writer wraped in
a hmac secured structure
*/
func RespondWithJSONHMAC(w http.ResponseWriter, code int, payload interface{}, secretKey string) {
	encoded, _ := json.Marshal(payload)

	hmac := NewHash(string(encoded), secretKey)

	w.Header().Set("Service-Content-Hash", hmac)
	RespondWithJSON(w, code, payload)
}
