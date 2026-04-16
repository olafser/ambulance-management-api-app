package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehicleServiceCreate_ValidatesBusinessRules(t *testing.T) {
	tests := []struct {
		name string
		mut  func(*model.VehicleCreateRequest)
	}{
		{
			name: "required fields",
			mut: func(req *model.VehicleCreateRequest) {
				req.CallSign = ""
			},
		},
		{
			name: "invalid status",
			mut: func(req *model.VehicleCreateRequest) {
				req.Status = model.VehicleStatus("BROKEN")
			},
		},
		{
			name: "negative mileage",
			mut: func(req *model.VehicleCreateRequest) {
				req.MileageKm = -1
			},
		},
		{
			name: "invalid service date format",
			mut: func(req *model.VehicleCreateRequest) {
				req.LastServiceDate = "10-03-2026"
			},
		},
	}

	svc := service.NewVehicleService(vehicleRepoStub{})

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := validCreateRequest()
			tc.mut(&req)

			_, err := svc.Create(context.Background(), req)
			if !errors.Is(err, service.ErrBadRequest) {
				t.Fatalf("expected ErrBadRequest, got %v", err)
			}
		})
	}
}

func TestVehicleServiceDeleteByID_RejectsNonPositiveID(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{})

	err := svc.DeleteByID(context.Background(), 0)
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
}
