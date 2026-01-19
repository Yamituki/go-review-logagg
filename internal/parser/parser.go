package parser

import "github.com/Yamituki/go-review-logagg/pkg/models"

// LogParser はログ行を解析して LogEntry 構造体に変換するためのインターフェースです。
type LogParser interface {
	Parse(line string) (models.LogEntry, error)
}
