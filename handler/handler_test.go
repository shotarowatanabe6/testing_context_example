package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"testing_context_example/service"

	"go.uber.org/mock/gomock"
)

func TestHandler_ServeHTTP_Success(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := service.NewMockService(ctrl)

	// 期待: GetData が t.Context()（経由の r.Context）で呼ばれ、正常データを返す
	mockSvc.EXPECT().GetData(gomock.Any(), "key1").
		DoAndReturn(func(ctx context.Context, id string) (any, error) {
			// ctx は t.Context() を親にしており、まだキャンセルされていないためdefaultで正常データを返す
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return map[string]any{"id": id, "ok": true}, nil
			}
		})

	h := Handler{Svc: mockSvc}

	req := httptest.NewRequest(http.MethodGet, "/?id=key1", nil)
	req = req.WithContext(t.Context()) // t.Context をリクエストへ伝播
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status=%d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandler_ServeHTTP_Canceled(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := service.NewMockService(ctrl)

	// 期待: GetData はコンテキストのキャンセルに追従し、context.Canceled を返す
	mockSvc.EXPECT().GetData(gomock.Any(), "key1").
		DoAndReturn(func(ctx context.Context, id string) (any, error) {
			// コンテキストがキャンセルされるまで待つような処理を模倣
			<-ctx.Done()
			return nil, ctx.Err()
		})

	h := Handler{Svc: mockSvc}

	// t.Context を親に明示的な短いタイムアウトを設定して「途中でキャンセル」を再現
	ctx, cancel := context.WithTimeout(t.Context(), 50*time.Millisecond)
	t.Cleanup(cancel)

	req := httptest.NewRequest(http.MethodGet, "/?id=key1", nil).WithContext(ctx)
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	// ハンドラーはキャンセルを 503 にマップする仕様
	if rr.Code != http.StatusServiceUnavailable {
		t.Fatalf("want 503, got %d, body=%s", rr.Code, rr.Body.String())
	}
}

func TestHandler_ServeHTTP_InternalError(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockSvc := service.NewMockService(ctrl)

	mockSvc.EXPECT().GetData(gomock.Any(), "key1").
		Return(nil, errors.New("boom"))

	h := Handler{Svc: mockSvc}

	req := httptest.NewRequest(http.MethodGet, "/?id=key1", nil).
		WithContext(t.Context())
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("want 500, got %d", rr.Code)
	}
}
