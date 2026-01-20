package monitor

/*
 * context　パッケージは、キャンセル可能なコンテキストを提供します。
 * fmt パッケージはフォーマットされたI/Oを提供します
 * os パッケージはプラットフォーム非依存のOS機能を提供します。
 * sync パッケージは基本的な同期プリミティブを提供します。
 * time パッケージは時間の測定と表示を提供します。
 */
import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/Yamituki/go-review-logagg/internal/processor"
	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// FileMonitor はファイルシステムの監視を行うための構造体です。
type FileMonitor struct {
	// 監視対象のファイルパス
	filePath string
	// 監視間隔
	interval time.Duration
	// データ処理用のコンカレントプロセッサ
	processor processor.ConcurrentProcessor
	// 監視統計情報
	stats models.Stats
	// 最終更新時刻
	lastModTime time.Time
	// 監視の停止を制御するチャネル
	mutex sync.Mutex
	// コンテキスト
	ctx context.Context
	// コンテキストのキャンセル関数
	cancel context.CancelFunc
	// 監視の停止を通知するチャネル
	done chan struct{}
}

// NewFileMonitor は新しい FileMonitor インスタンスを作成します。
func NewFileMonitor(filePath string, interval time.Duration) *FileMonitor {
	// コンテキストとキャンセル関数の作成
	ctx, cancel := context.WithCancel(context.Background())

	return &FileMonitor{
		filePath:  filePath,
		interval:  interval,
		processor: *processor.NewConcurrentProcessor(2),
		stats:     models.Stats{},
		mutex:     sync.Mutex{},
		ctx:       ctx,
		cancel:    cancel,
		done:      make(chan struct{}),
	}
}

// Start はファイル監視を開始します。
func (fm *FileMonitor) Start() error {
	// 監視間隔
	ticker := time.NewTicker(fm.interval)

	// 監視goroutineの開始
	go func() {

		// Ticker のクリーンアップ
		defer ticker.Stop()

		// 監視ループ
		for {
			select {
			case <-fm.ctx.Done():
				// 停止完了を通知
				close(fm.done)
				return
			case <-ticker.C:
				// ファイルの変更をチェックして更新
				if err := fm.checkAndUpdate(); err != nil {
					fmt.Printf("監視エラー: %v\n", err)
				}
			}
		}
	}()

	return nil
}

// Stop はファイル監視を停止します。
func (fm *FileMonitor) Stop() error {
	// 監視の停止を通知
	fm.cancel()
	// 停止完了を待機
	<-fm.done

	return nil
}

// GetStats は現在の監視統計情報を取得します。
func (fm *FileMonitor) GetStats() (models.Stats, error) {
	// ミューテックスのロック
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	return fm.stats, nil
}

// checkAndUpdate はファイルの変更をチェックし、必要に応じて更新を行います。
func (fm *FileMonitor) checkAndUpdate() error {
	// ファイルの情報を取得
	fileInfo, err := os.Stat(fm.filePath)
	if err != nil {
		return err
	}

	// 変更時刻を取得
	modTime := fileInfo.ModTime()

	if fm.lastModTime.IsZero() {
		// 初回実行時は最終更新時刻を設定
		fm.lastModTime = modTime
	} else if !modTime.After(fm.lastModTime) {
		// 変更がない場合は終了
		return nil
	}

	// ファイルの読み込み
	stats, err := fm.processor.ProcessFiles([]string{fm.filePath})
	if err != nil {
		return err
	}

	// ミューテックスのロック
	fm.mutex.Lock()
	defer fm.mutex.Unlock()

	// 統計情報の更新
	fm.stats = stats
	// 最終更新時刻の更新
	fm.lastModTime = modTime

	return nil
}
