package handler_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"poketier/apps/season/internal/application/usecase"
	"poketier/apps/season/internal/presentation/handler"
	"poketier/apps/season/internal/presentation/response"
	"poketier/pkg/errs"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestListSeasonsHandler_Handle(t *testing.T) {
	t.Parallel()

	gin.SetMode(gin.TestMode)

	tests := []struct {
		caseName       string
		mockSetup      func(*MockListSeasonsUseCase)
		expectedStatus int
		expectedBody   interface{}
	}{
		{
			caseName: "正常系: シーズン一覧が正常に取得される",
			mockSetup: func(mockUC *MockListSeasonsUseCase) {
				startDate := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
				endDate := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC)
				result := &usecase.ListSeasonsResult{
					Seasons: []usecase.LSSeason{
						{
							SeasonID:  "season-1",
							Name:      "2024 シーズン",
							StartDate: startDate,
							EndDate:   &endDate,
							IsActive:  true,
						},
						{
							SeasonID:  "season-2",
							Name:      "2023 シーズン",
							StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
							EndDate:   nil,
							IsActive:  false,
						},
					},
				}
				mockUC.EXPECT().Execute(gomock.Any()).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: response.ListSeasonsResponse{
				Total: 2,
				Seasons: []response.LSSeason{
					{
						SeasonID:  "season-1",
						Name:      "2024 シーズン",
						StartDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   func() *time.Time { t := time.Date(2024, 12, 31, 23, 59, 59, 0, time.UTC); return &t }(),
						IsActive:  true,
					},
					{
						SeasonID:  "season-2",
						Name:      "2023 シーズン",
						StartDate: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
						EndDate:   nil,
						IsActive:  false,
					},
				},
			},
		},
		{
			caseName: "正常系: 空のシーズン一覧が返される",
			mockSetup: func(mockUC *MockListSeasonsUseCase) {
				result := &usecase.ListSeasonsResult{
					Seasons: []usecase.LSSeason{},
				}
				mockUC.EXPECT().Execute(gomock.Any()).Return(result, nil)
			},
			expectedStatus: http.StatusOK,
			expectedBody: response.ListSeasonsResponse{
				Total:   0,
				Seasons: []response.LSSeason{},
			},
		},
		{
			caseName: "異常系: UseCaseでエラーが発生した場合",
			mockSetup: func(mockUC *MockListSeasonsUseCase) {
				mockUC.EXPECT().Execute(gomock.Any()).Return(nil, errors.New("usecase error"))
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody: errs.ErrorResponse{
				Title:  "Internal Server Error",
				Status: http.StatusInternalServerError,
				Detail: "An internal server error occurred.",
			},
		},
		{
			caseName: "異常系: UseCaseでドメインエラーが発生した場合",
			mockSetup: func(mockUC *MockListSeasonsUseCase) {
				domainErr := errs.NewNotFoundError("seasons not found", nil)
				mockUC.EXPECT().Execute(gomock.Any()).Return(nil, domainErr)
			},
			expectedStatus: http.StatusNotFound,
			expectedBody: errs.ErrorResponse{
				Title:  "Not Found",
				Status: http.StatusNotFound,
				Detail: "The requested resource was not found.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.caseName, func(t *testing.T) {
			t.Parallel()

			// Arrange
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			mockUC := NewMockListSeasonsUseCase(ctrl)
			tt.mockSetup(mockUC)

			handler := handler.NewListSeasonsHandler(mockUC)

			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request = httptest.NewRequest(http.MethodGet, "/seasons", nil)
			c.Request = c.Request.WithContext(context.Background())

			// Act
			handler.Handle(c)

			// Assert
			assert.Equal(t, tt.expectedStatus, w.Code, "status code should match expected")

			var actualBody interface{}
			err := json.Unmarshal(w.Body.Bytes(), &actualBody)
			assert.NoError(t, err, "response body should be valid JSON")

			expectedJSON, err := json.Marshal(tt.expectedBody)
			assert.NoError(t, err, "expected body should be marshallable to JSON")

			var expectedBodyMap interface{}
			err = json.Unmarshal(expectedJSON, &expectedBodyMap)
			assert.NoError(t, err, "expected body should be valid JSON")

			assert.Equal(t, expectedBodyMap, actualBody, "response body should match expected")
		})
	}
}

func TestNewListSeasonsHandler(t *testing.T) {
	t.Parallel()

	t.Run("正常系: ハンドラーが正常に作成される", func(t *testing.T) {
		t.Parallel()

		// Arrange
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockUC := NewMockListSeasonsUseCase(ctrl)

		// Act
		handler := handler.NewListSeasonsHandler(mockUC)

		// Assert
		assert.NotNil(t, handler, "handler should not be nil")
	})
}
