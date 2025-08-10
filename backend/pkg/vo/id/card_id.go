package id

import "github.com/google/uuid"

// cardEntity はCard集約のマーカー型
type cardEntity struct{}

// CardID はカードの一意識別子
type CardID = ID[cardEntity]

// NewCardID は新しいCardIDを生成
func NewCardID() CardID {
	return new[cardEntity]()
}

// CardIDFromString は文字列からCardIDを再作成
func CardIDFromString(s string) (CardID, error) {
	return fromString[cardEntity](s)
}

// CardIDFromUUID はuuid.UUIDからCardIDを作成
func CardIDFromUUID(u uuid.UUID) CardID {
	return fromUUID[cardEntity](u)
}
