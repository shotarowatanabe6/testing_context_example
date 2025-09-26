//go:generate mockgen -source=service.go -destination=mock_service.go -package=service
package service

import (
	"context"
	"encoding/json"
	"net/http"
)

type Service interface {
	GetData(ctx context.Context, id string) (any, error)
}

type Handler struct {
	Svc Service
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	id := r.URL.Query().Get("id")
	data, err := h.Svc.GetData(ctx, id)
	if err != nil {
		// コンテキスト由来のエラーは 503、それ以外は 500 の例
		select {
		case <-ctx.Done():
			http.Error(w, "request canceled", http.StatusServiceUnavailable)
		default:
			http.Error(w, "internal error", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
