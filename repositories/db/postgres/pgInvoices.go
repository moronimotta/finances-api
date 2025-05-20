package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgInvoices interface {
	CreateInvoice(invoice *entities.Invoices) error
	GetInvoiceByID(id string) (*entities.Invoices, error)
	UpdateInvoice(invoice *entities.Invoices) error
	DeleteInvoice(id string) error

	GetInvoicesByUserIDAndProductID(userID, productID string) ([]entities.Invoices, error)
	GetAllInvoices() ([]entities.Invoices, error)
	GetInvoicesByProductID(productID string) ([]entities.Invoices, error)
	GetInvoicesByCustomerID(customerID string) ([]entities.Invoices, error)
	GetInvoicesByPaymentStatus(paymentStatus string) ([]entities.Invoices, error)
	GetInvoicesByPaymentMethod(paymentMethod string) ([]entities.Invoices, error)
}

type pgInvoicesRepository struct {
	db db.Database
}

func NewPgInvoicesRepository(db db.Database) PgInvoices {
	return &pgInvoicesRepository{
		db: db,
	}
}

func (r *pgInvoicesRepository) CreateInvoice(invoice *entities.Invoices) error {
	if err := r.db.GetDB().Create(invoice).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgInvoicesRepository) GetInvoiceByID(id string) (*entities.Invoices, error) {
	invoice := &entities.Invoices{}
	if err := r.db.GetDB().Where("id = ?", id).First(invoice).Error; err != nil {
		return nil, err
	}
	return invoice, nil
}

func (r *pgInvoicesRepository) UpdateInvoice(invoice *entities.Invoices) error {
	if err := r.db.GetDB().Model(&entities.Invoices{}).Where("id = ?", invoice.ID).Updates(invoice).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgInvoicesRepository) DeleteInvoice(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.Invoices{}).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgInvoicesRepository) GetInvoicesByUserIDAndProductID(userID, productID string) ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Where("user_id = ? AND product_id = ?", userID, productID).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
func (r *pgInvoicesRepository) GetAllInvoices() ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
func (r *pgInvoicesRepository) GetInvoicesByProductID(productID string) ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Where("product_id = ?", productID).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
func (r *pgInvoicesRepository) GetInvoicesByCustomerID(customerID string) ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Where("customer_id = ?", customerID).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
func (r *pgInvoicesRepository) GetInvoicesByPaymentStatus(paymentStatus string) ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Where("payment_status = ?", paymentStatus).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
func (r *pgInvoicesRepository) GetInvoicesByPaymentMethod(paymentMethod string) ([]entities.Invoices, error) {
	invoices := []entities.Invoices{}
	if err := r.db.GetDB().Where("payment_method = ?", paymentMethod).Find(&invoices).Error; err != nil {
		return nil, err
	}
	return invoices, nil
}
