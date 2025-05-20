package usecases

import (
	"finances-api/db"
	"finances-api/repositories"
	postgresRepository "finances-api/repositories/db/postgres"
)

type DbUsecase struct {
	Repository repositories.FinancialRepository
}

func NewPgUsecase(db db.Database) *DbUsecase {

	repository := postgresRepository.NewPostgresRepository(db)

	return &DbUsecase{
		Repository: repository,
	}
}
