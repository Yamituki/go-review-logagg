package models

import "time"

// Stats はログの統計情報を表す構造体です。
type Stats struct {
	// 総ログ数
	TotalCount int `json:"total_count"`
	// INFOレベルのログ数
	InfoCount int `json:"info_count"`
	// WARNレベルのログ数
	WarnCount int `json:"warn_count"`
	// ERRORレベルのログ数
	ErrorCount int `json:"error_count"`
	// 最初のログ時刻
	FirstTimestamp time.Time `json:"first_timestamp"`
	// 最後のログ時刻
	LastTimestamp time.Time `json:"last_timestamp"`
}
