package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"poketier/apps/season/internal/application/usecase"
	"poketier/apps/season/internal/domain/entity"
	"poketier/pkg/vo/id"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListSeasonUsecase_Execute(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(*MockLSSeasonRepository)
		wantResult  *usecase.ListSeasonsResult
		wantErr     bool
		errContains string
	}{
		{
			caseName: "正常系: シーズンが存在する場合、シーズン一覧を返す",
			setupMock: func(mockRepo *MockLSSeasonRepository) {
				seasons := []*entity.Season{
					createTestSeason(t, "A1a", time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)),
					createTestSeason(t, "A1b", time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)),
				}
				mockRepo.EXPECT().FindAll(gomock.Any()).Return(seasons, nil)
			},
			wantResult: &usecase.ListSeasonsResult{
				Seasons: []usecase.LSSeason{
					{
						SeasonID:  "550e8400-e29b-41d4-a716-446655440000",
						Name:      "A1a",
						StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
						IsActive:  true,
					},
					{
						SeasonID:  "550e8400-e29b-41d4-a716-446655440000",
						Name:      "A1b",
						StartDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC),
						IsActive:  false,
					},
				},
			},
			wantErr: false,
		},
		{
			caseName: "正常系: シーズンが存在しない場合、空のシーズン一覧を返す",
			setupMock: func(mockRepo *MockLSSeasonRepository) {
				mockRepo.EXPECT().FindAll(gomock.Any()).Return([]*entity.Season{}, nil)
			},
			wantResult: &usecase.ListSeasonsResult{
				Seasons: []usecase.LSSeason{},
			},
			wantErr: false,
		},
		{
			caseName: "異常系: リポジトリでエラーが発生した場合、エラーを返す",
			setupMock: func(mockRepo *MockLSSeasonRepository) {
				mockRepo.EXPECT().FindAll(gomock.Any()).Return(nil, errors.New("repository error"))
			},
			wantResult:  nil,
			wantErr:     true,
			errContains: "repository error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockRepo := NewMockLSSeasonRepository(ctrl)
			tt.setupMock(mockRepo)

			usecase := usecase.NewListSeasonsUsecase(mockRepo)
			ctx := context.Background()

			// Act
			got, err := usecase.Execute(ctx)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains, "error message does not contain expected text")
				}
				return
			}

			assert.NoError(t, err, "unexpected error occurred")
			assert.Equal(t, tt.wantResult, got, "result does not match expected value")
		})
	}
}

// createTestSeason はテスト用のSeasonエンティティを作成するヘルパー関数
func createTestSeason(t *testing.T, name string, endDate time.Time) *entity.Season {
	t.Helper()

	seasonID, err := id.SeasonIDFromString("550e8400-e29b-41d4-a716-446655440000")
	assert.NoError(t, err, "failed to create season ID")

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	season, err := entity.NewSeason(seasonID, name, startDate, endDate)
	assert.NoError(t, err, "failed to create season entity")

	return season
}
