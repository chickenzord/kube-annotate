package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickenzord/kube-annotate/pkg/config"
)

//RulesHandler handles rules
func RulesHandler(w http.ResponseWriter, r *http.Request) {
	payload, err := json.Marshal(config.Rules)
	if err != nil {
		log.WithError(err).Error("cannot encode rules")
		http.Error(w, fmt.Sprintf("cannot encode rules: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(payload))
}
