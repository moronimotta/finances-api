package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgUserProducts interface {
	CreateUserProduct(userID, productID string) error
	GetUserProductByID(id string) (*entities.UserProducts, error)
	UpdateUserProduct(userProduct *entities.UserProducts) error
	UpdateUserProductStatus(userProduct *entities.UserProducts) error
	DeleteUserProduct(id string) error
	GetUserProductsByUserID(userID string) ([]entities.UserProducts, error)
}

type pgUserProductsRepository struct {
	db db.Database
}

func NewPgUserProductsRepository(db db.Database) PgUserProducts {
	return &pgUserProductsRepository{
		db: db,
	}
}
func (r *pgUserProductsRepository) CreateUserProduct(userID, productID string) error {
	userProduct := &entities.UserProducts{
		UserID:    userID,
		ProductID: productID,
	}
	if err := r.db.GetDB().Create(userProduct).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgUserProductsRepository) GetUserProductByID(id string) (*entities.UserProducts, error) {
	userProduct := &entities.UserProducts{}
	if err := r.db.GetDB().Where("id = ?", id).First(userProduct).Error; err != nil {
		return nil, err
	}
	return userProduct, nil
}
func (r *pgUserProductsRepository) UpdateUserProduct(userProduct *entities.UserProducts) error {
	if err := r.db.GetDB().Model(&entities.UserProducts{}).Where("id = ?", userProduct.ID).Updates(userProduct).Error; err != nil {
		return err
	}
	return nil
}

func (r *pgUserProductsRepository) UpdateUserProductStatus(userProduct *entities.UserProducts) error {
	if err := r.db.GetDB().Model(&entities.UserProducts{}).Where("id = ?", userProduct.ID).Updates(userProduct).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgUserProductsRepository) DeleteUserProduct(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.UserProducts{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *pgUserProductsRepository) GetUserProductsByUserID(userID string) ([]entities.UserProducts, error) {
	userProducts := []entities.UserProducts{}
	if err := r.db.GetDB().Where("user_id = ?", userID).Find(&userProducts).Error; err != nil {
		return nil, err
	}
	return userProducts, nil
}
