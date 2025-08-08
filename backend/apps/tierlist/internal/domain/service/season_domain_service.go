package service

import (
	"errors"
	"poketier/apps/tierlist/internal/domain/entity"
)

// SeasonDomainService はSeasonに関するドメインサービス
type SeasonDomainService struct{}

// NewSeasonDomainService は新しいSeasonDomainServiceインスタンスを作成する
func NewSeasonDomainService() *SeasonDomainService {
	return &SeasonDomainService{}
}

// EnsureUniqueActiveSeason はアクティブなシーズンが同時に1つのみであることを保証する
func (s *SeasonDomainService) EnsureUniqueActiveSeason(
	existingSeasons []*entity.Season,
	newSeason *entity.Season,
) error {
	// newSeasonのnilチェック
	if newSeason == nil {
		return errors.New("new season cannot be nil")
	}

	// 新しいシーズンがアクティブでない場合はチェック不要
	if !newSeason.IsActive() {
		return nil
	}

	// 既存のシーズンにアクティブなものがあるかチェック
	for _, season := range existingSeasons {
		if season != nil && season.IsActive() {
			return errors.New("active season already exists")
		}
	}

	return nil
}
