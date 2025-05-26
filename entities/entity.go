package entities

import (
	"finances-api/utils/meta"
	"time"

	"github.com/google/uuid"
)

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

func (p *Products) BeforeCreate() {
	p.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	p.CreatedAt = now
	p.UpdatedAt = now
}

func (u *UserProducts) BeforeCreate() {
	u.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	u.CreatedAt = now
	u.UpdatedAt = now
}

func (i *Invoices) BeforeCreate() {
	i.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	i.CreatedAt = now
	i.UpdatedAt = now
}

func (t *Transactions) BeforeCreate() {
	t.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	t.CreatedAt = now
	t.UpdatedAt = now
}

func (g *Gateway) BeforeCreate() {
	g.ID = uuid.New().String()
	now := time.Now().UTC().Format(time.RFC3339)
	g.CreatedAt = now
	g.UpdatedAt = now
}
