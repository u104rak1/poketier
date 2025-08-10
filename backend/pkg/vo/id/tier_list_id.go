package id

import "github.com/google/uuid"

// tierListEntity はTierList集約のマーカー型
type tierListEntity struct{}

// TierListID はティアリストの一意識別子
type TierListID = ID[tierListEntity]

// NewTierListID は新しいTierListIDを生成
func NewTierListID() TierListID {
	return new[tierListEntity]()
}

// TierListIDFromString は文字列からTierListIDを再作成
func TierListIDFromString(s string) (TierListID, error) {
	return fromString[tierListEntity](s)
}

// TierListIDFromUUID はuuid.UUIDからTierListIDを作成
func TierListIDFromUUID(u uuid.UUID) TierListID {
	return fromUUID[tierListEntity](u)
}
