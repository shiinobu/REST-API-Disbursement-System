package services

import (
	"errors"
	"testing"

	"rest-api-disbursement-system/internal/models"
)

type stubDisbursementRepo struct {
	findAllPage   int
	findAllLimit  int
	findAllSearch string
	findAllStatus string

	findAllResult []models.Disbursement
	findAllTotal  int64

	findAllErr error
	exportErr  error
}

func (s *stubDisbursementRepo) Create(disbursement *models.Disbursement) error {
	return nil
}

func (s *stubDisbursementRepo) FindAll(page, limit int, search, status string) ([]models.Disbursement, int64, error) {
	s.findAllPage = page
	s.findAllLimit = limit
	s.findAllSearch = search
	s.findAllStatus = status

	return s.findAllResult, s.findAllTotal, s.findAllErr
}

func (s *stubDisbursementRepo) FindByID(id uint) (*models.Disbursement, error) {
	return nil, errors.New("not implemented")
}

func (s *stubDisbursementRepo) Update(disbursement *models.Disbursement) error {
	return errors.New("not implemented")
}

func (s *stubDisbursementRepo) Delete(id uint) error {
	return errors.New("not implemented")
}

func (s *stubDisbursementRepo) Export(status string) ([]models.Disbursement, error) {
	return nil, s.exportErr
}

func TestDisbursementServiceListFiltersByStatus(t *testing.T) {
	repo := &stubDisbursementRepo{
		findAllResult: []models.Disbursement{{ID: 1, Status: models.StatusPending}},
		findAllTotal:  1,
	}

	service := NewDisbursementService(repo)

	result, err := service.List(2, 20, "budi", string(models.StatusPending))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if repo.findAllPage != 2 || repo.findAllLimit != 20 || repo.findAllSearch != "budi" || repo.findAllStatus != string(models.StatusPending) {
		t.Fatalf("unexpected repository args: %+v", repo)
	}

	if result.Total != 1 || len(result.Data) != 1 {
		t.Fatalf("unexpected result: %+v", result)
	}
}

func TestDisbursementServiceListRejectsInvalidStatus(t *testing.T) {
	repo := &stubDisbursementRepo{}
	service := NewDisbursementService(repo)

	_, err := service.List(1, 10, "", "INVALID")
	if !errors.Is(err, ErrInvalidDisbursementStatus) {
		t.Fatalf("expected ErrInvalidDisbursementStatus, got %v", err)
	}
}

func TestDisbursementServiceExportRejectsInvalidStatus(t *testing.T) {
	repo := &stubDisbursementRepo{}
	service := NewDisbursementService(repo)

	_, err := service.Export("INVALID")
	if !errors.Is(err, ErrInvalidDisbursementStatus) {
		t.Fatalf("expected ErrInvalidDisbursementStatus, got %v", err)
	}
}
