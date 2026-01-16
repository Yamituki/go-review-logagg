package main

import (
	"log"

	"github.com/Yamituki/go-review-logagg/internal/server"
)

func main() {
	srv := server.NewServer(":8080")
	srv.SetupRoutes()
	log.Println("サーバーを起動しています: http://localhost:8080")
	if err := srv.Start(); err != nil {
		log.Fatalf("サーバーの起動に失敗: %v", err)
	}
}
