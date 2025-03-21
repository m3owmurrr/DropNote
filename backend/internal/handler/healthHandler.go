package handler

import "net/http"

type HealthHandler struct{}

// NewHealthHandler создаёт новый экземпляр HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

func (hh *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Its alive!"))
}
