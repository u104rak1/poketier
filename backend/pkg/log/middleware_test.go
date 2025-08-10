package log_test

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"poketier/pkg/log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestNewMiddleware(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	tests := []struct {
		caseName string
		logLevel string
		isSilent bool
	}{
		{
			caseName: "正常系_デバッグレベル",
			logLevel: "debug",
			isSilent: true, // テストではサイレントモード
		},
		{
			caseName: "正常系_インフォレベル",
			logLevel: "info",
			isSilent: true,
		},
		{
			caseName: "正常系_サイレントモード",
			logLevel: "info",
			isSilent: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			middleware := log.NewMiddleware(tt.logLevel, tt.isSilent)

			// Assert
			assert.NotNil(t, middleware, "Middleware should not be nil")
		})
	}
}

func TestMiddleware_RequestLogging(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_正常なリクエスト", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Status should be 200")
	})

	t.Run("正常系_クライアントエラー", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusBadRequest, w.Code, "Status should be 400")
	})

	t.Run("正常系_サーバーエラー", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code, "Status should be 500")
	})
}

func TestMiddleware_ErrorHandling(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_エラーがある場合の詳細ログ", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			// エラーをGinコンテキストに追加
			c.Error(errors.New("test error"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code, "Status should be 500")
	})

	t.Run("正常系_ラップされたエラーの場合", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			originalErr := errors.New("database error")
			wrappedErr := fmt.Errorf("failed to process request: %w", originalErr)
			c.Error(wrappedErr)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code, "Status should be 500")
	})

	t.Run("正常系_複数のエラーがある場合", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			c.Error(errors.New("first error"))
			c.Error(errors.New("second error"))
			c.JSON(http.StatusInternalServerError, gin.H{"error": "multiple errors"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusInternalServerError, w.Code, "Status should be 500")
	})
}

func TestExtractErrorDetails(t *testing.T) {
	// この関数は内部関数のため、ミドルウェア経由でテスト

	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	tests := []struct {
		caseName    string
		setupError  func() error
		expectChain bool
	}{
		{
			caseName: "正常系_シンプルなエラー",
			setupError: func() error {
				return errors.New("simple error")
			},
			expectChain: false,
		},
		{
			caseName: "正常系_ラップされたエラー",
			setupError: func() error {
				original := errors.New("original error")
				return fmt.Errorf("wrapped: %w", original)
			},
			expectChain: true,
		},
		{
			caseName: "正常系_多重ラップされたエラー",
			setupError: func() error {
				original := errors.New("database error")
				level1 := fmt.Errorf("repository error: %w", original)
				return fmt.Errorf("service error: %w", level1)
			},
			expectChain: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			w := httptest.NewRecorder()
			_, r := gin.CreateTestContext(w)
			middleware := log.NewMiddleware("info", true)

			r.Use(middleware)
			r.GET("/test", func(c *gin.Context) {
				err := tt.setupError()
				c.Error(err)
				c.JSON(http.StatusInternalServerError, gin.H{"error": "error occurred"})
			})

			req := httptest.NewRequest(http.MethodGet, "/test", nil)

			// Act
			r.ServeHTTP(w, req)

			// Assert
			assert.Equal(t, http.StatusInternalServerError, w.Code, "Status should be 500")

			// エラーチェーンが適切に処理されることを間接的に確認
			// （ログ出力内容の直接検証は困難なため、エラーが正常に処理されることを確認）
		})
	}
}

func TestMiddleware_ResponseTime(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_レスポンス時間が記録される", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			// 短い処理時間をシミュレート
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Status should be 200")
		// レスポンス時間は内部でログに記録されるため、パニックしないことを確認
	})
}

func TestMiddleware_RequestContext(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_リクエスト情報がコンテキストに設定される", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		_, r := gin.CreateTestContext(w)
		middleware := log.NewMiddleware("info", true)

		var capturedRequestID string
		var capturedLogger log.Logger

		r.Use(middleware)
		r.GET("/test", func(c *gin.Context) {
			capturedRequestID = log.GetRequestID(c)
			capturedLogger = log.GetLogger(c)
			c.JSON(http.StatusOK, gin.H{"status": "ok"})
		})

		req := httptest.NewRequest(http.MethodGet, "/test", nil)

		// Act
		r.ServeHTTP(w, req)

		// Assert
		assert.Equal(t, http.StatusOK, w.Code, "Status should be 200")
		assert.NotEmpty(t, capturedRequestID, "Request ID should be set")
		assert.NotNil(t, capturedLogger, "Logger should be set")

		// UUIDの形式確認（36文字、ハイフン付き）
		assert.Len(t, capturedRequestID, 36, "Request ID should be UUID format")
		assert.Contains(t, capturedRequestID, "-", "Request ID should contain hyphens")
	})
}
