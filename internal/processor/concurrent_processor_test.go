package processor

import (
	"fmt"
	"os"
	"testing"
	"time"
)

// TestConcurrentProcessor_ProcessFiles_Success は ConcurrentProcessor の ProcessFiles メソッドの成功ケースをテストします。
func TestConcurrentProcessor_ProcessFiles_Success(t *testing.T) {
	// 一時的ログファイルを複数作成
	// ログの基本フォーマット: "YYYY-MM-DD HH:MM:SS [LEVEL] Message"
	tmpDir := t.TempDir()
	var filePaths []string
	var firstTimestamp time.Time
	var lastTimestamp time.Time
	for i := 0; i < 3; i++ {
		filePath := fmt.Sprintf("%s/logfile_%d.log", tmpDir, i)
		now := time.Now()
		infoTime := now.Add(time.Duration(i*15) * time.Minute)
		errorTime := now.Add(time.Duration(i*15+5) * time.Minute)
		warnTime := now.Add(time.Duration(i*15+10) * time.Minute)

		content := fmt.Sprintf(`%s [INFO] アプリケーションが起動しました。
%s [ERROR] データベース接続に失敗しました。
%s [WARN] メモリ使用量が高いです。
`, infoTime.Format("2006-01-02 15:04:05"), errorTime.Format("2006-01-02 15:04:05"), warnTime.Format("2006-01-02 15:04:05"))

		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("一時ログファイルの作成に失敗しました: %v", err)
		}

		// タイムスタンプの初期化
		if i == 0 {
			firstTimestamp = infoTime
			lastTimestamp = warnTime
		}

		// 最初のタイムスタンプを更新
		if infoTime.Before(firstTimestamp) {
			firstTimestamp = infoTime
		}

		// 最後のタイムスタンプを更新
		if warnTime.After(lastTimestamp) {
			lastTimestamp = warnTime
		}

		filePaths = append(filePaths, filePath)
	}

	// ConcurrentProcessor の初期化
	cp := NewConcurrentProcessor(2)

	// ファイルの処理
	stats, err := cp.ProcessFiles(filePaths)
	if err != nil {
		t.Fatalf("ProcessFiles メソッドがエラーを返しました: %v", err)
	}

	t.Logf("集約結果: %+v", stats)

	// 予期のタイムスタンプの計算
	expFirstTimestamp, err := time.Parse("2006-01-02 15:04:05", firstTimestamp.Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Fatalf("タイムスタンプのパースに失敗しました: %v", err)
	}

	expLastTimestamp, err := time.Parse("2006-01-02 15:04:05", lastTimestamp.Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Fatalf("タイムスタンプのパースに失敗しました: %v", err)
	}

	// 結果の検証
	expectedTotalEntries := 9 // 各ファイルに3エントリ、3ファイルで合計9エントリ
	expectedErrorEntries := 3 // 各ファイルに1つのERRORエントリ、3ファイルで合計3エントリ
	expectedWarnEntries := 3  // 各ファイルに1つのWARNエントリ、3ファイルで合計3エントリ
	expectedInfoEntries := 3  // 各ファイルに1つのINFOエントリ、3ファイルで合計3エントリ
	expectedFirstTimestamp := expFirstTimestamp
	expectedLastTimestamp := expLastTimestamp

	if stats.TotalCount != expectedTotalEntries {
		t.Errorf("TotalCount が期待値と異なります。期待値: %d, 実際: %d", expectedTotalEntries, stats.TotalCount)
	}

	if stats.ErrorCount != expectedErrorEntries {
		t.Errorf("ErrorCount が期待値と異なります。期待値: %d, 実際: %d", expectedErrorEntries, stats.ErrorCount)
	}

	if stats.WarnCount != expectedWarnEntries {
		t.Errorf("WarnCount が期待値と異なります。期待値: %d, 実際: %d", expectedWarnEntries, stats.WarnCount)
	}

	if stats.InfoCount != expectedInfoEntries {
		t.Errorf("InfoCount が期待値と異なります。期待値: %d, 実際: %d", expectedInfoEntries, stats.InfoCount)
	}

	if !stats.FirstTimestamp.Equal(expectedFirstTimestamp) {
		t.Errorf("FirstTimestamp が期待値と異なります。期待値: %v, 実際: %v", expectedFirstTimestamp, stats.FirstTimestamp)
	}

	if !stats.LastTimestamp.Equal(expectedLastTimestamp) {
		t.Errorf("LastTimestamp が期待値と異なります。期待値: %v, 実際: %v", expectedLastTimestamp, stats.LastTimestamp)
	}
}

