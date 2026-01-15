package models

/*
 * time パッケージは時間の操作を提供します。
 */
import "time"

// LogEntry はログエントリを表す構造体です。
type LogEntry struct {
	// Logのタイムスタンプ
	Timestamp time.Time `json:"timestamp"`
	// Logのレベル INFO, ERROR, WARN
	Level string `json:"level" enum:"INFO,ERROR,WARN"`
	// Logのメッセージ内容
	Message string `json:"message"`
	// Logの発生源 (例: サービス名やホスト名)
	Source string `json:"source"`
}
