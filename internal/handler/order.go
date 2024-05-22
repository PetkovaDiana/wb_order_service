package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *Handler) orderByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	orderID := r.URL.Query().Get("order_id")
	if orderID == "" {
		http.Error(w, "Order ID is required", http.StatusBadRequest)
		return
	}

	order, err := h.services.GetOrderById(orderID)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Order not found", http.StatusBadRequest)
		return
	}

	resp, err := json.Marshal(order)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to encode order", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}
