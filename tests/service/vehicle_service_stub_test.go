package service_test

import (
	"context"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
)

type vehicleRepoStub struct {
	listFn         func(ctx context.Context, status, station string) ([]entity.VehicleEntity, error)
	createFn       func(ctx context.Context, vehicle entity.VehicleEntity) (entity.VehicleEntity, error)
	getByIDFn      func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error)
	updateByIDFn   func(ctx context.Context, vehicleID int64, vehicle entity.VehicleEntity) (entity.VehicleEntity, error)
	updateStatusFn func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error)
	deleteByIDFn   func(ctx context.Context, vehicleID int64) error
}

func (s vehicleRepoStub) List(ctx context.Context, status, station string) ([]entity.VehicleEntity, error) {
	if s.listFn != nil {
		return s.listFn(ctx, status, station)
	}
	return nil, nil
}

func (s vehicleRepoStub) Create(ctx context.Context, vehicle entity.VehicleEntity) (entity.VehicleEntity, error) {
	if s.createFn != nil {
		return s.createFn(ctx, vehicle)
	}
	return vehicle, nil
}

func (s vehicleRepoStub) GetByID(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, vehicleID)
	}
	return entity.VehicleEntity{}, nil
}

func (s vehicleRepoStub) UpdateByID(ctx context.Context, vehicleID int64, vehicle entity.VehicleEntity) (entity.VehicleEntity, error) {
	if s.updateByIDFn != nil {
		return s.updateByIDFn(ctx, vehicleID, vehicle)
	}
	return vehicle, nil
}

func (s vehicleRepoStub) UpdateStatusByID(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(ctx, vehicleID, status)
	}
	return entity.VehicleEntity{VehicleID: vehicleID, Status: status}, nil
}

func (s vehicleRepoStub) DeleteByID(ctx context.Context, vehicleID int64) error {
	if s.deleteByIDFn != nil {
		return s.deleteByIDFn(ctx, vehicleID)
	}
	return nil
}

func validCreateRequest() model.VehicleCreateRequest {
	return model.VehicleCreateRequest{
		CallSign:        "A-101",
		VehicleType:     "Type B",
		PlateNumber:     "BA123AA",
		Station:         "North",
		AssignedCrew:    "Crew 7",
		Status:          model.AVAILABLE,
		MileageKm:       15000,
		LastServiceDate: "2026-03-10",
		Notes:           "ready",
	}
}
