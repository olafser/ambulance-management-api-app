package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehicleServiceUpdateStatusByID_RejectsChangeWhileOnMission(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{
		getByIDFn: func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: vehicleID, Status: string(model.ON_MISSION), CallSign: "A-101"}, nil
		},
	})

	_, err := svc.UpdateStatusByID(context.Background(), 7, model.VehicleStatusUpdateRequest{Status: model.AVAILABLE})
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}

func TestVehicleServiceUpdateByID_RejectsStatusChangeWhileOnMission(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{
		getByIDFn: func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: vehicleID, Status: string(model.ON_MISSION), CallSign: "A-101"}, nil
		},
	})

	req := model.VehicleUpdateRequest{
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

	_, err := svc.UpdateByID(context.Background(), 7, req)
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
