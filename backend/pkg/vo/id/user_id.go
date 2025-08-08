package id

// userEntity はUser集約のマーカー型
type userEntity struct{}

// UserID はユーザーの一意識別子
type UserID = ID[userEntity]

// NewUserID は新しいUserIDを生成
func NewUserID() UserID {
	return New[userEntity]()
}

// ReNewUserID は文字列からUserIDを再作成
func ReNewUserID(s string) (UserID, error) {
	return ReNew[userEntity](s)
}
