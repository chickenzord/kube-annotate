package handler

import (
	"net/http"

	"github.com/gamegos/jsend"
)

func writeJsend(w http.ResponseWriter, msg string, data interface{}, code int) {
	jsend.Wrap(w).
		Status(code).
		Message(msg).
		Data(data).
		Send()
}
