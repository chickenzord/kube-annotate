package annotator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickenzord/kube-annotate/config"
)

//MutateHandler handles admission mutation
func MutateHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	admissionReview, err := parseBody(r)
	if err != nil {
		log.WithError(err).Error("cannot parse body")
		http.Error(w, "cannot parse body", http.StatusBadRequest)
		return
	}

	var result = mutate(admissionReview)
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

//RulesHandler handles rules
func RulesHandler(w http.ResponseWriter, r *http.Request) {
	logRequest(r)

	payload, err := json.Marshal(config.Rules)
	if err != nil {
		log.WithError(err).Error("cannot encode rules")
		http.Error(w, fmt.Sprintf("cannot encode rules: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(payload))
}
