package id

// seasonEntity はSeason集約のマーカー型
type seasonEntity struct{}

// SeasonID はシーズンの一意識別子
type SeasonID = ID[seasonEntity]

// NewSeasonID は新しいSeasonIDを生成
func NewSeasonID() SeasonID {
	return New[seasonEntity]()
}

// ReNewSeasonID は文字列からSeasonIDを再作成
func ReNewSeasonID(s string) (SeasonID, error) {
	return ReNew[seasonEntity](s)
}
