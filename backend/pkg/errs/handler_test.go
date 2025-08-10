package errs_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"poketier/pkg/errs"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandleError(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	tests := []struct {
		caseName           string
		setupError         func() error
		expectedStatusCode int
		expectedResponse   errs.ErrorResponse
	}{
		{
			caseName: "正常系_BadRequestドメインエラー",
			setupError: func() error {
				return errs.NewValidationError("invalid input", errors.New("field required"))
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedResponse: errs.ErrorResponse{
				Title:  "Bad Request",
				Status: http.StatusBadRequest,
				Detail: "The request is invalid.",
			},
		},
		{
			caseName: "正常系_Unauthorizedドメインエラー",
			setupError: func() error {
				return errs.NewUnauthorizedError("authentication failed", errors.New("invalid token"))
			},
			expectedStatusCode: http.StatusUnauthorized,
			expectedResponse: errs.ErrorResponse{
				Title:  "Unauthorized",
				Status: http.StatusUnauthorized,
				Detail: "Authentication is required.",
			},
		},
		{
			caseName: "正常系_Forbiddenドメインエラー",
			setupError: func() error {
				return errs.NewForbiddenError("access denied", errors.New("insufficient permissions"))
			},
			expectedStatusCode: http.StatusForbidden,
			expectedResponse: errs.ErrorResponse{
				Title:  "Forbidden",
				Status: http.StatusForbidden,
				Detail: "You do not have permission to perform this action.",
			},
		},
		{
			caseName: "正常系_NotFoundドメインエラー",
			setupError: func() error {
				return errs.NewNotFoundError("user not found", errors.New("no rows"))
			},
			expectedStatusCode: http.StatusNotFound,
			expectedResponse: errs.ErrorResponse{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "The requested resource was not found.",
			},
		},
		{
			caseName: "正常系_Timeoutドメインエラー",
			setupError: func() error {
				return errs.NewTimeoutError("request timeout", errors.New("context deadline exceeded"))
			},
			expectedStatusCode: http.StatusRequestTimeout,
			expectedResponse: errs.ErrorResponse{
				Title:  "Request Timeout",
				Status: http.StatusRequestTimeout,
				Detail: "The request timed out.",
			},
		},
		{
			caseName: "正常系_Conflictドメインエラー",
			setupError: func() error {
				return errs.NewConflictError("resource conflict", errors.New("unique constraint"))
			},
			expectedStatusCode: http.StatusConflict,
			expectedResponse: errs.ErrorResponse{
				Title:  "Conflict",
				Status: http.StatusConflict,
				Detail: "A resource conflict occurred.",
			},
		},
		{
			caseName: "正常系_未知のドメインエラー",
			setupError: func() error {
				return &errs.DomainError{
					Type:    errors.New("unknown error"),
					Message: "unknown domain error",
					Cause:   errors.New("some cause"),
				}
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: errs.ErrorResponse{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
				Detail: "An internal server error occurred.",
			},
		},
		{
			caseName: "正常系_一般的なエラー",
			setupError: func() error {
				return errors.New("some general error")
			},
			expectedStatusCode: http.StatusInternalServerError,
			expectedResponse: errs.ErrorResponse{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
				Detail: "An internal server error occurred.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			err := tt.setupError()

			// Act
			errs.HandleError(c, err)

			// Assert
			assert.Equal(t, tt.expectedStatusCode, w.Code, "Status code should match")

			var response errs.ErrorResponse
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err, "Response should be valid JSON")

			assert.Equal(t, tt.expectedResponse.Title, response.Title, "Title should match")
			assert.Equal(t, tt.expectedResponse.Status, response.Status, "Status should match")
			assert.Equal(t, tt.expectedResponse.Detail, response.Detail, "Detail should match")

			// Ginのエラーコンテキストに追加されていることを確認
			assert.Len(t, c.Errors, 1, "Error should be added to Gin context")
		})
	}
}

func TestHandleValidationError(t *testing.T) {
	t.Parallel()

	// Ginのテスト用モードに設定
	gin.SetMode(gin.TestMode)

	t.Run("正常系_バリデーションエラーを処理", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		validationErrors := []error{
			errors.New("field name is required"),
			errors.New("field email is invalid"),
		}

		// Act
		errs.HandleValidationError(c, validationErrors)

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "Status code should be 422")

		var response errs.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")

		assert.Equal(t, "Validation Error", response.Title, "Title should match")
		assert.Equal(t, http.StatusUnprocessableEntity, response.Status, "Status should match")
		assert.Equal(t, "The provided data is invalid.", response.Detail, "Detail should match")

		require.NotNil(t, response.Errors, "Errors field should not be nil")
		assert.Len(t, *response.Errors, 2, "Should have 2 error messages")
		assert.Contains(t, *response.Errors, "field name is required", "Should contain first error")
		assert.Contains(t, *response.Errors, "field email is invalid", "Should contain second error")
	})

	t.Run("正常系_空のバリデーションエラーリスト", func(t *testing.T) {
		t.Parallel()

		// Arrange
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		validationErrors := []error{}

		// Act
		errs.HandleValidationError(c, validationErrors)

		// Assert
		assert.Equal(t, http.StatusUnprocessableEntity, w.Code, "Status code should be 422")

		var response errs.ErrorResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		require.NoError(t, err, "Response should be valid JSON")

		require.NotNil(t, response.Errors, "Errors field should not be nil")
		assert.Len(t, *response.Errors, 0, "Should have 0 error messages")
	})
}
