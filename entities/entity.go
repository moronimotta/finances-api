package entities

import (
	"finances-api/utils/meta"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Products struct {
	ID                       string `json:"id"`
	Name                     string `json:"name"`
	Description              string `json:"description"`
	Price                    int64  `json:"price"`
	Currency                 string `json:"currency"`
	GatewayProductExternalID string `json:"external_id"`
	GatewayPriceExternalID   string `json:"gateway_price_external_id"`
	GatewayID                string `json:"gateway_id"`
	CreatedAt                string `json:"created_at"`
	UpdatedAt                string `json:"updated_at"`
	DeletedAt                string `json:"deleted_at"`
	GatewayName              string `json:"gateway_name"`
}

type Checkout struct {
	CustomerID  string    `json:"customer_id"`
	PriceID     []string  `json:"price_id"`
	SessionID   string    `json:"session_id"`
	GatewayName string    `json:"gateway_name"`
	Meta        meta.Meta `json:"meta"`
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
	ID             string            `json:"id"`
	UserExternalID string            `json:"user_external_id"`
	Items          []TransactionItem `json:"items" gorm:"foreignKey:TransactionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	AmountPayed    int64             `json:"amount_paid"`
	AmountTotal    int64             `json:"amount_total"`
	AmountRefunded int64             `json:"amount_refunded"`
	Status         string            `json:"status"`          // paid, pending, failed, refunded
	Type           string            `json:"type"`            // product, subscription
	PaymentMethod  string            `json:"payment_method"`  // card, bank_transfer
	PaymentDetails string            `json:"payment_details"` // cardbrand...
	ReceiptURL     string            `json:"receipt_url"`
	ExternalID     string            `json:"external_id"` // e.g Stripe will be charge_id
	Currency       string            `json:"currency"`
	Meta           meta.Meta         `json:"meta"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
	// InvoiceID      string `json:"invoice_id"` // only for subscriptions
}

type TransactionItem struct {
	ID                string `json:"id"`
	TransactionID     string `json:"transaction_id" gorm:"index;not null"`
	ProductExternalID string `json:"product_external_id"`
	Quantity          int64  `json:"quantity"`
	UnitAmount        int64  `json:"unit_amount"`
	Currency          string `json:"currency"`

	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	DeletedAt string `json:"deleted_at"`
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

func (g *Gateway) BeforeCreate(tx *gorm.DB) error {
	g.ID = uuid.New().String()
	g.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	g.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

// Repeat for Products, UserProducts, Invoices, Transactions:
func (p *Products) BeforeCreate(tx *gorm.DB) error {
	p.ID = uuid.New().String()
	p.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (p *Products) BeforeUpdate(tx *gorm.DB) error {
	p.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (u *UserProducts) BeforeCreate(tx *gorm.DB) error {
	u.ID = uuid.New().String()
	u.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	u.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (i *Invoices) BeforeCreate(tx *gorm.DB) error {
	i.ID = uuid.New().String()
	i.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	i.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (t *Transactions) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New().String()
	t.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	t.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}

func (ti *TransactionItem) BeforeCreate(tx *gorm.DB) error {
	ti.ID = uuid.New().String()
	ti.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	ti.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	return nil
}
