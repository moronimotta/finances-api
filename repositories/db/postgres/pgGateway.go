package postgresRepository

import (
	"finances-api/db"
	"finances-api/entities"
)

type PgGateway interface {
	CreateGateway(name string) error
	GetGatewayByID(id string) (*entities.Gateway, error)
	GetAllGateways() ([]entities.Gateway, error)
	UpdateGateway(gateway *entities.Gateway) error
	DeleteGateway(id string) error
}

type pgGatewayRepository struct {
	db db.Database
}

func NewPgGatewayRepository(db db.Database) PgGateway {
	return &pgGatewayRepository{
		db: db,
	}
}
func (r *pgGatewayRepository) CreateGateway(name string) error {
	gateway := &entities.Gateway{
		Name: name,
	}

	if err := r.db.GetDB().Create(gateway).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgGatewayRepository) GetGatewayByID(id string) (*entities.Gateway, error) {
	gateway := &entities.Gateway{}
	if err := r.db.GetDB().Where("id = ?", id).First(gateway).Error; err != nil {
		return nil, err
	}
	return gateway, nil
}
func (r *pgGatewayRepository) GetAllGateways() ([]entities.Gateway, error) {
	gateways := []entities.Gateway{}
	if err := r.db.GetDB().Find(&gateways).Error; err != nil {
		return nil, err
	}
	return gateways, nil
}
func (r *pgGatewayRepository) UpdateGateway(gateway *entities.Gateway) error {
	if err := r.db.GetDB().Model(&entities.Gateway{}).Where("id = ?", gateway.ID).Updates(gateway).Error; err != nil {
		return err
	}
	return nil
}
func (r *pgGatewayRepository) DeleteGateway(id string) error {
	if err := r.db.GetDB().Where("id = ?", id).Delete(&entities.Gateway{}).Error; err != nil {
		return err
	}
	return nil
}
