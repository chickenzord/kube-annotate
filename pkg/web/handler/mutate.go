package handler

import (
	"fmt"
	"io/ioutil"
	"mime"
	"net/http"

	"github.com/chickenzord/kube-annotate/pkg/mutator"
)

//MutateHandler handles admission mutation
func MutateHandler(w http.ResponseWriter, r *http.Request) {
	// Parse request
	if r.ContentLength == 0 {
		writeJsend(w, "empty body", nil, http.StatusBadRequest)
		return
	}
	contentType, _, _ := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if contentType != "application/json" {
		msg := fmt.Sprintf("invalid content type: %s", contentType)
		writeJsend(w, msg, nil, http.StatusBadRequest)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil || data == nil {
		writeJsend(w, "cannot read body", err.Error(), http.StatusBadRequest)
		return
	}

	// Write response
	response, err := mutator.MutateBytes(data)
	if err != nil {
		writeJsend(w, "cannot mutate payload", err.Error(), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(response); err != nil {
		writeJsend(w, "cannot write response", err.Error(), http.StatusInternalServerError)
		return
	}
}
