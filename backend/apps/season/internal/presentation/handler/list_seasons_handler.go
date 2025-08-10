package handler

import (
	"context"
	"net/http"
	"poketier/apps/season/internal/application/usecase"
	"poketier/apps/season/internal/presentation/response"
	"poketier/pkg/errs"

	"github.com/gin-gonic/gin"
)

type listSeasonsHandler struct {
	uc ListSeasonsUseCase
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

	ctx.JSON(http.StatusOK, response.NewListSeasonsResponse(result))
}
