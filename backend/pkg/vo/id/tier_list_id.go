package id

// tierListEntity はTierList集約のマーカー型
type tierListEntity struct{}

// TierListID はティアリストの一意識別子
type TierListID = ID[tierListEntity]

// NewTierListID は新しいTierListIDを生成
func NewTierListID() TierListID {
	return New[tierListEntity]()
}

// ReNewTierListID は文字列からTierListIDを再作成
func ReNewTierListID(s string) (TierListID, error) {
	return ReNew[tierListEntity](s)
}
