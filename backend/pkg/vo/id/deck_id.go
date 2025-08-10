package id

import "github.com/google/uuid"

// deckEntity はDeck集約のマーカー型
type deckEntity struct{}

// DeckID はデッキの一意識別子
type DeckID = ID[deckEntity]

// NewDeckID は新しいDeckIDを生成
func NewDeckID() DeckID {
	return new[deckEntity]()
}

// DeckIDFromString は文字列からDeckIDを再作成
func DeckIDFromString(s string) (DeckID, error) {
	return fromString[deckEntity](s)
}

// DeckIDFromUUID はuuid.UUIDからDeckIDを作成
func DeckIDFromUUID(u uuid.UUID) DeckID {
	return fromUUID[deckEntity](u)
}
