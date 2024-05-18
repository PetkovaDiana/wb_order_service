package handler

import (
	"fmt"
	"net/http"
	"projectOrder/internal/service"
	"time"
)

const (
	basePath = "/api"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *http.ServeMux {
	mx := http.NewServeMux()

	mx.HandleFunc(fmt.Sprintf("%s/order-info", basePath), h.baseHandler(h.orderByID)) // GET

	return mx
}

func (h *Handler) baseHandler(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(time.Now())
		fmt.Println(r.Method)

		handlerFunc(w, r)
	}
}
