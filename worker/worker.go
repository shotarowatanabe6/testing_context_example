package main

import (
	"context"
	"time"
)

// 50ms周期で処理を行い、ctx.Done() で停止するワーカー
func worker(ctx context.Context, produced chan<- struct{}) {
	ticker := time.NewTicker(50 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// workerが動いていることを外部から観測できるよう、チャネルへ送信する
			select {
			case produced <- struct{}{}:
			default:
			}
		}
	}
}
