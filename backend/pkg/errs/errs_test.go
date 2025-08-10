package errs_test

import (
	"errors"
	"testing"

	"poketier/pkg/errs"

	"github.com/stretchr/testify/assert"
)

func TestDomainError_Error(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName string
		setupErr func() *errs.DomainError
		want     string
	}{
		{
			caseName: "正常系_原因エラーがない場合",
			setupErr: func() *errs.DomainError {
				return &errs.DomainError{
					Type:    errs.ErrNotFound,
					Message: "user not found",
					Cause:   nil,
				}
			},
			want: "user not found",
		},
		{
			caseName: "正常系_原因エラーがある場合",
			setupErr: func() *errs.DomainError {
				return &errs.DomainError{
					Type:    errs.ErrNotFound,
					Message: "user not found",
					Cause:   errors.New("database connection failed"),
				}
			},
			want: "user not found: database connection failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			domainErr := tt.setupErr()

			// Act
			got := domainErr.Error()

			// Assert
			assert.Equal(t, tt.want, got, "Error message should match expected value")
		})
	}
}

func TestDomainError_Unwrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		caseName string
		setupErr func() *errs.DomainError
		want     error
	}{
		{
			caseName: "正常系_原因エラーがない場合",
			setupErr: func() *errs.DomainError {
				return &errs.DomainError{
					Type:    errs.ErrNotFound,
					Message: "user not found",
					Cause:   nil,
				}
			},
			want: nil,
		},
		{
			caseName: "正常系_原因エラーがある場合",
			setupErr: func() *errs.DomainError {
				originalErr := errors.New("database connection failed")
				return &errs.DomainError{
					Type:    errs.ErrNotFound,
					Message: "user not found",
					Cause:   originalErr,
				}
			},
			want: errors.New("database connection failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			domainErr := tt.setupErr()

			// Act
			got := domainErr.Unwrap()

			// Assert
			if tt.want == nil {
				assert.Nil(t, got, "Unwrap should return nil when no cause")
			} else {
				assert.Equal(t, tt.want.Error(), got.Error(), "Unwrap should return the cause error")
			}
		})
	}
}

func TestNewValidationError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_バリデーションエラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "invalid input"
		cause := errors.New("field is required")

		// Act
		domainErr := errs.NewValidationError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrBadRequest, domainErr.Type, "Type should be ErrBadRequest")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}

func TestNewUnauthorizedError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_認証エラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "authentication failed"
		cause := errors.New("invalid token")

		// Act
		domainErr := errs.NewUnauthorizedError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrUnauthorized, domainErr.Type, "Type should be ErrUnauthorized")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}

func TestNewForbiddenError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_認可エラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "access denied"
		cause := errors.New("insufficient permissions")

		// Act
		domainErr := errs.NewForbiddenError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrForbidden, domainErr.Type, "Type should be ErrForbidden")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}

func TestNewNotFoundError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_NotFoundエラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "user not found"
		cause := errors.New("no rows in result set")

		// Act
		domainErr := errs.NewNotFoundError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrNotFound, domainErr.Type, "Type should be ErrNotFound")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}

func TestNewTimeoutError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_タイムアウトエラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "request timeout"
		cause := errors.New("context deadline exceeded")

		// Act
		domainErr := errs.NewTimeoutError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrTimeout, domainErr.Type, "Type should be ErrTimeout")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}

func TestNewConflictError(t *testing.T) {
	t.Parallel()

	t.Run("正常系_競合エラーを作成", func(t *testing.T) {
		t.Parallel()

		// Arrange
		message := "resource conflict"
		cause := errors.New("unique constraint violation")

		// Act
		domainErr := errs.NewConflictError(message, cause)

		// Assert
		assert.Equal(t, errs.ErrConflict, domainErr.Type, "Type should be ErrConflict")
		assert.Equal(t, message, domainErr.Message, "Message should match")
		assert.Equal(t, cause, domainErr.Cause, "Cause should match")
	})
}
