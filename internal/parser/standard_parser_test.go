package parser

import (
	"errors"
	"testing"
	"time"
)

// TestStandardParser_Parser_Success は StandardParser の Parse メソッドが正しく動作することを確認します。
func TestStandardParser_Parser_Success(t *testing.T) {
	// テスト用のログ行
	logLine := "2024-06-15 14:23:45 [INFO] Application started successfully."

	// StandardParser のインスタンスを作成
	parser := NewStandardParser()

	t.Logf("パーサーのインスタンスが作成されました: %T", parser)

	// Parse メソッドを呼び出し
	entry, err := parser.Parse(logLine)
	if err != nil {
		t.Fatalf("Parse メソッドでエラーが発生しました: %v", err)
	}

	t.Logf("解析結果: %+v", entry)

	// 期待される結果と比較
	expectedTimestamp := "2024-06-15 14:23:45"
	expectedLevel := "INFO"
	expectedMessage := "Application started successfully."

	if entry.Timestamp.Format("2006-01-02 15:04:05") != expectedTimestamp {
		t.Errorf("Timestamp が期待値と異なります。期待: %s, 実際: %s", expectedTimestamp, entry.Timestamp.Format("2006-01-02 15:04:05"))
	}

	if entry.Level != expectedLevel {
		t.Errorf("Level が期待値と異なります。期待: %s, 実際: %s", expectedLevel, entry.Level)
	}

	if entry.Message != expectedMessage {
		t.Errorf("Message が期待値と異なります。期待: %s, 実際: %s", expectedMessage, entry.Message)
	}
}

// TestStandardParser_Parse_InvalidFormat は StandardParser の Parse メソッドが不正なフォーマットのログ行に対してエラーを返すことを確認します。
func TestStandardParser_Parse_InvalidFormat(t *testing.T) {
	// 不正なフォーマットのログ行
	logLine := "Invalid log line format"

	// StandardParser のインスタンスを作成
	parser := NewStandardParser()

	t.Logf("パーサーのインスタンスが作成されました: %T", parser)

	// Parse メソッドを呼び出し
	_, err := parser.Parse(logLine)
	if err == nil {
		t.Fatalf("不正なフォーマットのログ行に対してエラーが発生することを期待しましたが、エラーはありませんでした")
	}

	t.Logf("期待通りエラーが発生しました: %v", err)

	var parseErr *time.ParseError
	if !errors.As(err, &parseErr) {
		t.Errorf("予期しないエラーが発生しました: %v", err)
		return
	}
}

// TestStandardParser_Parse_EmptyLine は StandardParser の Parse メソッドが空のログ行に対してエラーを返すことを確認します。
func TestStandardParser_Parse_EmptyLine(t *testing.T) {
	// 空のログ行
	logLine := ""

	// StandardParser のインスタンスを作成
	parser := NewStandardParser()

	t.Logf("パーサーのインスタンスが作成されました: %T", parser)

	// Parse メソッドを呼び出し
	_, err := parser.Parse(logLine)
	if err == nil {
		t.Fatalf("空のログ行に対してエラーが発生することを期待しましたが、エラーはありませんでした")
	}

	t.Logf("期待通りエラーが発生しました: %v", err)

	// エラーメッセージの確認
	expectedErrMsg := "空のログ行は解析できません"
	if err.Error() != expectedErrMsg {
		t.Errorf("予期しないエラーメッセージが発生しました。期待: %s, 実際: %s", expectedErrMsg, err.Error())
	}
}

// TestStandardParser_Parse_DifferentLevels は StandardParser の Parse メソッドが異なるログレベルに対して正しく動作することを確認します。
func TestStandardParser_Parse_DifferentLevels(t *testing.T) {
	// テスト用のログ行と期待されるレベル
	invalidLevelLogLine := "2024-06-15 14:23:45 [DEBUG] This is a debug message."

	// StandardParser のインスタンスを作成
	parser := NewStandardParser()

	t.Logf("パーサーのインスタンスが作成されました: %T", parser)

	// Parse メソッドを呼び出し
	_, err := parser.Parse(invalidLevelLogLine)
	if err == nil {
		t.Fatalf("不明なログレベルに対してエラーが発生することを期待しましたが、エラーはありませんでした")
	}

	t.Logf("期待通りエラーが発生しました: %v", err)

	expectedErrMsg := "不明なログレベル: DEBUG"
	if err.Error() != expectedErrMsg {
		t.Errorf("予期しないエラーメッセージが発生しました。期待: %s, 実際: %s", expectedErrMsg, err.Error())
	}
}
