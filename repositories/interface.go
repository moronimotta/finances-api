package repositories

import "finances-api/entities"

type GatewayRepository interface {
	CreateProduct(name, description string) (string, error)
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
	CreateCheckoutSession(productID, priceID, successURL, cancelURL string) (string, error)
}

type FinancialRepository interface {
	// Products
	CreateProduct(name, description, externalID, gatewayName string, price int64) error
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
	CreateGateway(name string) error
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
