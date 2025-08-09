package sqlc

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"poketier/config/env"
)

// NewPgxPool は 環境変数からPostgreSQLへの接続プールを作成
func NewPgxPool(ctx context.Context, env *env.Env) (*pgxpool.Pool, error) {
	// PostgreSQL接続文字列を構築
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		env.POSTGRES_HOST,
		env.POSTGRES_PORT,
		env.POSTGRES_USER,
		env.POSTGRES_PASSWORD,
		env.POSTGRES_DBNAME,
		env.POSTGRES_SSLMODE,
	)

	// 接続プール設定
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to parse database config: %w", err)
	}

	// プール設定の調整
	config.MaxConns = 25                      // 最大接続数
	config.MinConns = 5                       // 最小接続数
	config.MaxConnLifetime = time.Hour        // 接続の最大生存時間
	config.MaxConnIdleTime = time.Minute * 30 // アイドル接続の最大時間

	// 接続プールを作成
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// 接続テスト
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return pool, nil
}

// ClosePgxPool は 接続プールを安全に閉じる
func ClosePgxPool(pool *pgxpool.Pool) {
	if pool != nil {
		pool.Close()
	}
}
