package id

import "github.com/google/uuid"

// userEntity はUser集約のマーカー型
type userEntity struct{}

// UserID はユーザーの一意識別子
type UserID = ID[userEntity]

// NewUserID は新しいUserIDを生成
func NewUserID() UserID {
	return new[userEntity]()
}

// UserIDFromString は文字列からUserIDを再作成
func UserIDFromString(s string) (UserID, error) {
	return fromString[userEntity](s)
}

// UserIDFromUUID はuuid.UUIDからUserIDを作成
func UserIDFromUUID(u uuid.UUID) UserID {
	return fromUUID[userEntity](u)
}
