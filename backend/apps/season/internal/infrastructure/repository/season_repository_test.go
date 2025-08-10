package repository_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"poketier/apps/season/internal/domain/entity"
	"poketier/apps/season/internal/infrastructure/repository"
	"poketier/pkg/vo/id"
	"poketier/sqlc/db"
)

var (
	seasonID  = id.NewSeasonID()
	seasonID2 = id.NewSeasonID()
)

func TestSeasonRepository_FindByID(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		seasonID    id.SeasonID
		want        *entity.Season
		expectError bool
	}{
		{
			caseName: "正常系: 指定されたIDのSeasonが取得できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				seasonUUID := pgtype.UUID{
					Bytes: seasonID.UUID(),
					Valid: true,
				}
				dbSeason := db.Season{
					SeasonID: seasonUUID,
					Name:     "A2b",
					StartDate: pgtype.Date{
						Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
				mockQuerier.EXPECT().GetSeason(gomock.Any(), seasonUUID).Return(dbSeason, nil)
			},
			seasonID: seasonID,
			want: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"A2b",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				)
				return season
			}(),
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				seasonUUID := pgtype.UUID{
					Bytes: seasonID.UUID(),
					Valid: true,
				}
				mockQuerier.EXPECT().GetSeason(gomock.Any(), seasonUUID).Return(db.Season{}, errors.New("db error"))
			},
			seasonID:    seasonID,
			want:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			got, err := repo.FindByID(context.Background(), tt.seasonID)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.Equal(t, tt.want.ID(), got.ID(), "season ID does not match")
			assert.Equal(t, tt.want.Name(), got.Name(), "season name does not match")
			assert.Equal(t, tt.want.StartDate(), got.StartDate(), "start date does not match")
			assert.Equal(t, tt.want.EndDate(), got.EndDate(), "end date does not match")
			assert.Equal(t, tt.want.IsActive(), got.IsActive(), "is active does not match")
		})
	}
}

func TestSeasonRepository_FindActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		want        *entity.Season
		expectError bool
	}{
		{
			caseName: "正常系: アクティブなSeasonが取得できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				dbSeason := db.Season{
					SeasonID: pgtype.UUID{
						Bytes: seasonID.UUID(),
						Valid: true,
					},
					Name: "S1",
					StartDate: pgtype.Date{
						Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
						Valid: true,
					},
				}
				mockQuerier.EXPECT().GetActiveSeason(gomock.Any()).Return(dbSeason, nil)
			},
			want: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				)
				return season
			}(),
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				mockQuerier.EXPECT().GetActiveSeason(gomock.Any()).Return(db.Season{}, errors.New("db error"))
			},
			want:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			got, err := repo.FindActive(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.Equal(t, tt.want.ID(), got.ID(), "season ID does not match")
			assert.Equal(t, tt.want.Name(), got.Name(), "season name does not match")
			assert.Equal(t, tt.want.StartDate(), got.StartDate(), "start date does not match")
			assert.Equal(t, tt.want.EndDate(), got.EndDate(), "end date does not match")
			assert.Equal(t, tt.want.IsActive(), got.IsActive(), "is active does not match")
		})
	}
}

func TestSeasonRepository_FindAll(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		want        []*entity.Season
		expectError bool
	}{
		{
			caseName: "正常系: 全てのSeasonが取得できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				dbSeasons := []db.Season{
					{
						SeasonID: pgtype.UUID{
							Bytes: seasonID.UUID(),
							Valid: true,
						},
						Name: "S1",
						StartDate: pgtype.Date{
							Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							Valid: true,
						},
						EndDate: pgtype.Date{
							Time:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
							Valid: true,
						},
					},
					{
						SeasonID: pgtype.UUID{
							Bytes: seasonID2.UUID(),
							Valid: true,
						},
						Name: "S2",
						StartDate: pgtype.Date{
							Time:  time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
							Valid: true,
						},
						EndDate: pgtype.Date{
							Time:  time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
							Valid: true,
						},
					},
				}
				mockQuerier.EXPECT().ListSeasons(gomock.Any()).Return(dbSeasons, nil)
			},
			want: func() []*entity.Season {
				season1, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				)
				season2, _ := entity.NewSeason(
					seasonID2,
					"S2",
					time.Date(2023, 4, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 6, 30, 0, 0, 0, 0, time.UTC),
				)
				return []*entity.Season{season1, season2}
			}(),
			expectError: false,
		},
		{
			caseName: "正常系: 空のリストが返される事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				mockQuerier.EXPECT().ListSeasons(gomock.Any()).Return([]db.Season{}, nil)
			},
			want:        []*entity.Season{},
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				mockQuerier.EXPECT().ListSeasons(gomock.Any()).Return(nil, errors.New("db error"))
			},
			want:        nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			got, err := repo.FindAll(context.Background())

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.Len(t, got, len(tt.want), "seasons length does not match")
			for i, expectedSeason := range tt.want {
				assert.Equal(t, expectedSeason.ID(), got[i].ID(), "season ID at index %d does not match", i)
				assert.Equal(t, expectedSeason.Name(), got[i].Name(), "season name at index %d does not match", i)
				assert.Equal(t, expectedSeason.StartDate(), got[i].StartDate(), "start date at index %d does not match", i)
				assert.Equal(t, expectedSeason.EndDate(), got[i].EndDate(), "end date at index %d does not match", i)
				assert.Equal(t, expectedSeason.IsActive(), got[i].IsActive(), "is active at index %d does not match", i)
			}
		})
	}
}

