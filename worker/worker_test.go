package main

import (
	"sync"
	"testing"
	"time"
)

func TestWorker_ManagedWithTContext(t *testing.T) {
	// 1. テスト開始
	t.Parallel()

	// 2. t.Context() を用意する
	ctx := t.Context()

	// 3. ワーカーを起動する
	var wg sync.WaitGroup
	wg.Add(1)
	produced := make(chan struct{}, 1) // ワーカーが動いていることを外部から観測する用のチャネル

	go func() {
		// 8. ワーカーが終了する
		defer wg.Done()
		// 4. ワーカーが動く
		worker(ctx, produced) // テスト対象の関数。50ms周期で処理を行い、ctx.Done() で停止するワーカー
	}()

	// 5. workerが動いていることを確認し、200ms以内に観測できなければテストを失敗させる
	select {
	case <-produced:
		// 6. workerが動いており、チャネルへの送信がある
	case <-time.After(200 * time.Millisecond):
		t.Fatal("worker did not produce any output in time")
	}

	// 7. t.Cleanupの実行前に t.Context()が自動的にキャンセルされる
	// 9. t.CleanupとしてwaitGroupのWaitを実行する。既に goroutine が終了していることが保証されているため、ブロックされることはない
	t.Cleanup(wg.Wait)
}
