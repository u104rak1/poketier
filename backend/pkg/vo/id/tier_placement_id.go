package id

import "github.com/google/uuid"

// tierPlacementEntity はTierPlacement集約のマーカー型
type tierPlacementEntity struct{}

// TierPlacementID はティア配置の一意識別子
type TierPlacementID = ID[tierPlacementEntity]

// NewTierPlacementID は新しいTierPlacementIDを生成
func NewTierPlacementID() TierPlacementID {
	return new[tierPlacementEntity]()
}

// TierPlacementIDFromString は文字列からTierPlacementIDを再作成
func TierPlacementIDFromString(s string) (TierPlacementID, error) {
	return fromString[tierPlacementEntity](s)
}

// TierPlacementIDFromUUID はuuid.UUIDからTierPlacementIDを作成
func TierPlacementIDFromUUID(u uuid.UUID) TierPlacementID {
	return fromUUID[tierPlacementEntity](u)
}
