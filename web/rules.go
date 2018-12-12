package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chickenzord/kube-annotate/config"
)

//RulesHandler handles rules
func RulesHandler(w http.ResponseWriter, r *http.Request) {
	payload, _ := json.Marshal(config.Rules)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, string(payload))
}
