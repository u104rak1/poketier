-- シーズンのCRUD操作

-- name: SaveSeason :one
-- Upsert: 存在する場合は更新、しない場合は挿入
INSERT INTO seasons (
    season_id,
    name,
    start_date,
    end_date,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
) ON CONFLICT (season_id) 
DO UPDATE SET
    name = EXCLUDED.name,
    start_date = EXCLUDED.start_date,
    end_date = EXCLUDED.end_date,
    is_active = EXCLUDED.is_active
RETURNING *;

-- name: CreateSeason :one
INSERT INTO seasons (
    season_id,
    name,
    start_date,
    end_date,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetSeason :one
SELECT * FROM seasons
WHERE season_id = $1;

-- name: GetActiveSeason :one
SELECT * FROM seasons
WHERE is_active = true
LIMIT 1;

-- name: ListSeasons :many
SELECT * FROM seasons
ORDER BY start_date DESC;

-- name: UpdateSeason :one
UPDATE seasons
SET 
    name = $2,
    start_date = $3,
    end_date = $4,
    is_active = $5
WHERE season_id = $1
RETURNING *;

-- name: DeleteSeason :exec
DELETE FROM seasons
WHERE season_id = $1;

-- name: SetActiveSeason :exec
-- 既存のアクティブシーズンを無効化してから新しいシーズンを有効化
UPDATE seasons 
SET is_active = CASE 
    WHEN season_id = $1 THEN true 
    ELSE false 
END;

-- name: CountSeasons :one
SELECT COUNT(*) FROM seasons;

-- name: BulkCreateSeasons :copyfrom
INSERT INTO seasons (
    season_id,
    name,
    start_date,
    end_date,
    is_active
) VALUES (
    $1, $2, $3, $4, $5
);

-- name: DeleteAllSeasons :exec
-- 開発・テスト用: 全シーズンを削除
DELETE FROM seasons;

-- name: BulkDeleteSeasons :exec
-- 指定したIDリストのシーズンを一括削除
DELETE FROM seasons
WHERE season_id = ANY($1::uuid[]);
