package errs

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorResponse はRFC 9457に準拠したエラーレスポンスの構造体です。ただし、Typeフィールドは含みません。
// https://www.rfc-editor.org/rfc/rfc9457.html
type ErrorResponse struct {
	Title  string    `json:"title"`            // 人間が読めるエラーの要約
	Status int       `json:"status"`           // HTTPステータスコード
	Detail string    `json:"detail"`           // エラーの詳細説明
	Errors *[]string `json:"errors,omitempty"` // バリデーションエラーメッセージの配列（オプション）
}

// エラー種別とHTTPステータスコード、クライアント向けメッセージのマッピング
var errorMappings = map[error]struct {
	status int
	title  string
	detail string
}{
	ErrBadRequest:   {http.StatusBadRequest, "Bad Request", "The request is invalid."},
	ErrUnauthorized: {http.StatusUnauthorized, "Unauthorized", "Authentication is required."},
	ErrForbidden:    {http.StatusForbidden, "Forbidden", "You do not have permission to perform this action."},
	ErrNotFound:     {http.StatusNotFound, "Not Found", "The requested resource was not found."},
	ErrTimeout:      {http.StatusRequestTimeout, "Request Timeout", "The request timed out."},
	ErrConflict:     {http.StatusConflict, "Conflict", "A resource conflict occurred."},
}

// HandleError はエラーを受け取り、適切なHTTPレスポンスを返します
func HandleError(ctx *gin.Context, err error) {
	// Ginのエラーコンテキストにエラーを追加（ミドルウェアでログ出力される）
	ctx.Error(err)

	var domainErr *DomainError
	if errors.As(err, &domainErr) {
		// ドメインエラーの場合、マッピングを使用
		if mapping, exists := errorMappings[domainErr.Type]; exists {
			response := ErrorResponse{
				Title:  mapping.title,
				Status: mapping.status,
				Detail: mapping.detail,
			}
			ctx.JSON(mapping.status, response)
			return
		}
	}

	// 未知のエラーの場合は内部サーバーエラーとして扱う
	response := ErrorResponse{
		Title:  "Internal Server Error",
		Status: http.StatusInternalServerError,
		Detail: "An internal server error occurred.",
	}
	ctx.JSON(http.StatusInternalServerError, response)
}

// HandleValidationError はバリデーションエラーを処理します
func HandleValidationError(ctx *gin.Context, errors []error) {
	errorMessages := make([]string, len(errors))
	for i, err := range errors {
		errorMessages[i] = err.Error()
	}

	response := ErrorResponse{
		Title:  "Validation Error",
		Status: http.StatusUnprocessableEntity,
		Detail: "The provided data is invalid.",
		Errors: &errorMessages,
	}
	ctx.JSON(http.StatusUnprocessableEntity, response)
}
