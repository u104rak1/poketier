package id_test

import (
	"testing"

	"poketier/pkg/vo/id"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestID_New(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName  string
		generator func() string
	}{
		{
			caseName:  "正常系: SeasonIDが空文字列以外で生成される事",
			generator: func() string { return id.NewSeasonID().String() },
		},
		{
			caseName:  "正常系: CardIDが空文字列以外で生成される事",
			generator: func() string { return id.NewCardID().String() },
		},
		{
			caseName:  "正常系: DeckIDが空文字列以外で生成される事",
			generator: func() string { return id.NewDeckID().String() },
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			got := tt.generator()

			// Assert
			assert.NotEmpty(t, got, "ID should not be empty")
		})
	}
}

func TestID_FromUUID(t *testing.T) {
	t.Parallel()

	t.Run("正常系: uuid.UUIDからSeasonIDが作成される事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		originalUUID := uuid.Must(uuid.NewV7())

		// Act
		seasonID := id.SeasonIDFromUUID(originalUUID)

		// Assert
		assert.Equal(t, originalUUID, seasonID.UUID(), "created SeasonID should have same UUID")
		assert.Equal(t, originalUUID.String(), seasonID.String(), "created SeasonID string should match UUID string")
	})
}

func TestID_FromString(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName    string
		input       string
		wantErr     bool
		errContains string
	}{
		{
			caseName: "正常系: 有効なUUID文字列からSeasonIDが作成される事",
			input:    uuid.Must(uuid.NewV7()).String(),
			wantErr:  false,
		},
		{
			caseName:    "異常系: 無効なUUID文字列が渡された場合",
			input:       "invalid-uuid",
			wantErr:     true,
			errContains: "invalid UUID format",
		},
		{
			caseName:    "異常系: 空文字列が渡された場合",
			input:       "",
			wantErr:     true,
			errContains: "invalid UUID format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Act
			got, err := id.SeasonIDFromString(tt.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
				assert.Contains(t, err.Error(), tt.errContains, "error message should contain expected text")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.Equal(t, tt.input, got.String(), "regenerated ID should match input")
		})
	}
}

func TestID_String(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 生成されたIDから文字列が取得できる事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		testID := id.NewSeasonID()

		// Act
		got := testID.String()

		// Assert
		assert.NotEmpty(t, got, "String should not be empty")
		assert.Contains(t, got, "-", "String should be UUID format with hyphens")
	})
}

func TestID_UUID(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 生成されたIDからUUIDが取得できる事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		testID := id.NewSeasonID()

		// Act
		got := testID.UUID()

		// Assert
		assert.NotEqual(t, uuid.Nil, got, "UUID should not be nil")
		assert.Equal(t, testID.String(), got.String(), "UUID string should match ID string")
	})
}

func TestID_Equals(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 同じIDオブジェクトが等価である事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		id1 := id.NewSeasonID()
		id2 := id1 // 同じIDオブジェクト

		// Act
		got := id1.Equals(id2)

		// Assert
		assert.True(t, got, "same ID object should be equal to itself")
	})

	t.Run("正常系: 異なるIDオブジェクトが等価でない事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		id1 := id.NewSeasonID()
		id2 := id.NewSeasonID()

		// Act
		got := id1.Equals(id2)

		// Assert
		assert.False(t, got, "different ID objects should not be equal")
	})

	t.Run("正常系: 同じUUID文字列から作成されたIDが等価である事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		uuidStr := uuid.Must(uuid.NewV7()).String()
		id1, _ := id.SeasonIDFromString(uuidStr)
		id2, _ := id.SeasonIDFromString(uuidStr)

		// Act
		got := id1.Equals(id2)

		// Assert
		assert.True(t, got, "IDs created from same UUID string should be equal")
	})
}
