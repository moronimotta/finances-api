package usecases

import (
	"finances-api/db"
	"finances-api/messaging"
	"finances-api/repositories"
)

type PaymentAPIUsecases struct {
	Gateway repositories.GatewayRepository
	Db      repositories.FinancialRepository
	Pub     messaging.Publisher
}

func NewPaymentAPIUsecases(gatewayName string, db db.Database) *PaymentAPIUsecases {
	dbRepo := NewDbUsecase(db)
	gatewayRepo := NewGatewayUsecase(gatewayName)

	return &PaymentAPIUsecases{
		Gateway: gatewayRepo,
		Db:      dbRepo,
	}
}
