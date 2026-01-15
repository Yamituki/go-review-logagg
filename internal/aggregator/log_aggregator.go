package aggregator

import "github.com/Yamituki/go-review-logagg/pkg/models"

// LogAggregator はログデータを集約するための構造体です。
type LogAggregator struct {
	// 保持するログ一覧
	entries []models.LogEntry
	// 統計情報
	stats models.Stats
}

// NewLogAggregator は LogAggregator の新しいインスタンスを作成します。
func NewLogAggregator() *LogAggregator {
	return &LogAggregator{}
}

// Add は1つのログエントリを追加します。
func (la *LogAggregator) Add(entry models.LogEntry) error {
	la.entries = append(la.entries, entry)
	la.updateStats(entry)
	return nil
}

// GetStats は現在のログエントリに基づいて統計情報を取得します。
func (la *LogAggregator) GetStats() models.Stats {
	return la.stats
}

// Reset は集約されたログデータと統計情報をリセットします。
func (la *LogAggregator) Reset() {
	la.entries = []models.LogEntry{}
	la.stats = models.Stats{}
}

// 統計情報の更新メソッド
func (la *LogAggregator) updateStats(entry models.LogEntry) {
	// 総ログ数の更新
	la.stats.TotalCount++

	// レベル別ログ数の更新
	switch entry.Level {
	case "INFO":
		la.stats.InfoCount++
	case "WARN":
		la.stats.WarnCount++
	case "ERROR":
		la.stats.ErrorCount++
	}

	// 最初と最後のタイムスタンプの初期化
	if la.stats.TotalCount == 1 {
		la.stats.FirstTimestamp = entry.Timestamp
		la.stats.LastTimestamp = entry.Timestamp
		return
	}

	// 最初のタイムスタンプの更新
	if entry.Timestamp.Before(la.stats.FirstTimestamp) || la.stats.TotalCount == 1 {
		la.stats.FirstTimestamp = entry.Timestamp
	}

	// 最後のタイムスタンプの更新
	if entry.Timestamp.After(la.stats.LastTimestamp) {
		la.stats.LastTimestamp = entry.Timestamp
	}
}
