package aggregator

import (
	"testing"
	"time"

	"github.com/Yamituki/go-review-logagg/pkg/models"
)

// TestLogAggregator_Add_Success は LogAggregator の Add メソッドの成功ケースをテストします。
func TestLogAggregator_Add_Success(t *testing.T) {
	// テスト用のログエントリを作成
	entry := models.LogEntry{
		Timestamp: time.Now(),
		Level:     "INFO",
		Message:   "これはテストログエントリです。",
	}

	// LogAggregator のインスタンスを作成
	aggregator := NewLogAggregator()

	t.Logf("作成したインスタンス: %+v", aggregator)

	// Add メソッドを呼び出し
	err := aggregator.Add(entry)
	if err != nil {
		t.Fatalf("エラーは発生しないはずですが、エラーが発生しました: %v", err)
	}
}

// TestLogAggregator_GetStats は LogAggregator の GetStats メソッドをテストします。
func TestLogAggregator_GetStats(t *testing.T) {
	// LogAggregator のインスタンスを作成
	aggregator := NewLogAggregator()

	t.Logf("作成したインスタンス: %+v", aggregator)

	// 複数のログエントリを追加
	entries := []models.LogEntry{
		{Timestamp: time.Now().Add(-10 * time.Minute), Level: "INFO", Message: "終わった？"},
		{Timestamp: time.Now().Add(-5 * time.Minute), Level: "WARN", Message: "注意してください"},
		{Timestamp: time.Now(), Level: "ERROR", Message: "エラーが発生しました"},
	}

	// エントリを追加
	for _, entry := range entries {
		err := aggregator.Add(entry)
		if err != nil {
			t.Fatalf("エラーは発生しないはずですが、エラーが発生しました: %v", err)
		}
	}

	// 統計情報を取得
	stats := aggregator.GetStats()

	t.Logf("取得した統計情報: %+v", stats)

	// 統計情報の検証
	if stats.TotalCount != 3 {
		t.Errorf("期待される総ログ数は 3 ですが、実際の値は %d です", stats.TotalCount)
	}
}

// TestLogAggregator_Reset は LogAggregator の Reset メソッドをテストします。
func TestLogAggregator_Reset(t *testing.T) {
	// LogAggregator のインスタンスを作成
	aggregator := NewLogAggregator()

	t.Logf("作成したインスタンス: %+v", aggregator)

	// 複数のログエントリを追加
	entries := []models.LogEntry{
		{Timestamp: time.Now().Add(-10 * time.Minute), Level: "INFO", Message: "終わった？"},
		{Timestamp: time.Now().Add(-5 * time.Minute), Level: "WARN", Message: "注意してください"},
		{Timestamp: time.Now(), Level: "ERROR", Message: "エラーが発生しました"},
	}

	// エントリを追加
	for _, entry := range entries {
		err := aggregator.Add(entry)
		if err != nil {
			t.Fatalf("エラーは発生しないはずですが、エラーが発生しました: %v", err)
		}
	}

	// Reset メソッドを呼び出し
	aggregator.Reset()

	// 統計情報を取得
	stats := aggregator.GetStats()

	t.Logf("リセット後の統計情報: %+v", stats)

	// 統計情報の検証
	if stats.TotalCount != 0 {
		t.Errorf("リセット後の総ログ数は 0 であるべきですが、実際の値は %d です", stats.TotalCount)
	}
}

// TestLogAggregator_Timestamps は LogAggregator のタイムスタンプの更新をテストします。
func TestLogAggregator_Timestamps(t *testing.T) {
	// LogAggregator のインスタンスを作成
	aggregator := NewLogAggregator()

	t.Logf("作成したインスタンス: %+v", aggregator)

	// 複数のログエントリを追加
	entries := []models.LogEntry{
		{Timestamp: time.Now().Add(-15 * time.Minute), Level: "INFO", Message: "最初のログ"},
		{Timestamp: time.Now().Add(-10 * time.Minute), Level: "WARN", Message: "中間のログ"},
		{Timestamp: time.Now().Add(-5 * time.Minute), Level: "ERROR", Message: "最後のログ"},
	}

	// エントリを追加
	for _, entry := range entries {
		err := aggregator.Add(entry)
		if err != nil {
			t.Fatalf("エラーは発生しないはずですが、エラーが発生しました: %v", err)
		}
	}

	// 統計情報を取得
	stats := aggregator.GetStats()

	t.Logf("取得した統計情報: %+v", stats)

	// タイムスタンプの検証
	expectedFirst := entries[0].Timestamp
	expectedLast := entries[2].Timestamp

	if !stats.FirstTimestamp.Equal(expectedFirst) {
		t.Errorf("期待される最初のタイムスタンプは %v ですが、実際の値は %v です", expectedFirst, stats.FirstTimestamp)
	}

	if !stats.LastTimestamp.Equal(expectedLast) {
		t.Errorf("期待される最後のタイムスタンプは %v ですが、実際の値は %v です", expectedLast, stats.LastTimestamp)
	}
}
