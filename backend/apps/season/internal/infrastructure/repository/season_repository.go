package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"poketier/apps/season/internal/domain/entity"
	"poketier/pkg/vo/id"
	"poketier/sqlc/db"
)

// SeasonQuerier はデータベースクエリを定義するインターフェース
type SeasonQuerier interface {
	GetSeason(ctx context.Context, seasonID pgtype.UUID) (db.Season, error)
	GetActiveSeason(ctx context.Context) (db.Season, error)
	ListSeasons(ctx context.Context) ([]db.Season, error)
	SaveSeason(ctx context.Context, arg db.SaveSeasonParams) (db.Season, error)
	UpdateSeason(ctx context.Context, arg db.UpdateSeasonParams) (db.Season, error)
	DeleteSeason(ctx context.Context, seasonID pgtype.UUID) error
}

// seasonRepository はSeasonRepositoryの実装
type seasonRepository struct {
	queries SeasonQuerier
}

// NewSeasonRepository は新しいSeasonRepositoryを作成
func NewSeasonRepository(queries SeasonQuerier) *seasonRepository {
	return &seasonRepository{
		queries: queries,
	}
}

// FindByID は指定されたIDのSeasonを取得
func (r *seasonRepository) FindByID(ctx context.Context, seasonID id.SeasonID) (*entity.Season, error) {
	seasonUUID := pgtype.UUID{
		Bytes: seasonID.UUID(),
		Valid: true,
	}

	dbSeason, err := r.queries.GetSeason(ctx, seasonUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get season by ID: %w", err)
	}

	return r.toEntity(dbSeason)
}

// FindActive はアクティブなSeasonを取得
func (r *seasonRepository) FindActive(ctx context.Context) (*entity.Season, error) {
	activeSeason, err := r.queries.GetActiveSeason(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active season: %w", err)
	}

	return r.toEntity(activeSeason)
}

// FindAll は全てのSeasonを取得
func (r *seasonRepository) FindAll(ctx context.Context) ([]*entity.Season, error) {
	dbSeasons, err := r.queries.ListSeasons(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list seasons: %w", err)
	}

	seasons := make([]*entity.Season, 0, len(dbSeasons))
	for _, dbSeason := range dbSeasons {
		season, err := r.toEntity(dbSeason)
		if err != nil {
			return nil, err
		}
		seasons = append(seasons, season)
	}

	return seasons, nil
}

// Save は新しいSeasonを保存
func (r *seasonRepository) Save(ctx context.Context, season *entity.Season) error {
	params := r.toSaveParams(season)

	_, err := r.queries.SaveSeason(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to save season: %w", err)
	}

	return nil
}

// Update は既存のSeasonを更新
func (r *seasonRepository) Update(ctx context.Context, season *entity.Season) error {
	params := r.toUpdateParams(season)

	_, err := r.queries.UpdateSeason(ctx, params)
	if err != nil {
		return fmt.Errorf("failed to update season: %w", err)
	}

	return nil
}

// Delete は指定されたIDのSeasonを削除
func (r *seasonRepository) Delete(ctx context.Context, seasonID id.SeasonID) error {
	seasonUUID := pgtype.UUID{
		Bytes: seasonID.UUID(),
		Valid: true,
	}

	err := r.queries.DeleteSeason(ctx, seasonUUID)
	if err != nil {
		return fmt.Errorf("failed to delete season: %w", err)
	}

	return nil
}

// toEntity はデータベースモデルからエンティティに変換
func (r *seasonRepository) toEntity(dbSeason db.Season) (*entity.Season, error) {
	// SeasonIDを変換
	seasonID := id.SeasonIDFromUUID(dbSeason.SeasonID.Bytes)

	// 開始日を変換
	startDate := dbSeason.StartDate.Time

	// 終了日を変換（nilの場合もある）
	var endDate *time.Time
	if dbSeason.EndDate.Valid {
		endDate = &dbSeason.EndDate.Time
	}

	// エンティティを作成
	season, err := entity.NewSeason(seasonID, dbSeason.Name, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to create season entity: %w", err)
	}

	return season, nil
}

// toSaveParams はエンティティからSave用パラメータに変換
func (r *seasonRepository) toSaveParams(season *entity.Season) db.SaveSeasonParams {
	// 終了日の変換
	var endDate pgtype.Date
	if season.EndDate() != nil {
		endDate = pgtype.Date{
			Time:  *season.EndDate(),
			Valid: true,
		}
	}

	return db.SaveSeasonParams{
		SeasonID: pgtype.UUID{
			Bytes: season.ID().UUID(),
			Valid: true,
		},
		Name: season.Name(),
		StartDate: pgtype.Date{
			Time:  season.StartDate(),
			Valid: true,
		},
		EndDate: endDate,
		IsActive: pgtype.Bool{
			Bool:  season.IsActive(),
			Valid: true,
		},
	}
}

// toUpdateParams はエンティティからUpdate用パラメータに変換
func (r *seasonRepository) toUpdateParams(season *entity.Season) db.UpdateSeasonParams {
	// 終了日の変換
	var endDate pgtype.Date
	if season.EndDate() != nil {
		endDate = pgtype.Date{
			Time:  *season.EndDate(),
			Valid: true,
		}
	}

	return db.UpdateSeasonParams{
		SeasonID: pgtype.UUID{
			Bytes: season.ID().UUID(),
			Valid: true,
		},
		Name: season.Name(),
		StartDate: pgtype.Date{
			Time:  season.StartDate(),
			Valid: true,
		},
		EndDate: endDate,
		IsActive: pgtype.Bool{
			Bool:  season.IsActive(),
			Valid: true,
		},
	}
}
