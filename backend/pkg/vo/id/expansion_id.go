package id

import "github.com/google/uuid"

// expansionEntity はExpansion集約のマーカー型
type expansionEntity struct{}

// ExpansionID は拡張パックの一意識別子
type ExpansionID = ID[expansionEntity]

// NewExpansionID は新しいExpansionIDを生成
func NewExpansionID() ExpansionID {
	return new[expansionEntity]()
}

// ExpansionIDFromString は文字列からExpansionIDを再作成
func ExpansionIDFromString(s string) (ExpansionID, error) {
	return fromString[expansionEntity](s)
}

// ExpansionIDFromUUID はuuid.UUIDからExpansionIDを作成
func ExpansionIDFromUUID(u uuid.UUID) ExpansionID {
	return fromUUID[expansionEntity](u)
}
