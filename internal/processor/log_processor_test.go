package processor

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

// TestLogProcessor_ProcessFile_Success は LogProcessor の ProcessFile メソッドの正常系をテストします。
func TestLogProcessor_ProcessFile_Success(t *testing.T) {
	// テスト用の一時的なログファイルを作成
	// ログの基本フォーマット: "YYYY-MM-DD HH:MM:SS [LEVEL] Message"
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_log_processor.log")
	logContent := `2024-06-01 12:00:00 [INFO] アプリケーションが起動しました。
2024-06-01 12:05:00 [ERROR] データベース接続に失敗しました。
2024-06-01 12:10:00 [WARN] メモリ使用量が高くなっています。
`

	err := os.WriteFile(tmpFile, []byte(logContent), 0644)
	if err != nil {
		t.Fatalf("一時ログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	// LogProcessor のインスタンスを作成
	lp := NewLogProcessor()

	// ProcessFile メソッドを呼び出し
	stats, err := lp.ProcessFile(tmpFile)
	if err != nil {
		t.Errorf("ProcessFile メソッドの実行に失敗しました: %v", err)
	}

	t.Logf("ProcessFile メソッドの実行に成功しました。取得した統計情報: %+v", stats)

	// 期待される統計情報を定義
	expectedTotalEntries := 3
	expectedErrorEntries := 1
	expectedWarnEntries := 1

	// 統計情報の検証
	if stats.TotalCount != expectedTotalEntries {
		t.Errorf("期待される総エントリ数 %d, 実際の総エントリ数 %d", expectedTotalEntries, stats.TotalCount)
	}

	if stats.ErrorCount != expectedErrorEntries {
		t.Errorf("期待されるエラーエントリ数 %d, 実際のエラーエントリ数 %d", expectedErrorEntries, stats.ErrorCount)
	}

	if stats.WarnCount != expectedWarnEntries {
		t.Errorf("期待される警告エントリ数 %d, 実際の警告エントリ数 %d", expectedWarnEntries, stats.WarnCount)
	}
}

// TestLogProcessor_ProcessFile_FileNotFound は LogProcessor の ProcessFile メソッドの異常系をテストします。
func TestLogProcessor_ProcessFile_FileNotFound(t *testing.T) {
	// 存在しないファイルパスを指定
	nonExistentFile := "non_existent_log_file.log"

	// LogProcessor のインスタンスを作成
	lp := NewLogProcessor()

	// ProcessFile メソッドを呼び出し
	_, err := lp.ProcessFile(nonExistentFile)
	if err == nil {
		t.Fatalf("存在しないファイルに対してエラーが発生しませんでした")
	}

	t.Logf("期待されるエラーが発生しました: %v", err)

	// エラーメッセージの内容を確認
	var pathErr *os.PathError
	if !errors.As(err, &pathErr) {
		t.Errorf("期待されるエラータイプ os.PathError ではありません: %v", err)
	}
}

// TestLogProcessor_ProcessFile_EmptyFile は LogProcessor の ProcessFile メソッドが空のファイルを正しく処理できるかをテストします。
func TestLogProcessor_ProcessFile_EmptyFile(t *testing.T) {
	// テスト用の一時的な空のログファイルを作成
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty_log_processor.log")
	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("一時ログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	// LogProcessor のインスタンスを作成
	lp := NewLogProcessor()

	// ProcessFile メソッドを呼び出し
	stats, err := lp.ProcessFile(tmpFile)
	if err != nil {
		t.Fatalf("ProcessFile メソッドの実行に失敗しました: %v", err)
	}

	t.Logf("ProcessFile メソッドの実行に成功しました。取得した統計情報: %+v", stats)

	// 統計情報の検証
	if stats.TotalCount != 0 {
		t.Errorf("期待される総エントリ数 0, 実際の総エントリ数 %d", stats.TotalCount)
	}

	if stats.ErrorCount != 0 {
		t.Errorf("期待されるエラーエントリ数 0, 実際のエラーエントリ数 %d", stats.ErrorCount)
	}

	if stats.WarnCount != 0 {
		t.Errorf("期待される警告エントリ数 0, 実際の警告エントリ数 %d", stats.WarnCount)
	}
}
