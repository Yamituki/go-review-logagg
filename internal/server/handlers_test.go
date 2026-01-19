package server

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// TestHandleHealth は handleHealth ハンドラーのテストを行います。
func TestHandleHealth(t *testing.T) {
	// テストリクエスト
	testReq := httptest.NewRequest(http.MethodGet, "/health", nil)
	testRec := httptest.NewRecorder()

	// ハンドラーの呼び出し
	handleHealth(testRec, testReq)

	t.Logf("ステータスコード: %d", testRec.Code)
	t.Logf("レスポンスボディ: %s", testRec.Body.String())

	// ステータスコードの検証
	if testRec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコード %d, 実際のステータスコード %d", http.StatusOK, testRec.Code)
	}

	// レスポンスボディの検証
	expectedBody := `{"status":"ok"}`
	if testRec.Body.String() != expectedBody {
		t.Errorf("期待されるレスポンスボディ %s, 実際のレスポンスボディ %s", expectedBody, testRec.Body.String())
	}
}

// TestHandleAnalyze_Success は handleAnalyze ハンドラーの成功ケースのテストを行います。
func TestHandleAnalyze_Success(t *testing.T) {
	// 一時的なログファイルを作成
	// // ログの基本フォーマット: "YYYY-MM-DD HH:MM:SS [LEVEL] Message"
	tmpDir := t.TempDir()
	logFilePath := tmpDir + "/test.log"
	logFileContent := `2024-10-01 12:00:00 [INFO] アプリケーションが起動しました
2024-10-01 12:05:00 [ERROR] データベース接続に失敗しました
2024-10-01 12:10:00 [WARN] メモリ使用量が高くなっています
`

	if err := os.WriteFile(logFilePath, []byte(logFileContent), 0644); err != nil {
		t.Fatalf("一時的なログファイルの作成に失敗しました: %s", err.Error())
	}

	var reqBody jsonRequest
	reqBody.Filepath = logFilePath

	// テストリクエストの作成
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("リクエストボディの生成に失敗しました: %s", err.Error())
	}

	// リクエストの作成
	testReq := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(reqJSON))
	testRec := httptest.NewRecorder()

	// ハンドラーの呼び出し
	handleAnalyze(testRec, testReq)

	t.Logf("ステータスコード: %d", testRec.Code)
	t.Logf("レスポンスボディ: %s", testRec.Body.String())

	// ステータスコードの検証
	if testRec.Code != http.StatusOK {
		t.Errorf("期待されるステータスコード %d, 実際のステータスコード %d", http.StatusOK, testRec.Code)
	}

	// レスポンスボディの検証
	var resp jsonResponse
	respBody, err := io.ReadAll(testRec.Body)
	if err != nil {
		t.Fatalf("レスポンスボディの読み取りに失敗しました: %s", err.Error())
	}

	if err := json.Unmarshal(respBody, &resp); err != nil {
		t.Fatalf("レスポンスボディの解析に失敗しました: %s", err.Error())
	}

	if resp.Status != "ok" {
		t.Errorf("期待されるステータス 'ok', 実際のステータス %s", resp.Status)
	}

	// 期待されるログ解析結果の検証
	expectedTotalCount := 3
	expectedInfoCount := 1
	expectedWarnCount := 1
	expectedErrorCount := 1
	exceptedFirstTimestamp := "2024-10-01T12:00:00Z"
	exceptedLastTimestamp := "2024-10-01T12:10:00Z"
	exceptedFirstTimestampString := "2024-10-01 12:00:00 +0000 UTC"
	exceptedLastTimestampString := "2024-10-01 12:10:00 +0000 UTC"

	statsSource, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("レスポンスデータの型が不正です")
	}

	var stats models.Stats
	statsBytes, err := json.Marshal(statsSource)
	if err != nil {
		t.Fatalf("レスポンスデータのマーシャリングに失敗しました: %s", err.Error())
	}

	if err := json.Unmarshal(statsBytes, &stats); err != nil {
		t.Fatalf("レスポンスデータのアンマーシャリングに失敗しました: %s", err.Error())
	}

	if stats.TotalCount != expectedTotalCount {
		t.Errorf("期待される総ログ数 %d, 実際の総ログ数 %d", expectedTotalCount, stats.TotalCount)
	}

	if stats.InfoCount != expectedInfoCount {
		t.Errorf("期待されるINFOログ数 %d, 実際のINFOログ数 %d", expectedInfoCount, stats.InfoCount)
	}

	if stats.WarnCount != expectedWarnCount {
		t.Errorf("期待されるWARNログ数 %d, 実際のWARNログ数 %d", expectedWarnCount, stats.WarnCount)
	}

	if stats.ErrorCount != expectedErrorCount {
		t.Errorf("期待されるERRORログ数 %d, 実際のERRORログ数 %d", expectedErrorCount, stats.ErrorCount)
	}

	if stats.FirstTimestamp.String() != exceptedFirstTimestampString {
		t.Errorf("期待される最初のタイムスタンプ %s, 実際の最初のタイムスタンプ %s", exceptedFirstTimestamp, stats.FirstTimestamp.String())
	}

	if stats.LastTimestamp.String() != exceptedLastTimestampString {
		t.Errorf("期待される最後のタイムスタンプ %s, 実際の最後のタイムスタンプ %s", exceptedLastTimestamp, stats.LastTimestamp.String())
	}
}

