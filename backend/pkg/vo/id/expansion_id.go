package id

// expansionEntity はExpansion集約のマーカー型
type expansionEntity struct{}

// ExpansionID は拡張パックの一意識別子
type ExpansionID = ID[expansionEntity]

// NewExpansionID は新しいExpansionIDを生成
func NewExpansionID() ExpansionID {
	return New[expansionEntity]()
}

// ReNewExpansionID は文字列からExpansionIDを再作成
func ReNewExpansionID(s string) (ExpansionID, error) {
	return ReNew[expansionEntity](s)
}
