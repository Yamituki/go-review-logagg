package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Yamituki/go-review-logagg/internal/processor"
)

// jsonRequest は JSON リクエストの共通構造を表します。
type jsonRequest struct {
	Filepath string `json:"filepath"`
}

// jsonResponse は JSON レスポンスの共通構造を表します。
type jsonResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

// handleHealth はヘルスチェックのハンドラーです。
func handleHealth(w http.ResponseWriter, r *http.Request) {
	// 戻り値の型は jsonResponse を使用します。
	w.Header().Set("Content-Type", "application/json")

	// 処理結果の状態を返します。
	w.WriteHeader(http.StatusOK)

	// レスポンスボディを JSON 形式で返します。
	var resp jsonResponse
	resp.Status = "ok"
	resp.Data = nil
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","data":"レスポンスの生成に失敗しました: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}

// handleAnalyze はログ集約のハンドラーです。
func handleAnalyze(w http.ResponseWriter, r *http.Request) {
	// 戻り値の型は jsonResponse を使用します。
	w.Header().Set("Content-Type", "application/json")

	// リクエストボディの解析と処理を行います。
	body := r.Body

	// 終了時にボディを閉じます。
	defer body.Close()

	var req jsonRequest
	var resp jsonResponse

	// リクエストの解析

	decoder := json.NewDecoder(body)
	if err := decoder.Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","data":"リクエストの解析に失敗しました: %s"}`, err.Error()), http.StatusBadRequest)
		return
	}

	// ログファイルの解析処理
	ps := processor.NewLogProcessor()
	stats, err := ps.ProcessFile(req.Filepath)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","data":"ログファイルの解析に失敗しました: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	// 処理結果の状態を返します。
	w.WriteHeader(http.StatusOK)

	// レスポンスボディを JSON 形式で返します。
	resp.Status = "ok"
	resp.Data = stats
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"status":"error","data":"レスポンスの生成に失敗しました: %s"}`, err.Error()), http.StatusInternalServerError)
		return
	}

	w.Write(jsonResp)
}
