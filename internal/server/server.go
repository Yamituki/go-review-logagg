package server

import (
	"net/http"

	"github.com/Yamituki/go-review-logagg/internal/processor"
)

// Server はログ集約サーバーを表します。
type Server struct {
	port      string
	processor *processor.LogProcessor
}

// NewServer は新しい Server インスタンスを作成します。
func NewServer(port string) *Server {
	return &Server{
		port:      port,
		processor: processor.NewLogProcessor(),
	}
}

// Start はサーバーを起動します。
func (s *Server) Start() error {
	return http.ListenAndServe(s.port, nil)
}

// SetupRoutes はサーバーのルートを設定します。
func (s *Server) SetupRoutes() {
	// ヘルスチェックのエンドポイント
	http.HandleFunc("/health", handleHealth)

	// ログ集約のエンドポイント
	http.HandleFunc("/analyze", handleAnalyze)
}
