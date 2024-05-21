package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/quic-go/quic-go/http3"
)

func main() {
	// コマンドライン引数を取得
	args := os.Args

	// 引数の数を確認
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <filepath>")
		os.Exit(1)
	}

	// ファイル名を取得
	filePath := args[1]

	// ファイルを開く
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// フォームデータにファイルを追加
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		fmt.Println("Error creating form file:", err)
		return
	}

	// ファイルの内容をコピー
	_, err = io.Copy(part, file)
	if err != nil {
		fmt.Println("Error copying file:", err)
		return
	}
	writer.Close()

	// POSTリクエストを作成
	req, err := http.NewRequest("POST", "https://localhost:4433", &buf)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// クライアントTLS設定
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // サーバー証明書の検証をスキップ (テスト目的)
	}

	// HTTP/3クライアントを作成
	client := &http.Client{
		Transport: &http3.RoundTripper{
			TLSClientConfig: tlsConfig,
		},
	}

	// リクエストを送信
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// サーバーからのレスポンスを読み取る
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// サーバーのレスポンスを表示
	fmt.Println("Response from server:", string(body))
}
