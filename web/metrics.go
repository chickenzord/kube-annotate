package web

import (
	"fmt"
	"net/http"
)

//MetricsHandler handles prometheus metrics
func MetricsHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "status 1")
}
