package env_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"poketier/env"
)

func TestNewEnv(t *testing.T) {
	tests := []struct {
		caseName string
		envVars  map[string]string
		want     *env.Env
	}{
		{
			caseName: "正常系: 環境変数が設定されていない場合デフォルト値が使用される",
			envVars:  map[string]string{},
			want: &env.Env{
				APP_PORT:          "8080",
				POSTGRES_HOST:     "postgres",
				POSTGRES_DBNAME:   "poketierlocal",
				POSTGRES_USER:     "dbuser",
				POSTGRES_PASSWORD: "Password123",
				POSTGRES_PORT:     "5432",
				POSTGRES_SSLMODE:  "disable",
				LOG_LEVEL:         "debug",
				IS_SILENT_LOG:     false,
			},
		},
		{
			caseName: "正常系: 環境変数で設定した値が正しく取得される",
			envVars: map[string]string{
				"APP_PORT":          "9000",
				"POSTGRES_HOST":     "localhost",
				"POSTGRES_DBNAME":   "test_db",
				"POSTGRES_USER":     "test_user",
				"POSTGRES_PASSWORD": "test_password",
				"POSTGRES_PORT":     "5433",
				"POSTGRES_SSLMODE":  "require",
				"LOG_LEVEL":         "info",
				"IS_SILENT_LOG":     "true",
			},
			want: &env.Env{
				APP_PORT:          "9000",
				POSTGRES_HOST:     "localhost",
				POSTGRES_DBNAME:   "test_db",
				POSTGRES_USER:     "test_user",
				POSTGRES_PASSWORD: "test_password",
				POSTGRES_PORT:     "5433",
				POSTGRES_SSLMODE:  "require",
				LOG_LEVEL:         "info",
				IS_SILENT_LOG:     true,
			},
		},
		{
			caseName: "正常系: 一部の環境変数のみ設定した場合、設定した値とデフォルト値が混在する",
			envVars: map[string]string{
				"APP_PORT":        "3000",
				"POSTGRES_DBNAME": "custom_db",
			},
			want: &env.Env{
				APP_PORT:          "3000",
				POSTGRES_HOST:     "postgres",
				POSTGRES_DBNAME:   "custom_db",
				POSTGRES_USER:     "dbuser",
				POSTGRES_PASSWORD: "Password123",
				POSTGRES_PORT:     "5432",
				POSTGRES_SSLMODE:  "disable",
				LOG_LEVEL:         "debug",
				IS_SILENT_LOG:     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			// Arrange
			for key, value := range tt.envVars {
				t.Setenv(key, value)
			}

			// Act
			got := env.NewEnv()

			// Assert
			assert.Equal(t, tt.want, got, "environment variable values do not match expected values")
		})
	}
}

func TestNewEnv_DefaultValues(t *testing.T) {
	t.Run("正常系: 環境変数が未設定の場合、全てのフィールドでデフォルト値が使用される", func(t *testing.T) {
		// Act
		got := env.NewEnv()

		// Assert
		assert.Equal(t, "8080", got.APP_PORT, "APP_PORT default value is incorrect")
		assert.Equal(t, "postgres", got.POSTGRES_HOST, "POSTGRES_HOST default value is incorrect")
		assert.Equal(t, "poketierlocal", got.POSTGRES_DBNAME, "POSTGRES_DBNAME default value is incorrect")
		assert.Equal(t, "dbuser", got.POSTGRES_USER, "POSTGRES_USER default value is incorrect")
		assert.Equal(t, "Password123", got.POSTGRES_PASSWORD, "POSTGRES_PASSWORD default value is incorrect")
		assert.Equal(t, "5432", got.POSTGRES_PORT, "POSTGRES_PORT default value is incorrect")
		assert.Equal(t, "disable", got.POSTGRES_SSLMODE, "POSTGRES_SSLMODE default value is incorrect")
		assert.Equal(t, "debug", got.LOG_LEVEL, "LOG_LEVEL default value is incorrect")
		assert.Equal(t, false, got.IS_SILENT_LOG, "IS_SILENT_LOG default value is incorrect")
	})
}

func TestNewEnv_EnvironmentVariables(t *testing.T) {
	t.Run("正常系: 全ての環境変数を設定した場合、設定値が正しく読み込まれる", func(t *testing.T) {
		// Arrange
		t.Setenv("APP_PORT", "9999")
		t.Setenv("POSTGRES_HOST", "test-host")
		t.Setenv("POSTGRES_DBNAME", "test-database")
		t.Setenv("POSTGRES_USER", "test-username")
		t.Setenv("POSTGRES_PASSWORD", "test-secret")
		t.Setenv("POSTGRES_PORT", "5555")
		t.Setenv("POSTGRES_SSLMODE", "verify-full")
		t.Setenv("LOG_LEVEL", "error")
		t.Setenv("IS_SILENT_LOG", "true")

		// Act
		got := env.NewEnv()

		// Assert
		assert.Equal(t, "9999", got.APP_PORT, "APP_PORT environment variable is not set correctly")
		assert.Equal(t, "test-host", got.POSTGRES_HOST, "POSTGRES_HOST environment variable is not set correctly")
		assert.Equal(t, "test-database", got.POSTGRES_DBNAME, "POSTGRES_DBNAME environment variable is not set correctly")
		assert.Equal(t, "test-username", got.POSTGRES_USER, "POSTGRES_USER environment variable is not set correctly")
		assert.Equal(t, "test-secret", got.POSTGRES_PASSWORD, "POSTGRES_PASSWORD environment variable is not set correctly")
		assert.Equal(t, "5555", got.POSTGRES_PORT, "POSTGRES_PORT environment variable is not set correctly")
		assert.Equal(t, "verify-full", got.POSTGRES_SSLMODE, "POSTGRES_SSLMODE environment variable is not set correctly")
		assert.Equal(t, "error", got.LOG_LEVEL, "LOG_LEVEL environment variable is not set correctly")
		assert.Equal(t, true, got.IS_SILENT_LOG, "IS_SILENT_LOG environment variable is not set correctly")
	})
}
