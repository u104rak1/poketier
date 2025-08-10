package response

import (
	"poketier/apps/season/internal/application/usecase"
	"time"
)

type ListSeasonsResponse struct {
	Total   int        `json:"total"`
	Seasons []LSSeason `json:"seasons"`
}

type LSSeason struct {
	SeasonID  string    `json:"season_id"`
	Name      string    `json:"name"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsActive  bool      `json:"is_active"`
}

func NewListSeasonsResponse(result *usecase.ListSeasonsResult) ListSeasonsResponse {
	seasons := make([]LSSeason, len(result.Seasons))
	for i, s := range result.Seasons {
		seasons[i] = LSSeason{
			SeasonID:  s.SeasonID,
			Name:      s.Name,
			StartDate: s.StartDate,
			EndDate:   s.EndDate,
			IsActive:  s.IsActive,
		}
	}
	return ListSeasonsResponse{
		Total:   len(result.Seasons),
		Seasons: seasons,
	}
}
