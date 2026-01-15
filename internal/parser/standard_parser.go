package parser

/*
 * time パッケージは時間の操作を提供します。
 */
import (
	"fmt"
	"time"

	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// StandardParser は標準的なログ解析を行う構造体です。
type StandardParser struct{}

// NewStandardParser は StandardParser の新しいインスタンスを作成します。
func NewStandardParser() *StandardParser {
	return &StandardParser{}
}

// Parse はログ行を解析し、LogEntry 構造体に変換します。
func (sp *StandardParser) Parse(line string) (models.LogEntry, error) {

	var entry models.LogEntry

	// ログの基本フォーマット: "YYYY-MM-DD HH:MM:SS [LEVEL] Message"

	// 空行のチェック
	if len(line) == 0 {
		return entry, fmt.Errorf("空のログ行は解析できません")
	}

	// 日付と時間の解析
	timestampStr := line[0:19]
	timestamp, err := time.Parse("2006-01-02 15:04:05", timestampStr)
	if err != nil {
		return entry, err
	}
	entry.Timestamp = timestamp

	// レベルの解析 [INFO], [ERROR], [WARN]
	levelStart := 20
	levelEnd := 0

	// 残りの部分からレベルを抽出
	for i := levelStart; i < len(line); i++ {
		if line[i] == ']' {
			levelEnd = i
			break
		}
	}

	switch line[levelStart+1 : levelEnd] {
	case "INFO":
		entry.Level = "INFO"
	case "ERROR":
		entry.Level = "ERROR"
	case "WARN":
		entry.Level = "WARN"
	default:
		return entry, fmt.Errorf("不明なログレベル: %s", line[levelStart+1:levelEnd])
	}

	// メッセージの解析
	messageStart := levelEnd + 2

	// メッセージが存在しない場合のエラーハンドリング
	if messageStart >= len(line) {
		return entry, fmt.Errorf("ログ行のメッセージが存在しません: %s", line)
	}

	entry.Message = line[messageStart:]

	return entry, nil
}
