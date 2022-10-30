package controllers

import "net/http"

type HealthController struct{}

func (c *HealthController) GetHealthStatus(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Healthy"))
}
