package config

import "os"

var JWTSecret []byte

func Init() {
	// 環境変数からJWTシークレットを取得
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "your-secret-key" // デフォルト値（本番環境では必ず変更）
	}
	JWTSecret = []byte(secret)
}