// TestConcurrentProcessor_ProcessFiles_EmptyFiles は ConcurrentProcessor の ProcessFiles メソッドが空のログファイルを正しく処理できることをテストします。
func TestConcurrentProcessor_ProcessFiles_EmptyFiles(t *testing.T) {
	// 一時的な空のログファイルリストを作成
	var filePaths []string

	// ConcurrentProcessor の初期化
	cp := NewConcurrentProcessor(2)

	// ファイルの処理
	stats, err := cp.ProcessFiles(filePaths)
	if err != nil {
		t.Fatalf("ProcessFiles メソッドがエラーを返しました: %v", err)
	}

	t.Logf("集約結果: %+v", stats)

	// 結果の検証
	if stats.TotalCount != 0 {
		t.Fatalf("TotalCount が期待値と異なります。期待値: 0, 実際: %d", stats.TotalCount)
	}

	if stats.ErrorCount != 0 {
		t.Fatalf("ErrorCount が期待値と異なります。期待値: 0, 実際: %d", stats.ErrorCount)
	}

	if stats.WarnCount != 0 {
		t.Fatalf("WarnCount が期待値と異なります。期待値: 0, 実際: %d", stats.WarnCount)
	}

	if stats.InfoCount != 0 {
		t.Fatalf("InfoCount が期待値と異なります。期待値: 0, 実際: %d", stats.InfoCount)
	}

	if !stats.FirstTimestamp.IsZero() {
		t.Fatalf("FirstTimestamp が期待値と異なります。期待値: zero value, 実際: %v", stats.FirstTimestamp)
	}

	if !stats.LastTimestamp.IsZero() {
		t.Fatalf("LastTimestamp が期待値と異なります。期待値: zero value, 実際: %v", stats.LastTimestamp)
	}
}

// TestConcurrentProcessor_ProcessFiles_PartialErrors は ConcurrentProcessor の ProcessFiles メソッドが一部のファイルでエラーが発生した場合でも正しく処理できることをテストします。
func TestConcurrentProcessor_ProcessFiles_PartialErrors(t *testing.T) {
	// 一時的ログファイルを複数作成（1つは存在しないファイルパスを指定）
	tmpDir := t.TempDir()
	var filePaths []string

	infoTimestamp := time.Date(2024, 1, 1, 12, 0, 0, 0, time.Local)
	errorTimestamp := time.Date(2024, 1, 1, 12, 5, 0, 0, time.Local)

	// 正常なログファイルの作成
	filePath1 := fmt.Sprintf("%s/logfile_valid.log", tmpDir)
	content := fmt.Sprintf(`%s [INFO] アプリケーションが起動しました。
%s [ERROR] データベース接続に失敗しました。
`, infoTimestamp.Format("2006-01-02 15:04:05"), errorTimestamp.Format("2006-01-02 15:04:05"))

	if err := os.WriteFile(filePath1, []byte(content), 0644); err != nil {
		t.Fatalf("一時ログファイルの作成に失敗しました: %v", err)
	}

	// 存在しないログファイルのパス
	filePath2 := fmt.Sprintf("%s/logfile_nonexistent.log", tmpDir)

	filePaths = append(filePaths, filePath1, filePath2)

	// ConcurrentProcessor の初期化
	cp := NewConcurrentProcessor(2)

	// ファイルの処理
	stats, err := cp.ProcessFiles(filePaths)
	if err == nil {
		t.Fatalf("ProcessFiles メソッドがエラーを返すことを期待していましたが、nil が返されました")
	}

	t.Logf("集約結果: %+v", stats)

	// エラーメッセージの検証
	expectedErrorMsg := fmt.Sprintf("ファイルの読み込みに失敗しました: open %s: The system cannot find the file specified.", filePath2)
	if err.Error() != expectedErrorMsg {
		t.Fatalf("エラーメッセージが期待値と異なります。期待値: %s, 実際: %s", expectedErrorMsg, err.Error())
	}

	// 結果の検証
	expectedTotalEntries := 2 // 有効なファイルに2エントリ
	expectedErrorEntries := 1 // 有効なファイルに1つのERRORエントリ

	if stats.TotalCount != expectedTotalEntries {
		t.Fatalf("TotalCount が期待値と異なります。期待値: %d, 実際: %d", expectedTotalEntries, stats.TotalCount)
	}

	if stats.ErrorCount != expectedErrorEntries {
		t.Fatalf("ErrorCount が期待値と異なります。期待値: %d, 実際: %d", expectedErrorEntries, stats.ErrorCount)
	}

	if stats.WarnCount != 0 {
		t.Fatalf("WarnCount が期待値と異なります。期待値: 0, 実際: %d", stats.WarnCount)
	}

	if stats.InfoCount != 1 {
		t.Fatalf("InfoCount が期待値と異なります。期待値: 1, 実際: %d", stats.InfoCount)
	}

	expectedFirstTimestamp, err := time.Parse("2006-01-02 15:04:05", infoTimestamp.Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Fatalf("タイムスタンプのパースに失敗しました: %v", err)
	}

	expectedLastTimestamp, err := time.Parse("2006-01-02 15:04:05", errorTimestamp.Format("2006-01-02 15:04:05"))
	if err != nil {
		t.Fatalf("タイムスタンプのパースに失敗しました: %v", err)
	}

	if !stats.FirstTimestamp.Equal(expectedFirstTimestamp) {
		t.Fatalf("FirstTimestamp が期待値と異なります。期待値: %v, 実際: %v", expectedFirstTimestamp, stats.FirstTimestamp)
	}

	if !stats.LastTimestamp.Equal(expectedLastTimestamp) {
		t.Fatalf("LastTimestamp が期待値と異なります。期待値: %v, 実際: %v", expectedLastTimestamp, stats.LastTimestamp)
	}

}
