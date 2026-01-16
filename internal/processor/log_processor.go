package processor

import (
	"github.com/Yamituki/go-review-logagg/internal/aggregator"
	"github.com/Yamituki/go-review-logagg/internal/parser"
	"github.com/Yamituki/go-review-logagg/internal/reader"
	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// LogProcessor はログを処理するための構造体です。
type LogProcessor struct{}

// NewLogProcessor は新しい LogProcessor インスタンスを作成します。
func NewLogProcessor() *LogProcessor {
	return &LogProcessor{}
}

func (lp *LogProcessor) ProcessFile(filePath string) (models.Stats, error) {
	// ファイルリーダーの初期化
	fr := reader.NewFileReader(filePath)

	var stats models.Stats
	var le models.LogEntry
	var line string
	var err error

	// パーサーの初期化
	ps := parser.NewStandardParser()

	// アグリゲーターの初期化
	ag := aggregator.NewLogAggregator()

	// すべての行を読み込む
	var lines []string
	lines, err = fr.ReadAllLines()
	if err != nil {
		return stats, err
	}

	// 各行を処理
	for _, line = range lines {
		// ログ行の解析
		le, err = ps.Parse(line)
		if err != nil {
			return stats, err
		}

		// 統計情報の更新
		ag.Add(le)
	}

	// 最終的な統計情報を取得
	stats = ag.GetStats()

	return stats, nil
}
