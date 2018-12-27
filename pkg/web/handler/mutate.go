package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickenzord/kube-annotate/pkg/mutator"
)

//MutateHandler handles admission mutation
func MutateHandler(w http.ResponseWriter, r *http.Request) {
	admissionReview, err := mutator.ParseBody(r)
	if err != nil {
		log.WithError(err).Error("cannot parse body")
		http.Error(w, "cannot parse body", http.StatusBadRequest)
		return
	}

	var result = mutator.Mutate(admissionReview)
	resp, err := json.Marshal(result)
	if err != nil {
		log.WithError(err).Error("cannot encode response")
		http.Error(w, fmt.Sprintf("cannot encode response: %v", err), http.StatusInternalServerError)
		return
	}
	if _, err := w.Write(resp); err != nil {
		log.WithError(err).Error("cannot write response")
		http.Error(w, fmt.Sprintf("cannot write response: %v", err), http.StatusInternalServerError)
		return
	}
}
