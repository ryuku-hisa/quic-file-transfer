package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

func main() {
	// 証明書と秘密鍵のパスを定義
	const certPath = "tls/server.crt"
	const keyPath = "tls/server.key"

	// 証明書ファイルを読み込み
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		log.Fatalf("failed to load key pair: %v", err)
	}

	// TLS設定を初期化
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		NextProtos:   []string{"h3-29"},
	}

	// HTTPハンドラを設定
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// POSTメソッド以外はエラーを返す
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
			return
		}

		// ファイルをフォームデータから取得
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, "Unable to read file", http.StatusBadRequest)
			return
		}
		defer file.Close()

		// 受け取ったファイルを保存
		outFile, err := os.Create("received_file")
		if err != nil {
			http.Error(w, "Unable to create file", http.StatusInternalServerError)
			return
		}
		defer outFile.Close()

		// ファイルをコピー
		if _, err = io.Copy(outFile, file); err != nil {
			http.Error(w, "Unable to save file", http.StatusInternalServerError)
			return
		}

		// 成功メッセージを送信
		fmt.Fprintf(w, "File received successfully")
	})

	// QUICプロトコルの設定
	quicConfig := &quic.Config{
		MaxIdleTimeout:       30 * time.Second, // 最大アイドルタイムアウト
		HandshakeIdleTimeout: 10 * time.Second, // ハンドシェイクのアイドルタイムアウト
	}

	// HTTP/3サーバーを設定
	server := http3.Server{
		Addr:       ":4433",    // サーバーのアドレス
		QUICConfig: quicConfig, // QUIC設定
		TLSConfig:  tlsConfig,  // TLS設定
	}

	// サーバーを起動
	fmt.Println("Starting server on :4433")
	if err := server.ListenAndServeTLS(certPath, keyPath); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
