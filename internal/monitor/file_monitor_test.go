package monitor

import (
	"os"
	"testing"
	"time"
)

// TestFileMonitor_Start_Stop は FileMonitor の Start と Stop メソッドのテストを行います。
func TestFileMonitor_Start_Stop(t *testing.T) {
	// テスト用の一時ファイルを作成
	filePath, err := fileCreator()
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// FileMonitor の初期化
	monitoringInterval := 10 * time.Millisecond
	fileMonitor := NewFileMonitor(filePath, monitoringInterval)

	// 監視の開始
	err = fileMonitor.Start()
	if err != nil {
		t.Fatalf("FileMonitor の Start に失敗: %v", err)
	}

	// 一定時間待機する
	time.Sleep(50 * time.Millisecond)

	// 監視の停止
	err = fileMonitor.Stop()
	if err != nil {
		t.Fatalf("FileMonitor の Stop に失敗: %v", err)
	}

}

// TestFileMonitor_GetStats_FileChange は FileMonitor の GetStats メソッドとファイル変更の検出をテストします。
func TestFileMonitor_GetStats_FileChange(t *testing.T) {
	// テスト用の一時ファイルを作成
	filePath, err := fileCreator()
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// FileMonitor の初期化
	monitoringInterval := 10 * time.Millisecond
	fileMonitor := NewFileMonitor(filePath, monitoringInterval)

	// 監視の開始
	err = fileMonitor.Start()
	if err != nil {
		t.Fatalf("FileMonitor の Start に失敗: %v", err)
	}

	// 一定時間待機する
	time.Sleep(monitoringInterval * 3)

	// ファイルを変更
	err = fileModifier(filePath)
	if err != nil {
		t.Fatalf("ファイルの変更に失敗: %v", err)
	}

	// 変更が反映されるまで待機
	time.Sleep(monitoringInterval * 3)
	// 統計情報の取得
	stats, err := fileMonitor.GetStats()
	if err != nil {
		t.Fatalf("FileMonitor の GetStats に失敗: %v", err)
	}

	t.Logf("取得した統計情報: %+v", stats)

	// 期待されるカウント値
	expectedInfoCount := 2  // 元の1行 + 追加の1行
	expectedWarnCount := 1  // 元の1行
	expectedErrorCount := 2 // 元の1行 + 追加の1行

	// 統計情報の検証
	if stats.InfoCount != expectedInfoCount {
		t.Errorf("InfoCount が期待値と異なる: 期待値=%d, 実際=%d", expectedInfoCount, stats.InfoCount)
	}
	if stats.WarnCount != expectedWarnCount {
		t.Errorf("WarnCount が期待値と異なる: 期待値=%d, 実際=%d", expectedWarnCount, stats.WarnCount)
	}
	if stats.ErrorCount != expectedErrorCount {
		t.Errorf("ErrorCount が期待値と異なる: 期待値=%d, 実際=%d", expectedErrorCount, stats.ErrorCount)
	}

	// 監視の停止
	err = fileMonitor.Stop()
	if err != nil {
		t.Fatalf("FileMonitor の Stop に失敗: %v", err)
	}
}

// TestFileMonitor_GetStats_NoChange は FileMonitor の GetStats メソッドでファイル変更がない場合のテストを行います。
func TestFileMonitor_GetStats_NoChange(t *testing.T) {
	// テスト用の一時ファイルを作成
	filePath, err := fileCreator()
	if err != nil {
		t.Fatalf("一時ファイルの作成に失敗: %v", err)
	}

	// FileMonitor の初期化
	monitoringInterval := 10 * time.Millisecond
	fileMonitor := NewFileMonitor(filePath, monitoringInterval)

	// 監視の開始
	err = fileMonitor.Start()
	if err != nil {
		t.Fatalf("FileMonitor の Start に失敗: %v", err)
	}

	// 一定時間待機する
	time.Sleep(50 * time.Millisecond)

	// 統計情報の取得
	stats, err := fileMonitor.GetStats()
	if err != nil {
		t.Fatalf("FileMonitor の GetStats に失敗: %v", err)
	}

	t.Logf("取得した統計情報: %+v", stats)

	// 期待されるカウント値
	expectedInfoCount := 1  // 元の1行
	expectedWarnCount := 1  // 元の1行
	expectedErrorCount := 1 // 元の1行

	// 統計情報の検証
	if stats.InfoCount != expectedInfoCount {
		t.Errorf("InfoCount が期待値と異なる: 期待値=%d, 実際=%d", expectedInfoCount, stats.InfoCount)
	}

	if stats.WarnCount != expectedWarnCount {
		t.Errorf("WarnCount が期待値と異なる: 期待値=%d, 実際=%d", expectedWarnCount, stats.WarnCount)
	}

	if stats.ErrorCount != expectedErrorCount {
		t.Errorf("ErrorCount が期待値と異なる: 期待値=%d, 実際=%d", expectedErrorCount, stats.ErrorCount)
	}

	// 監視の停止
	err = fileMonitor.Stop()
	if err != nil {
		t.Fatalf("FileMonitor の Stop に失敗: %v", err)
	}
}

// ファイル作成者
func fileCreator() (string, error) {
	// テスト用の一時ファイルを作成
	// ログの基本フォーマット: "YYYY-MM-DD HH:MM:SS [LEVEL] Message"
	tempFile, err := os.CreateTemp("", "testfile_*.log")
	if err != nil {
		return "", err
	}

	// テスト用のログデータを書き込む
	logData := `2024-01-01 12:00:00 [INFO] アプリケーションが起動しました。
2024-01-01 12:05:00 [WARN] メモリ使用量が高くなっています。
2024-01-01 12:10:00 [ERROR] データベース接続に失敗しました。`

	// ファイルにログデータを書き込む
	_, err = tempFile.WriteString(logData)
	if err != nil {
		return "", err
	}

	// 一時ファイルのパスを返す
	return tempFile.Name(), nil
}

// ファイル変更者
func fileModifier(filePath string) error {
	// 既存のファイルを開く
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	// 追加するログデータ
	additionalLogData := `
2024-01-01 12:15:00 [INFO] ユーザーがログインしました。
2024-01-01 12:20:00 [ERROR] ファイルの読み込みに失敗しました。`

	// ファイルに追加のログデータを書き込む
	_, err = file.WriteString(additionalLogData)
	if err != nil {
		return err
	}

	return nil
}
