package repositories

type FinancialRepository interface {
	CreateProduct(name, description string) (string, error)
	CreatePrice(productID string, unitAmount int64, currency string) (string, error)
}
