package services

import (
	"errors"

	"gorm.io/gorm"

	"rest-api-disbursement-system/internal/models"
	"rest-api-disbursement-system/internal/repository"
)

var (
	ErrDisbursementNotFound         = errors.New("disbursement tidak ditemukan")
	ErrDisbursementAlreadyProcessed = errors.New("disbursement sudah diproses")
)

type DisbursementService interface {
	Create(input CreateDisbursementInput) (*models.Disbursement, error)
	List() ([]models.Disbursement, error)
	Detail(id uint) (*models.Disbursement, error)
	Approve(id uint, userID uint) (*models.Disbursement, error)
	Reject(id uint, userID uint, reason string) (*models.Disbursement, error)
}

type CreateDisbursementInput struct {
	RequesterID     uint
	BeneficiaryName string
	BankName        string
	AccountNumber   string
	Amount          float64
	Description     string
}

type disbursementService struct {
	disbursements repository.DisbursementRepository
}

func NewDisbursementService(disbursements repository.DisbursementRepository) DisbursementService {
	return &disbursementService{disbursements: disbursements}
}

func (s *disbursementService) Create(input CreateDisbursementInput) (*models.Disbursement, error) {
	disbursement := &models.Disbursement{
		RequesterID:     input.RequesterID,
		BeneficiaryName: input.BeneficiaryName,
		BankName:        input.BankName,
		AccountNumber:   input.AccountNumber,
		Amount:          input.Amount,
		Description:     input.Description,
		Status:          models.StatusPending,
	}

	if err := s.disbursements.Create(disbursement); err != nil {
		return nil, err
	}

	return s.Detail(disbursement.ID)
}

func (s *disbursementService) List() ([]models.Disbursement, error) {
	return s.disbursements.FindAll()
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

func (s *disbursementService) Approve(id uint, userID uint) (*models.Disbursement, error) {
	disbursement, err := s.Detail(id)
	if err != nil {
		return nil, err
	}
	if disbursement.Status != models.StatusPending {
		return nil, ErrDisbursementAlreadyProcessed
	}

	disbursement.Status = models.StatusApproved
	disbursement.ProcessedByID = &userID
	disbursement.RejectionReason = nil

	if err := s.disbursements.Update(disbursement); err != nil {
		return nil, err
	}

	return s.Detail(id)
}

func (s *disbursementService) Reject(id uint, userID uint, reason string) (*models.Disbursement, error) {
	disbursement, err := s.Detail(id)
	if err != nil {
		return nil, err
	}
	if disbursement.Status != models.StatusPending {
		return nil, ErrDisbursementAlreadyProcessed
	}

	disbursement.Status = models.StatusRejected
	disbursement.ProcessedByID = &userID
	disbursement.RejectionReason = &reason

	if err := s.disbursements.Update(disbursement); err != nil {
		return nil, err
	}

	return s.Detail(id)
}
