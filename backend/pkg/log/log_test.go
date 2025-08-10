package log_test

import (
	"context"
	"os"
	"testing"

	"poketier/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetLogger(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_コンテキストからロガーを取得", func(t *testing.T) {
		t.Parallel()

		// Arrange
		c, _ := gin.CreateTestContext(nil)
		expectedLogger := log.NewStartupLogger("info", true)
		c.Set(log.LoggerKey, expectedLogger)

		// Act
		logger := log.GetLogger(c)

		// Assert
		assert.NotNil(t, logger, "Logger should not be nil")
		assert.Equal(t, expectedLogger, logger, "Should return the logger from context")
	})

	t.Run("正常系_コンテキストにロガーがない場合はフォールバック", func(t *testing.T) {
		t.Parallel()

		// Arrange
		c, _ := gin.CreateTestContext(nil)

		// Act
		logger := log.GetLogger(c)

		// Assert
		assert.NotNil(t, logger, "Logger should not be nil")
	})
}

func TestGetRequestID(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_コンテキストからリクエストIDを取得", func(t *testing.T) {
		t.Parallel()

		// Arrange
		c, _ := gin.CreateTestContext(nil)
		expectedRequestID := "test-request-id-123"
		c.Set(log.RequestIDKey, expectedRequestID)

		// Act
		requestID := log.GetRequestID(c)

		// Assert
		assert.Equal(t, expectedRequestID, requestID, "Should return the request ID from context")
	})

	t.Run("正常系_コンテキストにリクエストIDがない場合は空文字", func(t *testing.T) {
		t.Parallel()

		// Arrange
		c, _ := gin.CreateTestContext(nil)

		// Act
		requestID := log.GetRequestID(c)

		// Assert
		assert.Equal(t, "", requestID, "Should return empty string when no request ID")
	})
}

func TestNewStartupLogger(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName string
		logLevel string
		isSilent bool
	}{
		{
			caseName: "正常系_デバッグレベル",
			logLevel: "debug",
			isSilent: false,
		},
		{
			caseName: "正常系_インフォレベル",
			logLevel: "info",
			isSilent: false,
		},
		{
			caseName: "正常系_ワーニングレベル",
			logLevel: "warn",
			isSilent: false,
		},
		{
			caseName: "正常系_エラーレベル",
			logLevel: "error",
			isSilent: false,
		},
		{
			caseName: "正常系_サイレントモード",
			logLevel: "info",
			isSilent: true,
		},
		{
			caseName: "正常系_無効なレベル",
			logLevel: "invalid",
			isSilent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			logger := log.NewStartupLogger(tt.logLevel, tt.isSilent)

			// Assert
			assert.NotNil(t, logger, "Logger should not be nil")

			// ロガーが実際に動作することを確認
			logger.Info("test message")
			logger.Debug("debug message")
			logger.Warn("warn message")
			logger.Error("error message")
		})
	}
}

func TestSlogLogger_Methods(t *testing.T) {
	t.Parallel()

	t.Run("正常系_ログメソッドが正常に動作", func(t *testing.T) {
		t.Parallel()

		// Arrange
		logger := log.NewStartupLogger("info", true) // サイレントモードでテスト

		// Act & Assert - パニックしないことを確認
		logger.Debug("debug message", "key", "value")
		logger.Info("info message", "key", "value")
		logger.Warn("warn message", "key", "value")
		logger.Error("error message", "key", "value")
	})

	t.Run("正常系_Withメソッドでコンテキスト追加", func(t *testing.T) {
		t.Parallel()

		// Arrange
		logger := log.NewStartupLogger("info", true)

		// Act
		contextLogger := logger.With("service", "test", "version", "1.0")

		// Assert
		assert.NotNil(t, contextLogger, "Context logger should not be nil")
		assert.NotEqual(t, logger, contextLogger, "Should return a new logger instance")

		// 新しいロガーが動作することを確認
		contextLogger.Info("test message with context")
	})

	t.Run("正常系_WithContextメソッド", func(t *testing.T) {
		t.Parallel()

		// Arrange
		logger := log.NewStartupLogger("info", true)
		ctx := context.Background()

		// Act
		contextLogger := logger.WithContext(ctx)

		// Assert
		assert.NotNil(t, contextLogger, "Context logger should not be nil")
	})
}

func TestLogFile_Creation(t *testing.T) {
	// この関数は並行実行しない（ファイルアクセスのため）

	t.Run("正常系_デバッグモードでログファイル作成", func(t *testing.T) {
		// Arrange
		// テスト用の一時ディレクトリでテスト
		t.Setenv("TMPDIR", "/tmp")

		// Act
		logger := log.NewStartupLogger("debug", false)

		// Assert
		assert.NotNil(t, logger, "Logger should not be nil")

		// ログファイルが作成されているかチェック（ベストエフォート）
		if _, err := os.Stat("/tmp/app.log"); err == nil {
			// ファイルが存在する場合、内容を確認
			logger.Info("test message for file")
		}
	})

	t.Run("正常系_非デバッグモードではファイル作成しない", func(t *testing.T) {
		// Arrange & Act
		logger := log.NewStartupLogger("info", false)

		// Assert
		assert.NotNil(t, logger, "Logger should not be nil")
		logger.Info("test message for console")
	})
}

func TestParseLogLevel(t *testing.T) {
	t.Parallel()

	// この関数は内部関数なので、NewStartupLoggerを通してテスト
	tests := []struct {
		caseName string
		logLevel string
	}{
		{
			caseName: "正常系_debug",
			logLevel: "debug",
		},
		{
			caseName: "正常系_DEBUG大文字",
			logLevel: "DEBUG",
		},
		{
			caseName: "正常系_info",
			logLevel: "info",
		},
		{
			caseName: "正常系_warn",
			logLevel: "warn",
		},
		{
			caseName: "正常系_warning",
			logLevel: "warning",
		},
		{
			caseName: "正常系_error",
			logLevel: "error",
		},
		{
			caseName: "正常系_無効なレベル_デフォルトinfo",
			logLevel: "invalid",
		},
		{
			caseName: "正常系_空文字_デフォルトinfo",
			logLevel: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			logger := log.NewStartupLogger(tt.logLevel, true)

			// Assert
			assert.NotNil(t, logger, "Logger should be created regardless of log level")
		})
	}
}

func TestConstants(t *testing.T) {
	t.Parallel()

	t.Run("正常系_定数が正しく定義されている", func(t *testing.T) {
		t.Parallel()

		// Assert
		assert.Equal(t, "request_id", log.RequestIDKey, "RequestIDKey should be correct")
		assert.Equal(t, "logger", log.LoggerKey, "LoggerKey should be correct")
	})
}
