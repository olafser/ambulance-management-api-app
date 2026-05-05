package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/repository"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehicleServiceDeleteByID_DeletesUnfinishedDispatchesBeforeVehicle(t *testing.T) {
	var cleanupCallSign string
	var vehicleDeleted bool

	svc := service.NewVehicleService(
		vehicleRepoStub{
			getByIDFn: func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{VehicleID: vehicleID, CallSign: "A-101"}, nil
			},
			deleteByIDFn: func(ctx context.Context, vehicleID int64) error {
				vehicleDeleted = true
				return nil
			},
		},
		dispatchRepoStub{
			deleteUnfinishedFn: func(ctx context.Context, callSign string) error {
				cleanupCallSign = callSign
				return nil
			},
		},
	)

	if err := svc.DeleteByID(context.Background(), 7); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if cleanupCallSign != "A-101" {
		t.Fatalf("expected cleanup for call sign A-101, got %q", cleanupCallSign)
	}
	if !vehicleDeleted {
		t.Fatal("expected vehicle deletion after dispatch cleanup")
	}
}

func TestVehicleServiceDeleteByID_StopsWhenDispatchCleanupFails(t *testing.T) {
	var vehicleDeleted bool

	svc := service.NewVehicleService(
		vehicleRepoStub{
			getByIDFn: func(ctx context.Context, vehicleID int64) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{VehicleID: vehicleID, CallSign: "A-101"}, nil
			},
			deleteByIDFn: func(ctx context.Context, vehicleID int64) error {
				vehicleDeleted = true
				return nil
			},
		},
		dispatchRepoStub{
			deleteUnfinishedFn: func(ctx context.Context, callSign string) error {
				return repository.ErrDispatchConflict
			},
		},
	)

	err := svc.DeleteByID(context.Background(), 7)
	if !errors.Is(err, service.ErrConflict) {
		t.Fatalf("expected ErrConflict, got %v", err)
	}
	if vehicleDeleted {
		t.Fatal("expected vehicle not to be deleted when dispatch cleanup fails")
	}
}
