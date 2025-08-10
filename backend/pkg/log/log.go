package log

import (
	"context"
	"io"
	"log/slog"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Logger インターフェースでテスト時のモック化を容易に
type Logger interface {
	Debug(msg string, args ...any)
	Info(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	With(args ...any) Logger
	WithContext(ctx context.Context) Logger
}

// slogLogger slogのラッパー実装
type slogLogger struct {
	logger *slog.Logger
}

const (
	RequestIDKey = "request_id"
	LoggerKey    = "logger"
)

// GetLogger コンテキストからロガーを取得
func GetLogger(c *gin.Context) Logger {
	if logger, exists := c.Get(LoggerKey); exists {
		return logger.(Logger)
	}
	// フォールバック用のデフォルトロガー
	return newLogger("info", false)
}

// GetRequestID コンテキストからリクエストIDを取得
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		return requestID.(string)
	}
	return ""
}

// NewStartupLogger アプリケーション起動時用のロガーを作成
func NewStartupLogger(logLevel string, isSilent bool) Logger {
	return newLogger(logLevel, isSilent)
}

// newLogger 環境変数から設定を取得してロガーを作成（パッケージ内限定）
func newLogger(logLevel string, isSilent bool) Logger {
	var handler slog.Handler

	if isSilent {
		// サイレントモード：何も出力しない
		handler = slog.NewJSONHandler(io.Discard, &slog.HandlerOptions{
			Level: slog.LevelError + 1,
		})
	} else {
		level := parseLogLevel(logLevel)

		// 開発環境（debugレベル）の場合は、tmpfsにログファイル出力
		if level == slog.LevelDebug {
			// tmpfs（メモリ上）にログファイルを作成、毎回上書き
			// コンテナ再起動時に自動削除される
			logFile, err := os.OpenFile("/tmp/app.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
			if err != nil {
				// ファイル作成に失敗した場合はコンソールのみ
				handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
					Level: level,
				})
			} else {
				// ファイルのみに出力（デバッガー経由でもログを確認可能）
				handler = slog.NewJSONHandler(logFile, &slog.HandlerOptions{
					Level: level,
				})
			}
		} else {
			// 本番環境：JSON形式でコンソール出力
			handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: level,
			})
		}
	}

	return &slogLogger{
		logger: slog.New(handler),
	}
}

// parseLogLevel 文字列からログレベルを解析
func parseLogLevel(level string) slog.Level {
	switch strings.ToLower(level) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func (l *slogLogger) Debug(msg string, args ...any) {
	l.logger.Debug(msg, args...)
}

func (l *slogLogger) Info(msg string, args ...any) {
	l.logger.Info(msg, args...)
}

func (l *slogLogger) Warn(msg string, args ...any) {
	l.logger.Warn(msg, args...)
}

func (l *slogLogger) Error(msg string, args ...any) {
	l.logger.Error(msg, args...)
}

func (l *slogLogger) With(args ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(args...),
	}
}

func (l *slogLogger) WithContext(ctx context.Context) Logger {
	// Ginのコンテキストからリクエスト情報を取得
	return l
}
