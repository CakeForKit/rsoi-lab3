package service

import (
	"context"
	"testing"
	"time"

	"github.com/CakeForKit/rsoi-lab2/lb-core/model"
	"github.com/CakeForKit/rsoi-lab2/ms-flight/internal/repository"
)

type flightRepositoryStub struct {
	repository.FlightRepository
	flights []model.Flight
}

func (s flightRepositoryStub) GetAll(context.Context) ([]model.Flight, error) {
	return s.flights, nil
}

type airportRepositoryStub struct {
	repository.AirportRepository
	airports []model.Airport
}

func (s airportRepositoryStub) GetAll(context.Context) ([]model.Airport, error) {
	return s.airports, nil
}

func TestGetBuildsFlightResponse(t *testing.T) {
	date := time.Date(2021, 10, 8, 20, 0, 0, 0, time.UTC)
	service := &flightService{
		flightRepository: flightRepositoryStub{flights: []model.Flight{{
			FlightNumber: "AFL031", DateTime: date, FromAirportID: 1, ToAirportID: 2, Price: 1500,
		}}},
		airportRepository: airportRepositoryStub{airports: []model.Airport{
			{BaseEntity: model.BaseEntity{ID: 1}, City: "Санкт-Петербург", Name: "Пулково"},
			{BaseEntity: model.BaseEntity{ID: 2}, City: "Москва", Name: "Шереметьево"},
		}},
	}

	flight, err := service.Get(context.Background(), "AFL031")
	if err != nil {
		t.Fatalf("Get returned an error: %v", err)
	}
	if flight.FromAirport != "Санкт-Петербург Пулково" || flight.ToAirport != "Москва Шереметьево" || flight.Date != "2021-10-08 20:00" {
		t.Fatalf("unexpected flight response: %#v", flight)
	}
}
