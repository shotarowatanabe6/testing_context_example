package handler

import (
	"encoding/json"
	"net/http"

	"testing_context_example/service"
)

type Handler struct {
	Svc service.Service
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")

	data, err := h.Svc.GetData(ctx, id)
	if err != nil {
		select {
		case <-ctx.Done():
			http.Error(w, "request canceled", http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}
