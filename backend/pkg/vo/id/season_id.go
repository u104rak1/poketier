package id

import "github.com/google/uuid"

// seasonEntity はSeason集約のマーカー型
type seasonEntity struct{}

// SeasonID はシーズンの一意識別子
type SeasonID = ID[seasonEntity]

// NewSeasonID は新しいSeasonIDを生成
func NewSeasonID() SeasonID {
	return new[seasonEntity]()
}

// SeasonIDFromString は文字列からSeasonIDを再作成
func SeasonIDFromString(s string) (SeasonID, error) {
	return fromString[seasonEntity](s)
}

// SeasonIDFromUUID はuuid.UUIDからSeasonIDを作成
func SeasonIDFromUUID(u uuid.UUID) SeasonID {
	return fromUUID[seasonEntity](u)
}
