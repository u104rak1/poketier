package usecase

import (
	"context"
	"poketier/apps/tierlist/internal/domain/entity"
	"time"
)

// ListSeasonsResult はシーズン一覧取得結果
type ListSeasonsResult struct {
	Seasons []LSSeason
}

type LSSeason struct {
	SeasonID  string
	Name      string
	StartDate time.Time
	EndDate   *time.Time
	IsActive  bool
}

type LSSeasonRepository interface {
	FindAll(ctx context.Context) ([]*entity.Season, error)
}

type ListSeasonsUsecase struct {
	seasonRepo LSSeasonRepository
}

func NewListSeasonsUsecase(seasonRepo LSSeasonRepository) *ListSeasonsUsecase {
	return &ListSeasonsUsecase{
		seasonRepo: seasonRepo,
	}
}

// Execute はジーズン一覧取得を実行
func (u *ListSeasonsUsecase) Execute(ctx context.Context) (*ListSeasonsResult, error) {
	seasons, err := u.seasonRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	return u.toResult(seasons), nil
}

func (u *ListSeasonsUsecase) toResult(seasons []*entity.Season) *ListSeasonsResult {
	lsSeasons := make([]LSSeason, 0, len(seasons))
	for _, season := range seasons {
		lsSeasons = append(lsSeasons, LSSeason{
			SeasonID:  season.ID().String(),
			Name:      season.Name(),
			StartDate: season.StartDate(),
			EndDate:   season.EndDate(),
			IsActive:  season.IsActive(),
		})
	}
	return &ListSeasonsResult{
		Seasons: lsSeasons,
	}
}
