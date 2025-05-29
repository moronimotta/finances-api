package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgProducts interface {
	CreateProduct(product *entities.Products) error
	GetProductByID(id string) (*entities.Products, error)
	GetAllProducts() ([]entities.Products, error)
	GetProductByExternalID(externalID string) (*entities.Products, error)
	UpdateProduct(product *entities.Products) error
	DeleteProduct(id string) error
}

type pgProductsRepository struct {
	db db.Database
}

func NewPgProductsRepository(db db.Database) PgProducts {
	return &pgProductsRepository{
		db: db,
	}
}

func (r *pgProductsRepository) CreateProduct(product *entities.Products) error {

	if err := r.db.GetDB().Create(&product).Error; err != nil {
		return err
	}
	return nil
}

func (r *pgProductsRepository) GetProductByID(id string) (*entities.Products, error) {
	product := &entities.Products{}
	if err := r.db.GetDB().Where("id = ?", id).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *pgProductsRepository) GetProductByExternalID(externalID string) (*entities.Products, error) {
	product := &entities.Products{}
	if err := r.db.GetDB().Where("external_id = ?", externalID).First(product).Error; err != nil {
		return nil, err
	}
	return product, nil
}

func (r *pgProductsRepository) GetAllProducts() ([]entities.Products, error) {
	products := []entities.Products{}
	if err := r.db.GetDB().Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *pgProductsRepository) UpdateProduct(product *entities.Products) error {
	if err := r.db.GetDB().Model(&entities.Products{}).Where("id = ?", product.ID).Updates(product).Error; err != nil {
		return err
	}

	if err := r.db.GetDB().First(product, "id = ?", product.ID).Error; err != nil {
		return err
	}

	return nil
}

func (r *pgProductsRepository) DeleteProduct(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.Products{}).Error; err != nil {
		return err
	}
	return nil
}
