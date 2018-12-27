package handler

import (
	"fmt"
	"net/http"
)

//HealthHandler handles health checks
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "OK")
}