func TestSeasonRepository_Save(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		season      *entity.Season
		expectError bool
	}{
		{
			caseName: "正常系: Seasonが保存できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				expectedParams := db.SaveSeasonParams{
					SeasonID: pgtype.UUID{
						Bytes: seasonID.UUID(),
						Valid: true,
					},
					Name: "S1",
					StartDate: pgtype.Date{
						Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
						Valid: true,
					},
				}
				mockQuerier.EXPECT().SaveSeason(gomock.Any(), expectedParams).Return(db.Season{}, nil)
			},
			season: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				)
				return season
			}(),
			expectError: false,
		},
		{
			caseName: "正常系: Seasonが保存できる事（別の終了日）",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				expectedParams := db.SaveSeasonParams{
					SeasonID: pgtype.UUID{
						Bytes: seasonID.UUID(),
						Valid: true,
					},
					Name: "S1",
					StartDate: pgtype.Date{
						Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
				mockQuerier.EXPECT().SaveSeason(gomock.Any(), expectedParams).Return(db.Season{}, nil)
			},
			season: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				)
				return season
			}(),
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				mockQuerier.EXPECT().SaveSeason(gomock.Any(), gomock.Any()).Return(db.Season{}, errors.New("db error"))
			},
			season: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				)
				return season
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			err := repo.Save(context.Background(), tt.season)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
		})
	}
}

func TestSeasonRepository_Update(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		season      *entity.Season
		expectError bool
	}{
		{
			caseName: "正常系: Seasonが更新できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				expectedParams := db.UpdateSeasonParams{
					SeasonID: pgtype.UUID{
						Bytes: seasonID.UUID(),
						Valid: true,
					},
					Name: "S1",
					StartDate: pgtype.Date{
						Time:  time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
					EndDate: pgtype.Date{
						Time:  time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
						Valid: true,
					},
				}
				mockQuerier.EXPECT().UpdateSeason(gomock.Any(), expectedParams).Return(db.Season{}, nil)
			},
			season: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 3, 31, 0, 0, 0, 0, time.UTC),
				)
				return season
			}(),
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				mockQuerier.EXPECT().UpdateSeason(gomock.Any(), gomock.Any()).Return(db.Season{}, errors.New("db error"))
			},
			season: func() *entity.Season {
				season, _ := entity.NewSeason(
					seasonID,
					"S1",
					time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					time.Date(2023, 12, 31, 23, 59, 59, 0, time.UTC),
				)
				return season
			}(),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			err := repo.Update(context.Background(), tt.season)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
		})
	}
}

func TestSeasonRepository_Delete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		setupMock   func(mockQuerier *MockSeasonQuerier)
		seasonID    id.SeasonID
		expectError bool
	}{
		{
			caseName: "正常系: Seasonが削除できる事",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				seasonUUID := pgtype.UUID{
					Bytes: seasonID.UUID(),
					Valid: true,
				}
				mockQuerier.EXPECT().DeleteSeason(gomock.Any(), seasonUUID).Return(nil)
			},
			seasonID:    seasonID,
			expectError: false,
		},
		{
			caseName: "異常系: DBエラーが発生した場合",
			setupMock: func(mockQuerier *MockSeasonQuerier) {
				seasonUUID := pgtype.UUID{
					Bytes: seasonID.UUID(),
					Valid: true,
				}
				mockQuerier.EXPECT().DeleteSeason(gomock.Any(), seasonUUID).Return(errors.New("db error"))
			},
			seasonID:    seasonID,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockQuerier := NewMockSeasonQuerier(ctrl)
			tt.setupMock(mockQuerier)
			repo := repository.NewSeasonRepository(mockQuerier)

			// Act
			err := repo.Delete(context.Background(), tt.seasonID)

			// Assert
			if tt.expectError {
				assert.Error(t, err, "expected error but got none")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
		})
	}
}
