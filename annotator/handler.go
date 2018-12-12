package annotator

import (
	"encoding/json"
	"fmt"
	"net/http"
)

//Annotator annotate pod admission
type Annotator struct {
}

func (annotator *Annotator) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
