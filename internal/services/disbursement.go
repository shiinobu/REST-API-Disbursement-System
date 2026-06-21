package services

import (
	"errors"
	"strings"

	"gorm.io/gorm"

	"math"
	"rest-api-disbursement-system/internal/models"
	"rest-api-disbursement-system/internal/repository"
)

var (
	ErrDisbursementNotFound         = errors.New("disbursement tidak ditemukan")
	ErrDisbursementAlreadyProcessed = errors.New("disbursement sudah diproses")
	ErrDisbursementAmountTooLow     = errors.New("jumlah disbursement harus lebih dari 10.000")
	ErrInvalidDisbursementStatus    = errors.New("status disbursement tidak valid")
	ErrForbidden                    = errors.New("forbidden")
	ErrDisbursementCannotBeDeleted  = errors.New("disbursement tidak dapat dihapus")
)

type DisbursementService interface {
	Create(input CreateDisbursementInput) (*models.Disbursement, error)
	List(page, limit int, search, status string) (*DisbursementListResult, error)
	Detail(id uint) (*models.Disbursement, error)
	UpdateStatus(id uint, userID uint, role string, status string, note string) (*models.Disbursement, error)
	Delete(id uint) error
	Export(status string) ([]models.Disbursement, error)
}

type CreateDisbursementInput struct {
	RequesterID   uint
	RecipientName string
	BankCode      string
	AccountNumber string
	Amount        float64
	Note          string
}

type DisbursementListResult struct {
	Data       []models.Disbursement
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type disbursementService struct {
	disbursements repository.DisbursementRepository
}

func NewDisbursementService(disbursements repository.DisbursementRepository) DisbursementService {
	return &disbursementService{disbursements: disbursements}
}

func (s *disbursementService) Create(input CreateDisbursementInput) (*models.Disbursement, error) {
	if input.Amount < 10000 {
		return nil, ErrDisbursementAmountTooLow
	}

	var adminFee float64
	if input.Amount >= 5000000 {
		adminFee = 5000
	} else {
		adminFee = 2500
	}

	disbursement := &models.Disbursement{
		RequesterID:   input.RequesterID,
		RecipientName: input.RecipientName,
		BankCode:      input.BankCode,
		AccountNumber: input.AccountNumber,
		Amount:        input.Amount,
		AdminFee:      adminFee,
		Note:          input.Note,
		Status:        models.StatusPending,
	}

	if err := s.disbursements.Create(disbursement); err != nil {
		return nil, err
	}

	return s.Detail(disbursement.ID)
}

func validateDisbursementStatus(status string) error {
	switch status {
	case "":
	case string(models.StatusPending):
	case string(models.StatusApproved):
	case string(models.StatusRejected):
	default:
		return ErrInvalidDisbursementStatus
	}

	return nil
}

func (s *disbursementService) List(page, limit int, search, status string) (*DisbursementListResult, error) {
	if err := validateDisbursementStatus(status); err != nil {
		return nil, err
	}

	data, total, err := s.disbursements.FindAll(page, limit, search, status)

	if err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &DisbursementListResult{
		Data:       data,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *disbursementService) Detail(id uint) (*models.Disbursement, error) {
	disbursement, err := s.disbursements.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDisbursementNotFound
		}
		return nil, err
	}

	return disbursement, nil
}

func (s *disbursementService) UpdateStatus(id, userID uint, role, status, note string) (*models.Disbursement, error) {
	if strings.ToUpper(role) != "ADMIN" {
		return nil, ErrForbidden
	}

	disbursement, err := s.Detail(id)
	if err != nil {
		return nil, err
	}

	if disbursement.Status != models.StatusPending {
		return nil, ErrDisbursementAlreadyProcessed
	}

	switch status {
	case string(models.StatusApproved):
		disbursement.Status = models.StatusApproved
		disbursement.ProcessedByID = &userID
		disbursement.RejectionReason = nil
	case string(models.StatusRejected):
		disbursement.Status = models.StatusRejected
		disbursement.ProcessedByID = &userID
		disbursement.RejectionReason = &note
	default:
		return nil, ErrInvalidDisbursementStatus
	}

	if err := s.disbursements.Update(disbursement); err != nil {
		return nil, err
	}

	return s.Detail(id)
}

func (s *disbursementService) Delete(id uint) error {
	disbursement, err := s.Detail(id)
	if err != nil {
		return err
	}

	if disbursement.Status != models.StatusPending {
		return ErrDisbursementCannotBeDeleted
	}

	return s.disbursements.Delete(id)
}

func (s *disbursementService) Export(status string) ([]models.Disbursement, error) {
	if err := validateDisbursementStatus(status); err != nil {
		return nil, err
	}

	return s.disbursements.Export(status)
}
