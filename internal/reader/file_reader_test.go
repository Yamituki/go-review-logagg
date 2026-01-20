package reader

import (
	"os"
	"path/filepath"
	"testing"
)

// TestFileReader_ReadLine_Success は FileReader の ReadLine メソッドの成功ケースをテストします。
func TestFileReader_ReadLine_Success(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_log.txt")
	content := "これはテストログの1行目です。\nこれは2行目です。"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("テスト用ログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	t.Logf("一時ログファイルが作成されました: %s", tmpFile)

	// FileReader を初期化
	fr := NewFileReader(tmpFile)

	// ReadLine を呼び出し
	line, err := fr.ReadLine()
	if err != nil {
		t.Fatalf("ReadLineに失敗しました: %v", err)
	}

	t.Logf("読み込んだ行: %s", line)

	// 期待される行と比較
	expected := "これはテストログの1行目です。"

	if line != expected {
		t.Errorf("期待される行と異なります。期待: '%s', 実際: '%s'", expected, line)
	}
}

// TestFileReader_ReadLine_FileNotFound は FileReader の ReadLine メソッドでファイルが存在しない場合のエラー処理をテストします。
func TestFileReader_ReadLine_FileNotFound(t *testing.T) {
	// 存在しないファイルパスを指定
	nonExistentFile := "non_existent_log.txt"

	// FileReader を初期化
	fr := NewFileReader(nonExistentFile)

	// ReadLine を呼び出し
	_, err := fr.ReadLine()
	if err == nil {
		t.Fatalf("存在しないファイルに対してエラーが発生することを期待しましたが、エラーはありませんでした")
	}

	t.Logf("エラーが発生しました: %v", err)

	// エラーメッセージを確認
	if !os.IsNotExist(err) {
		t.Errorf("ファイルが存在しないエラーを期待しましたが、実際のエラー: %v", err)
	}

}

// TestFileReader_ReadLine_EmptyFile は FileReader の ReadLine メソッドで空ファイルの処理をテストします。
func TestFileReader_ReadLine_EmptyFile(t *testing.T) {
	// テスト用の空ファイルを作成
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty_log.txt")
	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("空のログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	t.Logf("空の一時ログファイルが作成されました: %s", tmpFile)

	// FileReader を初期化
	fr := NewFileReader(tmpFile)

	// ReadLine を呼び出し
	line, err := fr.ReadLine()
	if err != nil {
		t.Fatalf("ReadLineに失敗しました: %v", err)
	}

	t.Logf("読み込んだ行: %s", line)

	// 空ファイルなので空文字列が返ることを確認
	if line != "" {
		t.Errorf("空ファイルの場合は空文字列が返ることを期待しましたが、実際は: '%s'", line)
	}

}

// TestFileReader_ReadAllLines_Success は FileReader の ReadAllLines メソッドの成功ケースをテストします。
func TestFileReader_ReadAllLines_Success(t *testing.T) {
	// テスト用の一時ファイルを作成
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_log_all.txt")
	content := "行1\n行2\n行3"
	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("テスト用ログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	t.Logf("一時ログファイルが作成されました: %s", tmpFile)

	// FileReader を初期化
	fr := NewFileReader(tmpFile)

	// ReadAllLines を呼び出し
	lines, err := fr.ReadAllLines()
	if err != nil {
		t.Fatalf("ReadAllLinesに失敗しました: %v", err)
	}

	t.Logf("読み込んだ行: %v", lines)

	// 期待される行と比較
	expected := []string{"行1", "行2", "行3"}

	if len(lines) != len(expected) {
		t.Fatalf("期待される行数と異なります。期待: %d, 実際: %d", len(expected), len(lines))
	}

	for i, line := range lines {
		if line != expected[i] {
			t.Errorf("行 %d が期待と異なります。期待: '%s', 実際: '%s'", i+1, expected[i], line)
		}
	}
}

// TestFileReader_ReadAllLines_FileNotFound は FileReader の ReadAllLines メソッドでファイルが存在しない場合のエラー処理をテストします。
func TestFileReader_ReadAllLines_FileNotFound(t *testing.T) {
	// 存在しないファイルパスを指定
	nonExistentFile := "non_existent_log_all.txt"

	// FileReader を初期化
	fr := NewFileReader(nonExistentFile)

	// ReadAllLines を呼び出し
	_, err := fr.ReadAllLines()
	if err == nil {
		t.Fatalf("存在しないファイルに対してエラーが発生することを期待しましたが、エラーはありませんでした")
	}

	t.Logf("エラーが発生しました: %v", err)

	// エラーメッセージを確認
	if !os.IsNotExist(err) {
		t.Errorf("ファイルが存在しないエラーを期待しましたが、実際のエラー: %v", err)
	}
}

// TestFileReader_ReadAllLines_EmptyFile は FileReader の ReadAllLines メソッドで空ファイルの処理をテストします。
func TestFileReader_ReadAllLines_EmptyFile(t *testing.T) {
	// テスト用の空ファイルを作成
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "empty_log_all.txt")
	err := os.WriteFile(tmpFile, []byte(""), 0644)
	if err != nil {
		t.Fatalf("空のログファイルの作成に失敗しました: %v", err)
	}

	// テスト終了後に一時ファイルを削除
	defer os.Remove(tmpFile)

	t.Logf("空の一時ログファイルが作成されました: %s", tmpFile)

	// FileReader を初期化
	fr := NewFileReader(tmpFile)

	// ReadAllLines を呼び出し
	lines, err := fr.ReadAllLines()
	if err != nil {
		t.Fatalf("ReadAllLinesに失敗しました: %v", err)
	}

	t.Logf("読み込んだ行: %v", lines)

	// 空ファイルなので空のスライスが返ることを確認
	if len(lines) != 0 {
		t.Errorf("空ファイルの場合は空のスライスが返ることを期待しましたが、実際は: %v", lines)
	}
}
