package monitor

import "github.com/Yamituki/go-review-logagg/pkg/models"

// Monitor はシステムの監視を行うためのインターフェースです。
type Monitor interface {
	// Start は監視を開始します。
	Start() error
	// Stop は監視を停止します。
	Stop() error
	// GetStats は現在の監視統計情報を取得します。
	GetStats() (models.Stats, error)
}
