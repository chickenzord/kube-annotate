package annotator

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickenzord/kube-annotate/config"
)

//Annotator annotate pod admission
type Annotator struct {
}

//MutateHandler handles admission mutation
func MutateHandler(w http.ResponseWriter, r *http.Request) {
	admissionReview, err := parseBody(r)
	if err != nil {
		http.Error(w, "Cannot parse body", http.StatusBadRequest)
		return
	}

	var result = mutate(admissionReview)
	resp, err := json.Marshal(result)
	if err != nil {
		log.Errorf("Can't encode response: %v", err)
		http.Error(w, fmt.Sprintf("could not encode response: %v", err), http.StatusInternalServerError)
	}
	if _, err := w.Write(resp); err != nil {
		log.Errorf("Can't write response: %v", err)
		http.Error(w, fmt.Sprintf("could not write response: %v", err), http.StatusInternalServerError)
	}
}

//RulesHandler handles rules
func RulesHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(config.Rules)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(payload))
}
