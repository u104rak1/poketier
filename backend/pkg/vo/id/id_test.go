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

func TestID_ReNew(t *testing.T) {
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
			got, err := id.ReNewSeasonID(tt.input)

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
		id1, _ := id.ReNewSeasonID(uuidStr)
		id2, _ := id.ReNewSeasonID(uuidStr)

		// Act
		got := id1.Equals(id2)

		// Assert
		assert.True(t, got, "IDs created from same UUID string should be equal")
	})
}

func TestID_Value(t *testing.T) {
	t.Parallel()

	t.Run("正常系: 生成されたIDからValueが取得できる事", func(t *testing.T) {
		t.Parallel()

		// Arrange
		testID := id.NewSeasonID()

		// Act
		got, err := testID.Value()

		// Assert
		assert.NoError(t, err, "Value should not return error")
		assert.Equal(t, testID.String(), got, "Value should match ID string")
	})
}

func TestID_Scan(t *testing.T) {
	t.Parallel()

	originalID := id.NewSeasonID()
	uuidStr := originalID.String()

	tests := []struct {
		caseName    string
		input       interface{}
		wantErr     bool
		errContains string
		want        string
	}{
		{
			caseName: "正常系: 文字列からScanできる事",
			input:    uuidStr,
			wantErr:  false,
			want:     uuidStr,
		},
		{
			caseName: "正常系: バイト配列からScanできる事",
			input:    []byte(uuidStr),
			wantErr:  false,
			want:     uuidStr,
		},
		{
			caseName: "正常系: nilからScanして空文字列になる事",
			input:    nil,
			wantErr:  false,
			want:     "",
		},
		{
			caseName:    "異常系: 無効なUUID文字列が渡された場合",
			input:       "invalid-uuid",
			wantErr:     true,
			errContains: "invalid UUID format",
		},
		{
			caseName:    "異常系: サポートされていない型が渡された場合",
			input:       123,
			wantErr:     true,
			errContains: "cannot scan",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var scannedID id.SeasonID

			// Act
			err := scannedID.Scan(tt.input)

			// Assert
			if tt.wantErr {
				assert.Error(t, err, "expected error but got none")
				assert.Contains(t, err.Error(), tt.errContains, "error message should contain expected text")
				return
			}
			assert.NoError(t, err, "unexpected error occurred")
			assert.Equal(t, tt.want, scannedID.String(), "scanned ID should match expected value")
		})
	}
}
