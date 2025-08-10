package entity

import (
	"errors"
	"time"

	"poketier/pkg/vo/id"
)

// Season はティアリストの環境期間を表すエンティティ
type Season struct {
	id        id.SeasonID
	name      string
	startDate time.Time
	endDate   time.Time
}

// NewSeason は新しいSeasonインスタンスを作成する
func NewSeason(id id.SeasonID, name string, startDate time.Time, endDate time.Time) (*Season, error) {
	season := &Season{
		id:        id,
		name:      name,
		startDate: startDate,
		endDate:   endDate,
	}

	if err := season.validate(); err != nil {
		return nil, err
	}

	return season, nil
}

// ID はSeasonのIDを返す
func (s *Season) ID() id.SeasonID {
	return s.id
}

// Name はSeasonの名前を返す
func (s *Season) Name() string {
	return s.name
}

// IsActive はSeasonがアクティブかどうかを返す（現在日時が期間内かどうか）
func (s *Season) IsActive() bool {
	now := time.Now()
	return !now.Before(s.startDate) && !now.After(s.endDate)
}

// StartDate はSeasonの開始日を返す
func (s *Season) StartDate() time.Time {
	return s.startDate
}

// EndDate はSeasonの終了日を返す
func (s *Season) EndDate() time.Time {
	return s.endDate
}

// End はSeasonの終了日を変更する
func (s *Season) End(endDate time.Time) error {
	if err := s.validEndDate(endDate); err != nil {
		return err
	}

	s.endDate = endDate

	return nil
}

// validate は全体のバリデーションを実行する
func (s *Season) validate() error {
	if err := s.validName(); err != nil {
		return err
	}

	if err := s.validStartDate(); err != nil {
		return err
	}

	if err := s.validEndDate(s.endDate); err != nil {
		return err
	}

	return nil
}

// validName は名前のバリデーションを行う
func (s *Season) validName() error {
	if s.name == "" {
		return errors.New("name cannot be empty")
	}
	var maxNameLength = 5
	if len(s.name) > maxNameLength {
		return errors.New("name must be 5 characters or less")
	}
	return nil
}

// validStartDate は開始日のバリデーションを行う
func (s *Season) validStartDate() error {
	if s.startDate.IsZero() {
		return errors.New("start date cannot be zero")
	}
	return nil
}

// validEndDate は終了日のバリデーションを行う
func (s *Season) validEndDate(endDate time.Time) error {
	if endDate.IsZero() {
		return errors.New("end date cannot be zero")
	}

	if endDate.Before(s.startDate) {
		return errors.New("end date must be after start date")
	}

	return nil
}
