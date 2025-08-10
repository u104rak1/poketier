package log

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// extractErrorDetails エラーチェーンを辿って詳細情報を抽出
func extractErrorDetails(err error) []any {
	var logArgs []any

	// エラーチェーンを辿る
	var errorChain []string
	for e := err; e != nil; e = errors.Unwrap(e) {
		errorChain = append(errorChain, e.Error())
	}

	if len(errorChain) > 1 {
		logArgs = append(logArgs, "error_chain", errorChain)
		logArgs = append(logArgs, "root_cause", errorChain[len(errorChain)-1])
	}

	// エラーの型情報も追加
	logArgs = append(logArgs, "error_type", fmt.Sprintf("%T", err))

	return logArgs
}

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

		// エラーがある場合は詳細情報を追加
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				errorDetails := extractErrorDetails(ginErr.Err)
				logArgs = append(logArgs, "gin_error_type", ginErr.Type)
				logArgs = append(logArgs, errorDetails...)
			}
		}

		var (
			clientErrorStatusCode = 400
			serverErrorStatusCode = 500
		)
		switch {
		case statusCode >= clientErrorStatusCode && statusCode < serverErrorStatusCode:
			// 4xx: クライアントエラー
			requestLogger.Warn("Request completed with client error", logArgs...)
		case statusCode >= serverErrorStatusCode:
			// 5xx: サーバーエラー
			requestLogger.Error("Request completed with server error", logArgs...)
		default:
			// その他のステータスコード（念のため）
			requestLogger.Info("Request completed", logArgs...)
		}
	}
}
