package repositories

import (
	"finances-api/entities"
	"finances-api/utils/meta"
)

type GatewayRepository interface {
	CreateProduct(name, description string, localProduct entities.Products) (string, error)
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
	UpdateProduct(productID, name, description string, meta meta.Meta) error
	CreateCheckoutSession(productID, priceID string) (string, error)
	DeactivateProduct(productID string) error
	CreateCustomer(name, email, localUserID string) (string, error)
	UpdateCustomer(customerID, name, email string) error
}

type FinancialRepository interface {
	// Products
	CreateProduct(product *entities.Products) error
	GetProductByID(id string) (*entities.Products, error)
	GetAllProducts() ([]entities.Products, error)
	GetProductByExternalID(externalID string) (*entities.Products, error)
	UpdateProduct(product *entities.Products) error
	DeleteProduct(id string) error

	// UserProducts
	CreateUserProduct(userID, productID string) error
	GetUserProductByID(id string) (*entities.UserProducts, error)
	UpdateUserProduct(userProduct *entities.UserProducts) error
	UpdateUserProductStatus(userProduct *entities.UserProducts) error
	DeleteUserProduct(id string) error
	GetUserProductsByUserID(userID string) ([]entities.UserProducts, error)

	// Gateway
	CreateGateway(gateway entities.Gateway) error
	GetGatewayByID(id string) (*entities.Gateway, error)
	GetAllGateways() ([]entities.Gateway, error)
	UpdateGateway(gateway *entities.Gateway) error
	DeleteGateway(id string) error

	// Transactions
	CreateTransaction(transaction *entities.Transactions) error
	GetTransactionByID(id string) (*entities.Transactions, error)
	UpdateTransaction(transaction *entities.Transactions) error
	DeleteTransaction(id string) error
	GetTransactionsByUserIDAndProductID(userID, productID string) ([]entities.Transactions, error)
	GetAllTransactions() ([]entities.Transactions, error)
	GetTransactionsByProductID(productID string) ([]entities.Transactions, error)

	// Invoices
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