// TestHandleAnalyze_InvalidJSON は handleAnalyze ハンドラーの無効なJSONリクエストのテストを行います。
func TestHandleAnalyze_InvalidJSON(t *testing.T) {
	// 無効なJSONリクエストの作成
	invalidJSON := `{"filepath": "/path/to/logfile.log",}` // 末尾のカンマが無効

	// リクエストの作成
	testReq := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBufferString(invalidJSON))
	testRec := httptest.NewRecorder()

	// ハンドラーの呼び出し
	handleAnalyze(testRec, testReq)

	t.Logf("ステータスコード: %d", testRec.Code)
	t.Logf("レスポンスボディ: %s", testRec.Body.String())

	// ステータスコードの検証
	if testRec.Code != http.StatusBadRequest {
		t.Errorf("期待されるステータスコード %d, 実際のステータスコード %d", http.StatusBadRequest, testRec.Code)
	}

	// レスポンスボディの検証
	// expectedBody := `{"status":"error","data":"リクエストの解析に失敗しました: invalid character '}' looking for beginning of object key string"}`
	// if testRec.Body.String() != expectedBody {
	// 	t.Errorf("期待されるレスポンスボディ %s, 実際のレスポンスボディ %s", expectedBody, testRec.Body.String())
	// }

	// エラーメッセージの解析
	responseBody := testRec.Body.String()
	var resp jsonResponse
	if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
		t.Fatalf("レスポンスボディの解析に失敗しました: %s", err.Error())
	}

	errorMessage, ok := resp.Data.(string)
	if !ok {
		t.Fatalf("レスポンスデータの型が不正です")
	}

	// エラーメッセージの内容の検証
	expectedErrorMessage := `リクエストの解析に失敗しました: invalid character '}' looking for beginning of object key string`
	if errorMessage != expectedErrorMessage {
		t.Errorf("期待されるエラーメッセージ %s, 実際のエラーメッセージ %s", expectedErrorMessage, errorMessage)
	}
}

// TestHandleAnalyze_FileNotFound は handleAnalyze ハンドラーのファイル未発見エラーテストを行います。
func TestHandleAnalyze_FileNotFound(t *testing.T) {
	// 存在しないログファイルパスを指定したリクエストの作成
	var reqBody jsonRequest
	reqBody.Filepath = "/non/existent/logfile.log"

	// リクエストの作成
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		t.Fatalf("リクエストボディの生成に失敗しました: %s", err.Error())
	}

	testReq := httptest.NewRequest(http.MethodPost, "/analyze", bytes.NewBuffer(reqJSON))
	testRec := httptest.NewRecorder()

	// ハンドラーの呼び出し
	handleAnalyze(testRec, testReq)

	t.Logf("ステータスコード: %d", testRec.Code)
	t.Logf("レスポンスボディ: %s", testRec.Body.String())

	// ステータスコードの検証
	if testRec.Code != http.StatusInternalServerError {
		t.Errorf("期待されるステータスコード %d, 実際のステータスコード %d", http.StatusInternalServerError, testRec.Code)
	}

	// レスポンスボディの解析
	responseBody := testRec.Body.String()
	var resp jsonResponse
	if err := json.Unmarshal([]byte(responseBody), &resp); err != nil {
		t.Fatalf("レスポンスボディの解析に失敗しました: %s", err.Error())
	}

	// エラーメッセージの内容の検証
	errorMessage, ok := resp.Data.(string)
	if !ok {
		t.Fatalf("レスポンスデータの型が不正です")
	}

	expectedErrorMessage := "ログファイルの解析に失敗しました: open /non/existent/logfile.log: The system cannot find the path specified."
	if errorMessage != expectedErrorMessage {
		t.Errorf("期待されるエラーメッセージ %s, 実際のエラーメッセージ %s", expectedErrorMessage, errorMessage)
	}
}
