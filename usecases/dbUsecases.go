package usecases

import (
	"finances-api/db"
	"finances-api/repositories"
	postgresRepository "finances-api/repositories/db/postgres"
)

type DbUsecase struct {
	repositories.FinancialRepository
}

func NewDbUsecase(db db.Database) *DbUsecase {

	repository := postgresRepository.NewPostgresRepository(db)

	return &DbUsecase{
		repository,
	}
}
