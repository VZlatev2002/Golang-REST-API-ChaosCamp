package helpers

import (
	"encoding/json"
	"log"
	"net/http"
)

// ErrorMessage Convenience function to redirect to the error message page
func ErrorMessage(writer http.ResponseWriter, request *http.Request, msg string) {
	http.Redirect(writer, request, "localhost:8080/error", 302)
}

// RespondWithJSON marshals the payload to a json and sends response via the ResponseWriter
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		log.Fatal(err)
	}
}