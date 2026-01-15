package reader

/*
 * bufio パッケージはバッファ付きの入出力を提供します。
 * io パッケージは基本的な入出力インターフェースを提供します。
 * os パッケージはOSの機能（ファイル操作など）を提供します。
 */
import (
	"bufio"
	"io"
	"os"
)

// FileReader はファイルからログを読み込むための構造体です。
type FileReader struct {
	// 読み込むファイルのパス
	filepath string
}

// NewFileReader は指定されたファイルパスでFileReaderを初期化します。
func NewFileReader(filepath string) *FileReader {
	return &FileReader{filepath: filepath}
}

// ReadLine はファイルから1行を読み込み、その行を文字列として返します。
func (fr *FileReader) ReadLine() (string, error) {
	// ファイルを開く
	file, err := os.Open(fr.filepath)
	if err != nil {
		return "", err
	}

	// 関数終了時にファイルを閉じる
	defer file.Close()

	// スキャンナーを使用してファイルを行ごとに読み込む
	scanner := bufio.NewScanner(file)

	// 最初の行を読み込む
	if scanner.Scan() {
		return scanner.Text(), nil
	}

	// 読み込み中にエラーが発生した場合はそれを返す
	if err := scanner.Err(); err != nil && err != io.EOF {
		return "", err
	}

	return "", nil
}
