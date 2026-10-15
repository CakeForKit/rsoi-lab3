package service

import (
	"context"
	"testing"

	"github.com/CakeForKit/rsoi-lab2/lb-core/model"
	"github.com/CakeForKit/rsoi-lab2/ms-bonus/internal/repository"
	"github.com/google/uuid"
)

type bonusRepositoryStub struct {
	repository.BonusRepository
	bonuses []model.Privilege
	updated []model.Privilege
}

func (s *bonusRepositoryStub) GetAll(context.Context) ([]model.Privilege, error) {
	return s.bonuses, nil
}

func (s *bonusRepositoryStub) Update(_ context.Context, bonuses []model.Privilege) ([]model.Privilege, error) {
	s.updated = bonuses
	return bonuses, nil
}

type historyRepositoryStub struct {
	repository.BonusHistoryRepository
	created []model.PrivilegeHistory
}

func (s *historyRepositoryStub) Create(_ context.Context, history []model.PrivilegeHistory) ([]model.PrivilegeHistory, error) {
	s.created = history
	return history, nil
}

func TestApplyDebitUpdatesBalanceAndHistory(t *testing.T) {
	bonusRepository := &bonusRepositoryStub{bonuses: []model.Privilege{{
		BaseEntity: model.BaseEntity{ID: 1}, Username: "anna", Balance: 1500, Status: "SILVER",
	}}}
	historyRepository := &historyRepositoryStub{}
	service := &bonusService{bonusRepository: bonusRepository, historyRepository: historyRepository}

	result, err := service.Apply(context.Background(), "anna", model.PrivilegeOperationRequest{
		TicketUID: uuid.New(), Amount: 500, OperationType: model.PrivilegeOperationDebit,
	})
	if err != nil {
		t.Fatalf("Apply returned an error: %v", err)
	}
	if result.Balance != 1000 || result.Status != "SILVER" || len(historyRepository.created) != 1 || historyRepository.created[0].BalanceDiff != -500 {
		t.Fatalf("unexpected operation result: %#v, history: %#v", result, historyRepository.created)
	}
}
