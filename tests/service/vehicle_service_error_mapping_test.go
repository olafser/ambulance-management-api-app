package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/repository"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehicleServiceCreate_MapsConflictFromRepository(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{
		createFn: func(ctx context.Context, vehicle entity.VehicleEntity) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{}, repository.ErrVehicleConflict
		},
	})

	_, err := svc.Create(context.Background(), validCreateRequest())
	if !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
}

func TestVehicleServiceGetByID_MapsNotFoundFromRepository(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{
		getByIDFn: func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{}, repository.ErrVehicleNotFound
		},
	})

	_, err := svc.GetByID(context.Background(), 42)
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
