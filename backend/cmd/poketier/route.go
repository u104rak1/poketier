package main

import (
	"poketier/env"
	"poketier/pkg/log"

	"github.com/gin-gonic/gin"
)

func startServer() {
	// 環境変数を読み込み
	envConfig := env.NewEnv()

	r := gin.Default()

	// ログミドルウェアを設定
	r.Use(log.NewMiddleware(envConfig.LOG_LEVEL, envConfig.IS_SILENT_LOG))

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// サーバー起動
	startupLogger := log.NewStartupLogger(envConfig.LOG_LEVEL, envConfig.IS_SILENT_LOG)
	startupLogger.Info("Starting server", "port", envConfig.APP_PORT)
	if err := r.Run(":" + envConfig.APP_PORT); err != nil {
		startupLogger.Error("Failed to start server", "error", err)
	}
}
