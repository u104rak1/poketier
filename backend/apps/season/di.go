//go:build wireinject
// +build wireinject

package season

import (
	"poketier/apps/season/internal/application/usecase"
	"poketier/apps/season/internal/infrastructure/repository"
	"poketier/apps/season/internal/presentation/handler"
	"poketier/sqlc/db"

	"github.com/google/wire"
)

// InitializeListSeasonsHandler はListSeasonsHandlerとその依存関係を初期化します
func InitializeListSeasonsHandler(queries db.Querier) *handler.ListSeasonsHandler {
	wire.Build(
		// Repository provider
		wire.Bind(new(repository.SeasonQuerier), new(db.Querier)),
		repository.NewSeasonRepository,
		wire.Bind(new(usecase.LSSeasonRepository), new(*repository.SeasonRepository)),

		// Usecase provider
		usecase.NewListSeasonsUsecase,
		wire.Bind(new(handler.ListSeasonsUseCase), new(*usecase.ListSeasonsUsecase)),

		// Handler provider
		handler.NewListSeasonsHandler,
	)
	return &handler.ListSeasonsHandler{}
}
