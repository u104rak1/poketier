package errs

import "errors"

// ドメインエラーの種類を定義
var (
	// リクエストエラー（400）
	ErrBadRequest = errors.New("bad request")

	// 認証エラー（401）
	ErrUnauthorized = errors.New("unauthorized")

	// 認可エラー（403）
	ErrForbidden = errors.New("forbidden")

	// リソースが見つからない（404）
	ErrNotFound = errors.New("not found")

	// タイムアウト（408）
	ErrTimeout = errors.New("request timeout")

	// 競合状態（409）
	ErrConflict = errors.New("conflict")

	// 内部サーバーエラー（500）
	ErrInternal = errors.New("internal server error")
)

// ドメイン固有のエラー型
type DomainError struct {
	Type    error  // ドメインエラー種別
	Message string // 内部用の詳細メッセージ
	Cause   error  // 元のエラー
}

func (e *DomainError) Error() string {
	if e.Cause != nil {
		return e.Message + ": " + e.Cause.Error()
	}
	return e.Message
}

func (e *DomainError) Unwrap() error {
	return e.Cause
}

func NewValidationError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrBadRequest,
		Message: message,
		Cause:   cause,
	}
}

func NewUnauthorizedError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrUnauthorized,
		Message: message,
		Cause:   cause,
	}
}

func NewForbiddenError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrForbidden,
		Message: message,
		Cause:   cause,
	}
}

func NewNotFoundError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrNotFound,
		Message: message,
		Cause:   cause,
	}
}

func NewTimeoutError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrTimeout,
		Message: message,
		Cause:   cause,
	}
}

func NewConflictError(message string, cause error) *DomainError {
	return &DomainError{
		Type:    ErrConflict,
		Message: message,
		Cause:   cause,
	}
}
