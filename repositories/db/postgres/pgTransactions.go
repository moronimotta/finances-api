package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgTransactions interface {
	CreateTransaction(transaction *entities.Transactions) error
	GetTransactionByID(id string) (*entities.Transactions, error)
	UpdateTransaction(transaction *entities.Transactions) error
	DeleteTransaction(id string) error
	GetTransactionsByUserIDAndProductID(userID, productID string) ([]entities.Transactions, error)
	GetAllTransactions() ([]entities.Transactions, error)
	GetTransactionsByProductID(productID string) ([]entities.Transactions, error)
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
	if err := r.db.GetDB().Where("id = ?", id).First(transaction).Error; err != nil {
		return nil, err
	}
	return transaction, nil
}
func (r *pgTransactionsRepository) UpdateTransaction(transaction *entities.Transactions) error {
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
	if err := r.db.GetDB().Where("user_id = ? AND product_id = ?", userID, productID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

func (r *pgTransactionsRepository) GetAllTransactions() ([]entities.Transactions, error) {
	transactions := []entities.Transactions{}
	if err := r.db.GetDB().Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
func (r *pgTransactionsRepository) GetTransactionsByProductID(productID string) ([]entities.Transactions, error) {
	transactions := []entities.Transactions{}
	if err := r.db.GetDB().Where("product_id = ?", productID).Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}
