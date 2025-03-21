package handler

import "net/http"

type HealthHandler struct{}

// NewHealthHandler создаёт новый экземпляр HealthHandler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// CheckHealth проверяет состояние сервера
//
// @Summary Проверка состояния сервиса
// @Description Эндпоинт для проверки работоспособности API
// @Tags health
// @Produce plain
// @Success 200 {string} string "Its alive!"
// @Router /health [get]
func (hh *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Its alive!"))
}
