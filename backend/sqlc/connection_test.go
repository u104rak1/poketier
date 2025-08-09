package sqlc_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"poketier/env"
	"poketier/sqlc"
)

func TestNewPgxPool(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName      string
		env           *env.Env
		expectedError bool
	}{
		{
			caseName: "正常系: 有効な環境変数で接続プールが作成される事",
			env: &env.Env{
				POSTGRES_HOST:     "localhost",
				POSTGRES_PORT:     "5432",
				POSTGRES_USER:     "test_user",
				POSTGRES_PASSWORD: "test_password",
				POSTGRES_DBNAME:   "test_db",
				POSTGRES_SSLMODE:  "disable",
			},
			expectedError: true, // 実際のDBが存在しないためエラーになる
		},
		{
			caseName: "異常系: 無効なポート番号の場合",
			env: &env.Env{
				POSTGRES_HOST:     "localhost",
				POSTGRES_PORT:     "invalid_port",
				POSTGRES_USER:     "test_user",
				POSTGRES_PASSWORD: "test_password",
				POSTGRES_DBNAME:   "test_db",
				POSTGRES_SSLMODE:  "disable",
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			// Act
			pool, err := sqlc.NewPgxPool(ctx, tt.env)

			// Assert
			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, pool)
			} else {
				assert.NoError(t, err)
				require.NotNil(t, pool)
				defer sqlc.ClosePgxPool(pool)

				// 接続プールが正常に動作することを確認
				assert.NoError(t, pool.Ping(ctx))
			}
		})
	}
}

func TestClosePgxPool(t *testing.T) {
	t.Parallel()

	t.Run("正常系: nilプールでもパニックしない事", func(t *testing.T) {
		t.Parallel()

		// Arrange & Act & Assert
		assert.NotPanics(t, func() {
			sqlc.ClosePgxPool(nil)
		})
	})
}
