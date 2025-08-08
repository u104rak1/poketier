package id

// deckEntity はDeck集約のマーカー型
type deckEntity struct{}

// DeckID はデッキの一意識別子
type DeckID = ID[deckEntity]

// NewDeckID は新しいDeckIDを生成
func NewDeckID() DeckID {
	return New[deckEntity]()
}

// ReNewDeckID は文字列からDeckIDを再作成
func ReNewDeckID(s string) (DeckID, error) {
	return ReNew[deckEntity](s)
}
