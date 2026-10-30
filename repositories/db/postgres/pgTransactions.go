package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgTransactions interface {
	CreateTransaction(transaction *entities.Transactions) error
	GetTransactionByID(id string) (*entities.Transactions, error)
	GetTransactionByExternalID(externalID string) (*entities.Transactions, error)
	GetTransactionByAppointmentID(appointmentID string) (*entities.Transactions, error)
	GetTransactionByPaymentIntentID(paymentIntentID string) (*entities.Transactions, error)
	GetTransactionByCheckoutAttemptID(attemptID string) (*entities.Transactions, error)
	UpdateTransaction(transaction *entities.Transactions) error
	DeleteTransaction(id string) error
	GetTransactionsByUserIDAndProductID(userID, productID string) ([]entities.Transactions, error)
	GetAllTransactions() ([]entities.Transactions, error)
	GetTransactionsByProductID(productID string) ([]entities.Transactions, error)

	CreateTransactionItems(items []entities.TransactionItem) error
	UpdateTransactionItems(items []entities.TransactionItem) error
}

type pgTransactionsRepository struct {
	db db.Database
}

func NewPgTransactionsRepository(db db.Database) PgTransactions {
	return &pgTransactionsRepository{
		db: db,
	}
}

func (r *pgTransactionsRepository) CreateTransaction(transaction *entities.Transactions) error {
	if err := r.db.GetDB().Create(transaction).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgTransactionsRepository) GetTransactionByID(id string) (*entities.Transactions, error) {
	transaction := &entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("id = ?", id).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *pgTransactionsRepository) GetTransactionByExternalID(externalID string) (*entities.Transactions, error) {
	transaction := &entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("external_id = ?", externalID).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *pgTransactionsRepository) GetTransactionByAppointmentID(appointmentID string) (*entities.Transactions, error) {
	transaction := &entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("appointment_id = ?", appointmentID).Order("created_at DESC").First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *pgTransactionsRepository) GetTransactionByPaymentIntentID(paymentIntentID string) (*entities.Transactions, error) {
	transaction := &entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("stripe_payment_intent_id = ?", paymentIntentID).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *pgTransactionsRepository) GetTransactionByCheckoutAttemptID(attemptID string) (*entities.Transactions, error) {
	transaction := &entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("checkout_attempt_id = ?", attemptID).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}

func (r *pgTransactionsRepository) UpdateTransaction(transaction *entities.Transactions) error {
	if transaction.Type == "barber_appointment" {
		return r.db.GetDB().Model(&entities.Transactions{}).Where("id = ?", transaction.ID).Select(
			"UserExternalID", "AmountPayed", "AmountTotal", "AmountRefunded", "StripeFee", "NetAmount", "StripeFeeKnown",
			"Status", "PaymentStatus", "AppointmentID", "StripeSessionID", "StripePaymentIntentID", "CheckoutAttemptID", "CheckoutURL",
			"PaymentMethod", "PaymentDetails", "ReceiptURL", "ExternalID", "Currency", "Meta",
		).Updates(transaction).Error
	}
	if err := r.db.GetDB().Model(&entities.Transactions{}).Where("id = ?", transaction.ID).Updates(transaction).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgTransactionsRepository) DeleteTransaction(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.Transactions{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgTransactionsRepository) GetTransactionsByUserIDAndProductID(userID, productID string) ([]entities.Transactions, error) {
	transactions := []entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("user_id = ? AND product_id = ?", userID, productID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *pgTransactionsRepository) GetAllTransactions() ([]entities.Transactions, error) {
	transactions := []entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
func (r *pgTransactionsRepository) GetTransactionsByProductID(productID string) ([]entities.Transactions, error) {
	transactions := []entities.Transactions{}
	if err := r.db.GetDB().Preload("Items").Where("product_id = ?", productID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *pgTransactionsRepository) CreateTransactionItems(items []entities.TransactionItem) error {
	if len(items) == 0 {
		return nil
	}
	if err := r.db.GetDB().Create(&items).Error; err != nil {
		return err
	}
	return nil
}

func (r *pgTransactionsRepository) UpdateTransactionItems(items []entities.TransactionItem) error {
	if len(items) == 0 {
		return nil
	}
	tx := r.db.GetDB().Begin()
	if tx.Error != nil {
		return tx.Error
	}
	for _, it := range items {
		if err := tx.Model(&entities.TransactionItem{}).Where("id = ?", it.ID).Updates(it).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	return nil
}
