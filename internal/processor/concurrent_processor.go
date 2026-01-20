package processor

/*
 * fmt パッケージはフォーマットされたI/Oを提供します
 * sync パッケージは基本的な同期プリミティブを提供します。
 */
import (
	"fmt"
	"sync"

	"github.com/Yamituki/go-review-logagg/internal/aggregator"
	"github.com/Yamituki/go-review-logagg/internal/parser"
	"github.com/Yamituki/go-review-logagg/internal/reader"
	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// ConcurrentProcessor は並行処理を行うプロセッサの構造体です。
type ConcurrentProcessor struct {
	workers int
}

// NewConcurrentProcessor は ConcurrentProcessor の新しいインスタンスを作成します。
func NewConcurrentProcessor(workers int) *ConcurrentProcessor {
	return &ConcurrentProcessor{
		workers: workers,
	}
}

// ProcessFiles は指定されたファイルパスのログファイルを並行して処理します。
func (cp *ConcurrentProcessor) ProcessFiles(filePaths []string) (models.Stats, error) {

	// ファイルのパスを受け取るチャネル
	fileChan := make(chan string, len(filePaths))

	// 結果を受け取るチャネル
	resultChan := make(chan models.Stats, len(filePaths))

	// 同期用のWaitGroup
	var waitGroup sync.WaitGroup

	// 統計情報の初期化
	var stats models.Stats

	// ワーカーの数だけWaitGroupにカウントを追加
	waitGroup.Add(cp.workers)

	// エラーを収集
	var firstError error
	var errorMutex sync.Mutex

	// ワーカーを起動
	for i := 0; i < cp.workers; i++ {
		go func() {
			// ワーカーが終了したらWaitGroupのカウントをデクリメント
			defer waitGroup.Done()

			// ファイルパスをチャネルから受け取る
			for filepath := range fileChan {

				// リーダーの初期化
				reader := reader.NewFileReader(filepath)

				// ファイルの内容を読み込む
				var lines []string
				var err error
				lines, err = reader.ReadAllLines()
				if err != nil {
					errorMutex.Lock()
					if firstError == nil {
						firstError = fmt.Errorf("ファイルの読み込みに失敗しました: %v", err)
					}
					errorMutex.Unlock()
					continue
				}

				// パーサーの初期化
				parser := parser.NewStandardParser()

				// 集約器の初期化
				aggregator := aggregator.NewLogAggregator()

				// 各行をパースして集計
				var entry models.LogEntry
				for _, line := range lines {
					entry, err = parser.Parse(line)
					if err != nil {
						errorMutex.Lock()
						if firstError == nil {
							firstError = fmt.Errorf("ログのパースに失敗しました: %v", err)
						}
						errorMutex.Unlock()
						continue
					}

					// 集約器に追加
					aggregator.Add(entry)

				}

				// 結果をチャネルに送信
				resultChan <- aggregator.GetStats()

			}

		}()
	}

	// ファイルをチャネルに送信
	for _, path := range filePaths {
		fileChan <- path
	}

	// チャネルを閉じる
	close(fileChan)

	waitGroup.Wait()

	// チャネルを閉じる
	close(resultChan)

	// 結果をチャネルから受信して集約
	for range filePaths {
		result := <-resultChan

		// フィルター: 結果が空の場合はスキップ
		if result.InfoCount == 0 && result.WarnCount == 0 && result.ErrorCount == 0 {
			continue
		}

		// 統計情報の集約
		stats.TotalCount += result.TotalCount
		stats.ErrorCount += result.ErrorCount
		stats.WarnCount += result.WarnCount
		stats.InfoCount += result.InfoCount

		// タイムスタンプの初期化
		if stats.FirstTimestamp.IsZero() && !result.FirstTimestamp.IsZero() {
			stats.FirstTimestamp = result.FirstTimestamp
		}

		if stats.LastTimestamp.IsZero() && !result.LastTimestamp.IsZero() {
			stats.LastTimestamp = result.LastTimestamp
		}

		// 最小タイムスタンプの更新
		if result.FirstTimestamp.Before(stats.FirstTimestamp) {
			stats.FirstTimestamp = result.FirstTimestamp
		}

		// 最大タイムスタンプの更新
		if result.LastTimestamp.After(stats.LastTimestamp) {
			stats.LastTimestamp = result.LastTimestamp
		}
	}

	return stats, firstError
}
