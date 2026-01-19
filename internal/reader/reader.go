package reader

// LogReader はログを読み込むためのインターフェースです。
type LogReader interface {
	// 1行を読む
	ReadLine() (string, error)
	// 全行を読む
	ReadAllLines() ([]string, error)
}
