package handler

import (
	"context"
	"net/http"
	"poketier/apps/tierlist/internal/usecase"
	"poketier/pkg/errs"

	"github.com/gin-gonic/gin"
)

type listSeasonsHandler struct {
	uc ListSeasonsUseCase
}

// Response構造体
type listSeasonsResponse struct {
	Total   int                `json:"total"`
	Seasons []usecase.LSSeason `json:"seasons"`
}

type ListSeasonsUseCase interface {
	Execute(ctx context.Context) (*usecase.ListSeasonsResult, error)
}

func NewListSeasonsHandler(uc ListSeasonsUseCase) *listSeasonsHandler {
	return &listSeasonsHandler{
		uc: uc,
	}
}

func (h *listSeasonsHandler) Handle(ctx *gin.Context) {
	result, err := h.uc.Execute(ctx.Request.Context())
	if err != nil {
		errs.HandleError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, h.toResponse(result))
}

func (h *listSeasonsHandler) toResponse(result *usecase.ListSeasonsResult) listSeasonsResponse {
	return listSeasonsResponse{
		Total:   len(result.Seasons),
		Seasons: result.Seasons,
	}
}
