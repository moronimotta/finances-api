package entities

import "finances-api/utils/meta"

type Products struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	Currency    string `json:"currency"`
	ExternalID  string `json:"external_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
	GatewayName string `json:"gateway_name"`
}

type Checkout struct {
	CustomerID  string `json:"customer_id"`
	PriceID     string `json:"price_id"`
	SuccessURL  string `json:"success_url"`
	CancelURL   string `json:"cancel_url"`
	SessionID   string `json:"session_id"`
	GatewayName string `json:"gateway_name"`
}

type UserProducts struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	ProductID string `json:"product_id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	Status    string `json:"status"`
}

type Invoices struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	DeletedAt     string    `json:"deleted_at"`
	ExternalID    string    `json:"external_id"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod string    `json:"payment_method"`
	Currency      string    `json:"currency"`
	Meta          meta.Meta `json:"meta"`
}

type Transactions struct {
	ID          string `json:"id"`
	AmountPayed int64  `json:"amount_paid"`
	AmountTotal int64  `json:"amount_total"`
	InvoiceID   string `json:"invoice_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
	ExternalID  string `json:"external_id"`
}

type Gateway struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	SecretKey    string `json:"secret_key"`
	ClientSecret string `json:"client_secret"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
	DeletedAt    string `json:"deleted_at"`
}
