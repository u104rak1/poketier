package main

import (
	"context"
	"poketier/apps/season"
	"poketier/env"
	"poketier/pkg/cors"
	"poketier/pkg/log"
	"poketier/sqlc"
	"poketier/sqlc/db"

	"github.com/gin-gonic/gin"
)

func startServer() {
	// 環境変数を読み込み
	envConfig := env.NewEnv()

	// データベース接続プールを初期化
	pool, err := sqlc.NewPgxPool(context.Background(), envConfig)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	// Querierを作成
	queries := db.New(pool)

	r := gin.Default()

	// CORSミドルウェアを設定
	cors.SetCORS(r, envConfig.ALLOW_ORIGINS, envConfig.APP_ENV)

	// ログミドルウェアを設定
	r.Use(log.NewMiddleware(envConfig.LOG_LEVEL, envConfig.IS_SILENT_LOG))

	// ヘルスチェックエンドポイント
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	v1 := r.Group("/v1")

	// WireでDIされたハンドラーを使用
	newSeasonHandler(v1, queries)

	// サーバー起動
	startupLogger := log.NewStartupLogger(envConfig.LOG_LEVEL, envConfig.IS_SILENT_LOG)
	startupLogger.Info("Starting server", "port", envConfig.APP_PORT)
	if err := r.Run(":" + envConfig.APP_PORT); err != nil {
		startupLogger.Error("Failed to start server", "error", err)
	}
}

func newSeasonHandler(engine *gin.RouterGroup, queries *db.Queries) {
	// Wireで生成されたDIコードを使用してハンドラーを初期化
	seasonHandler := season.InitializeListSeasonsHandler(queries)

	// シーズン関連のエンドポイントを登録
	engine.GET("/seasons", seasonHandler.Handle)
}
