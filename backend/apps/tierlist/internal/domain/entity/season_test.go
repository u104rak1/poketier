package entity_test

import (
	"testing"
	"time"

	"poketier/apps/tierlist/internal/domain/entity"
	"poketier/pkg/vo/id"

	"github.com/stretchr/testify/assert"
)

const testSeasonName = "A2b"

func TestNewSeason(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName  string
		id        id.SeasonID
		name      string
		startDate time.Time
		endDate   *time.Time
		wantErr   bool
	}{
		{
			caseName:  "正常系: 有効なパラメータでSeasonが作成される",
			id:        id.NewSeasonID(),
			name:      "A2b",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   nil,
			wantErr:   false,
		},
		{
			caseName:  "正常系: 終了日が設定されたSeasonが作成される",
			id:        id.NewSeasonID(),
			name:      "A2a",
			startDate: time.Date(2024, 12, 1, 0, 0, 0, 0, time.UTC),
			endDate:   &[]time.Time{time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)}[0],
			wantErr:   false,
		},
		{
			caseName:  "異常系: 空のシーズン名が渡された場合",
			id:        id.NewSeasonID(),
			name:      "",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   nil,
			wantErr:   true,
		},
		{
			caseName:  "異常系: 5文字を超えるシーズン名が渡された場合",
			id:        id.NewSeasonID(),
			name:      "Season1",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   nil,
			wantErr:   true,
		},
		{
			caseName:  "正常系: 5文字のシーズン名が渡された場合",
			id:        id.NewSeasonID(),
			name:      "Sea5n",
			startDate: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			endDate:   nil,
			wantErr:   false,
		},
		{
			caseName:  "異常系: 開始日がゼロ値の場合",
			id:        id.NewSeasonID(),
			name:      "A2b",
			startDate: time.Time{},
			endDate:   nil,
			wantErr:   true,
		},
		{
			caseName:  "異常系: 開始日が終了日より後の場合",
			id:        id.NewSeasonID(),
			name:      "A2b",
			startDate: time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC),
			endDate:   &[]time.Time{time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)}[0],
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			season, err := entity.NewSeason(tt.id, tt.name, tt.startDate, tt.endDate)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
				assert.Nil(t, season, "expected season to be nil when error occurs")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.NotNil(t, season, "expected season to be not nil")
			assert.Equal(t, tt.id.String(), season.ID(), "ID does not match")
			assert.Equal(t, tt.name, season.Name(), "Name does not match")
			assert.Equal(t, tt.startDate, season.StartDate(), "StartDate does not match")
			if tt.endDate != nil {
				assert.NotNil(t, season.EndDate(), "EndDate should not be nil")
				assert.Equal(t, *tt.endDate, *season.EndDate(), "EndDate does not match")
				assert.False(t, season.IsActive(), "season with end date should not be active")
			} else {
				assert.Nil(t, season.EndDate(), "EndDate should be nil")
				assert.True(t, season.IsActive(), "season without end date should be active")
			}
		})
	}
}

func TestSeason_IsActive(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName string
		endDate  *time.Time
		want     bool
	}{
		{
			caseName: "正常系: 終了日がnilの場合はアクティブ",
			endDate:  nil,
			want:     true,
		},
		{
			caseName: "正常系: 終了日が設定されている場合は非アクティブ",
			endDate:  &[]time.Time{time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)}[0],
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			seasonID := id.NewSeasonID()
			name := testSeasonName
			startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
			season, _ := entity.NewSeason(seasonID, name, startDate, tt.endDate)

			// Act
			got := season.IsActive()

			// Assert
			assert.Equal(t, tt.want, got, "IsActive does not match expected value")
		})
	}
}

func TestSeason_End(t *testing.T) {
	t.Parallel()

	t.Run("正常系: アクティブなシーズンが終了される", func(t *testing.T) {
		t.Parallel()

		// Arrange
		seasonID := id.NewSeasonID()
		name := testSeasonName
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
		season, _ := entity.NewSeason(seasonID, name, startDate, nil)

		// Act
		err := season.End(endDate)

		// Assert
		assert.NoError(t, err, "unexpected error occurred")
		assert.False(t, season.IsActive(), "season should not be active after ending")
		assert.NotNil(t, season.EndDate(), "EndDate should not be nil after ending")
		assert.Equal(t, endDate, *season.EndDate(), "EndDate does not match")
	})

	t.Run("異常系: 既に終了しているシーズンは終了できない", func(t *testing.T) {
		t.Parallel()

		// Arrange
		seasonID := id.NewSeasonID()
		name := testSeasonName
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		endDate := time.Date(2024, 6, 30, 23, 59, 59, 0, time.UTC)
		season, _ := entity.NewSeason(seasonID, name, startDate, &endDate)
		newEndDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)

		// Act
		err := season.End(newEndDate)

		// Assert
		assert.Error(t, err, "expected error but got none")
		assert.Equal(t, endDate, *season.EndDate(), "EndDate should not change")
	})

	t.Run("異常系: 開始日より前の終了日は設定できない", func(t *testing.T) {
		t.Parallel()

		// Arrange
		seasonID := id.NewSeasonID()
		name := testSeasonName
		startDate := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
		season, _ := entity.NewSeason(seasonID, name, startDate, nil)
		endDate := time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC)

		// Act
		err := season.End(endDate)

		// Assert
		assert.Error(t, err, "expected error but got none")
		assert.True(t, season.IsActive(), "season should remain active")
		assert.Nil(t, season.EndDate(), "EndDate should remain nil")
	})

	t.Run("異常系: ゼロ値の終了日は設定できない", func(t *testing.T) {
		t.Parallel()

		// Arrange
		seasonID := id.NewSeasonID()
		name := testSeasonName
		startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
		season, _ := entity.NewSeason(seasonID, name, startDate, nil)

		// Act
		err := season.End(time.Time{})

		// Assert
		assert.Error(t, err, "expected error but got none")
		assert.True(t, season.IsActive(), "season should remain active")
		assert.Nil(t, season.EndDate(), "EndDate should remain nil")
	})
}
