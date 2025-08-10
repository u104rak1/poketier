package log

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func NewMiddleware(logLevel string, isSilent bool) gin.HandlerFunc {
	logger := newLogger(logLevel, isSilent)

	return func(c *gin.Context) {
		start := time.Now()

		// リクエストIDを生成
		requestID := uuid.New().String()
		c.Set(RequestIDKey, requestID)

		// リクエスト情報を含むロガーを作成してコンテキストに設定
		requestLogger := logger.With(
			"request_id", requestID,
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.ClientIP(),
		)
		c.Set(LoggerKey, requestLogger)

		// リクエスト開始ログ
		requestLogger.Info("Request started")

		// 次のハンドラーを実行
		c.Next()

		// レスポンス時間を計算
		duration := time.Since(start)
		statusCode := c.Writer.Status()

		// ステータスコードに応じてログレベルを変更
		logArgs := []any{
			"status_code", statusCode,
			"duration_ms", duration.Milliseconds(),
		}

		switch {
		case statusCode >= 400 && statusCode < 500:
			// 4xx: クライアントエラー
			requestLogger.Warn("Request completed with client error", logArgs...)
		case statusCode >= 500:
			// 5xx: サーバーエラー
			requestLogger.Error("Request completed with server error", logArgs...)
		default:
			// その他のステータスコード（念のため）
			requestLogger.Info("Request completed", logArgs...)
		}
	}
}
