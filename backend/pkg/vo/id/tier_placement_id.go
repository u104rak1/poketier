package id

// tierPlacementEntity はTierPlacement集約のマーカー型
type tierPlacementEntity struct{}

// TierPlacementID はティア配置の一意識別子
type TierPlacementID = ID[tierPlacementEntity]

// NewTierPlacementID は新しいTierPlacementIDを生成
func NewTierPlacementID() TierPlacementID {
	return New[tierPlacementEntity]()
}

// ReNewTierPlacementID は文字列からTierPlacementIDを再作成
func ReNewTierPlacementID(s string) (TierPlacementID, error) {
	return ReNew[tierPlacementEntity](s)
}
