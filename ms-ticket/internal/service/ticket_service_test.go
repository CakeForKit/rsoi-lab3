package service

import (
	"context"
	"testing"

	"github.com/CakeForKit/rsoi-lab2/lb-core/model"
	"github.com/CakeForKit/rsoi-lab2/ms-ticket/internal/repository"
	"github.com/google/uuid"
)

type ticketRepositoryStub struct {
	repository.TicketRepository
	tickets []model.Ticket
}

func (s ticketRepositoryStub) GetAll(context.Context) ([]model.Ticket, error) {
	return s.tickets, nil
}

func TestGetAllReturnsOnlyUserTickets(t *testing.T) {
	service := &ticketService{repository: ticketRepositoryStub{tickets: []model.Ticket{
		{TicketUID: uuid.New(), Username: "anna", FlightNumber: "AFL031"},
		{TicketUID: uuid.New(), Username: "ivan", FlightNumber: "AFL032"},
	}}}

	tickets, err := service.GetAll(context.Background(), "anna")
	if err != nil {
		t.Fatalf("GetAll returned an error: %v", err)
	}
	if len(tickets) != 1 || tickets[0].Username != "anna" {
		t.Fatalf("unexpected tickets: %#v", tickets)
	}
}
