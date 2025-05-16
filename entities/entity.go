package entities

import "finances-api/utils/meta"

// Product
// Invoice
// Customer (save in the user-api service)

type Product struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Price       int64  `json:"price"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
	ExternalID  string `json:"external_id"`
}

type Invoice struct {
	ID            string    `json:"id"`
	CustomerID    string    `json:"customer_id"`
	Amount        int64     `json:"amount"`
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
	DeletedAt     string    `json:"deleted_at"`
	ExternalID    string    `json:"external_id"`
	PaymentStatus string    `json:"payment_status"`
	PaymentMethod string    `json:"payment_method"`
	Meta          meta.Meta `json:"meta"`
}
