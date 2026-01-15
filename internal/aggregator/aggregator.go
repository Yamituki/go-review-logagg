package aggregator

import "github.com/Yamituki/go-review-logagg/pkg/models"

// Aggregator はログデータを集約するためのインターフェースです。
type Aggregator interface {
	// 1つのエントリを追加
	Add(entry models.LogEntry) error
	// 統計情報を取得
	GetStats() models.Stats
	// 統計をリセット
	Reset()
}
