package postgresRepository

import (
	"finances-api/db"
	"finances-api/repositories"
)

type PostgresRepository struct {
	PgProducts
	PgUserProducts
	PgGateway
	PgTransactions
	PgInvoices
}

func NewPostgresRepository(db db.Database) repositories.FinancialRepository {
	return PostgresRepository{
		PgProducts:     NewPgProductsRepository(db),
		PgUserProducts: NewPgUserProductsRepository(db),
		PgGateway:      NewPgGatewayRepository(db),
		PgTransactions: NewPgTransactionsRepository(db),
		PgInvoices:     NewPgInvoicesRepository(db),
	}
}
