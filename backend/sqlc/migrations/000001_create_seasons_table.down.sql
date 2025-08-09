-- トリガーとファンクションを削除
DROP TRIGGER IF EXISTS update_seasons_updated_at ON seasons;
DROP FUNCTION IF EXISTS update_updated_at_column();

-- インデックスを削除
DROP INDEX IF EXISTS idx_seasons_active;

-- テーブルを削除
DROP TABLE IF EXISTS seasons;
