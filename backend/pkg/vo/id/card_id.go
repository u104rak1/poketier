package id

// cardEntity はCard集約のマーカー型
type cardEntity struct{}

// CardID はカードの一意識別子
type CardID = ID[cardEntity]

// NewCardID は新しいCardIDを生成
func NewCardID() CardID {
	return New[cardEntity]()
}

// ReNewCardID は文字列からCardIDを再作成
func ReNewCardID(s string) (CardID, error) {
	return ReNew[cardEntity](s)
}
