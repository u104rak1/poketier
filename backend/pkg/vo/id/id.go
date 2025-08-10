// Package id はドメインエンティティのID値オブジェクトを提供します
package id

import (
	"fmt"

	"github.com/google/uuid"
)

// ID はUUID v7を使用するジェネリクスなID値オブジェクト
type ID[T comparable] struct {
	value uuid.UUID
}

// new はUUID v7を生成して新しいIDを作成
func new[T comparable]() ID[T] {
	return ID[T]{value: uuid.Must(uuid.NewV7())}
}

// fromUUID は内部的にUUIDからIDを作成（パッケージ内使用）
func fromUUID[T comparable](u uuid.UUID) ID[T] {
	return ID[T]{value: u}
}

// fromString は文字列からIDを再作成
func fromString[T comparable](s string) (ID[T], error) {
	parsedUUID, err := uuid.Parse(s)
	if err != nil {
		return ID[T]{}, fmt.Errorf("invalid UUID format: %w", err)
	}
	return ID[T]{value: parsedUUID}, nil
}

// String はUUIDの文字列表現を返す
func (id ID[T]) String() string {
	return id.value.String()
}

// UUID はUUIDオブジェクトを返す
func (id ID[T]) UUID() uuid.UUID {
	return id.value
}

// Equals は別のIDとの等価性を判定
func (id ID[T]) Equals(other ID[T]) bool {
	return id.value == other.value
}
