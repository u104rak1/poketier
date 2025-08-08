package id

import (
	"database/sql/driver"
	"fmt"

	"github.com/google/uuid"
)

// ID はUUID v7を使用するジェネリクスなID値オブジェクト
type ID[T comparable] struct {
	value string
}

// New はUUID v7を生成して新しいIDを作成
func New[T comparable]() ID[T] {
	return ID[T]{value: uuid.Must(uuid.NewV7()).String()}
}

// ReNew は文字列からIDを再作成
func ReNew[T comparable](s string) (ID[T], error) {
	_, err := uuid.Parse(s)
	if err != nil {
		return ID[T]{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return ID[T]{value: s}, nil
}

// String はUUIDの文字列表現を返す
func (id ID[T]) String() string {
	return id.value
}

// Equals は別のIDとの等価性を判定
func (id ID[T]) Equals(other ID[T]) bool {
	return id.value == other.value
}

// Value はdatabase/sql/driver.Valuerインターフェースの実装
func (id ID[T]) Value() (driver.Value, error) {
	return id.value, nil
}

// Scan はsql.Scannerインターフェースの実装
func (id *ID[T]) Scan(value interface{}) error {
	if value == nil {
		*id = ID[T]{}
		return nil
	}

	switch v := value.(type) {
	case string:
		_, err := uuid.Parse(v)
		if err != nil {
			return fmt.Errorf("invalid UUID format: %w", err)
		}
		*id = ID[T]{value: v}
		return nil
	case []byte:
		s := string(v)
		_, err := uuid.Parse(s)
		if err != nil {
			return fmt.Errorf("invalid UUID format: %w", err)
		}
		*id = ID[T]{value: s}
		return nil
	default:
		return fmt.Errorf("cannot scan %T into ID", value)
	}
}
