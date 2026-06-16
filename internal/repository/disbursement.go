package repository

import (
	"gorm.io/gorm"

	"rest-api-disbursement-system/internal/models"
)

type DisbursementRepository interface {
	Create(disbursement *models.Disbursement) error
	FindAll() ([]models.Disbursement, error)
	FindByID(id uint) (*models.Disbursement, error)
	Update(disbursement *models.Disbursement) error
}

type disbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) DisbursementRepository {
	return &disbursementRepository{db: db}
}

func (r *disbursementRepository) Create(disbursement *models.Disbursement) error {
	return r.db.Create(disbursement).Error
}

func (r *disbursementRepository) FindAll() ([]models.Disbursement, error) {
	var disbursements []models.Disbursement
	err := r.db.
		Preload("Requester").
		Preload("ProcessedBy").
		Order("created_at DESC").
		Find(&disbursements).
		Error

	return disbursements, err
}

func (r *disbursementRepository) FindByID(id uint) (*models.Disbursement, error) {
	var disbursement models.Disbursement
	err := r.db.
		Preload("Requester").
		Preload("ProcessedBy").
		First(&disbursement, id).
		Error
	if err != nil {
		return nil, err
	}
	return &disbursement, nil
}

func (r *disbursementRepository) Update(disbursement *models.Disbursement) error {
	return r.db.Save(disbursement).Error
}
