# go-review-logagg

Go言語復習用のログ集約ツール（標準ライブラリのみ使用）

## 概要
複数のログファイルを監視・解析し、HTTP APIで結果を提供するツールです。

## 開発状況
v1.0.0 - 全機能実装完了

## 実装予定機能
- [x] ログファイル読み込み
- [x] ログパース機能
- [x] データ集約機能
- [x] 統合処理（ファイル処理パイプライン）
- [x] HTTP API
- [x] 並行処理
- [x] リアルタイム監視

## 開発方針
- Git Flowを使用
- 標準ライブラリのみ使用
- テスト駆動開発

## 使い方

### サーバー起動
```bash
go run cmd/logagg/main.go
```

### API使用例
```bash
# ヘルスチェック
curl http://localhost:8080/health

# ログ解析
curl -X POST http://localhost:8080/analyze \
  -H "Content-Type: application/json" \
  -d '{"filepath": "sample.log"}'
```

## 制限事項
- リアルタイム監視はファイル全体を再読み込みするため、非常に大きなファイル（数GB以上）には向きません