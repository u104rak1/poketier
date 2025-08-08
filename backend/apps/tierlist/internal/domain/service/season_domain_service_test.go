package service_test

import (
"testing"
"time"

"poketier/apps/tierlist/internal/domain/entity"
"poketier/apps/tierlist/internal/domain/service"
"poketier/pkg/vo/id"

"github.com/stretchr/testify/assert"
)

func TestSeasonDomainService_EnsureUniqueActiveSeason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName        string
		existingSeasons []*entity.Season
		newSeason       *entity.Season
		wantErr         bool
	}{
		{
			caseName:        "正常系: アクティブなシーズンが存在しない場合は新しいアクティブシーズンを作成可能",
			existingSeasons: []*entity.Season{},
			newSeason:       createTestSeason(t, "A2b", nil),
			wantErr:         false,
		},
		{
			caseName: "正常系: 既存のシーズンが全て終了している場合は新しいアクティブシーズンを作成可能",
			existingSeasons: []*entity.Season{
				createTestSeason(t, "A2a", &[]time.Time{time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)}[0]),
			},
			newSeason: createTestSeason(t, "A2b", nil),
			wantErr:   false,
		},
		{
			caseName: "正常系: 新しいシーズンが既に終了している場合は作成可能",
			existingSeasons: []*entity.Season{
				createTestSeason(t, "A2a", nil),
			},
			newSeason: createTestSeason(t, "A2b", &[]time.Time{time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)}[0]),
			wantErr:   false,
		},
		{
			caseName: "異常系: 既にアクティブなシーズンが存在する場合は新しいアクティブシーズンを作成不可",
			existingSeasons: []*entity.Season{
				createTestSeason(t, "A2a", nil),
			},
			newSeason: createTestSeason(t, "A2b", nil),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
t.Parallel()

			// Arrange
			domainService := service.NewSeasonDomainService()

			// Act
			err := domainService.EnsureUniqueActiveSeason(tt.existingSeasons, tt.newSeason)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
			} else {
				assert.NoError(t, err, "unexpected error occurred")
			}
		})
	}
}

// createTestSeason はテスト用のSeasonを作成するヘルパー関数
func createTestSeason(t *testing.T, name string, endDate *time.Time) *entity.Season {
	t.Helper()
	seasonID := id.NewSeasonID()
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	// 終了日が設定されている場合は、開始日より後になるように調整
	if endDate != nil && endDate.Before(startDate) {
		adjustedEndDate := startDate.Add(24 * time.Hour)
		endDate = &adjustedEndDate
	}

	season, err := entity.NewSeason(seasonID, name, startDate, endDate)
	assert.NoError(t, err, "failed to create test season")
	return season
}
