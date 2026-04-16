package service_test

import (
	"context"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehicleServiceUpdateStatusByID_UpdatesStatus(t *testing.T) {
	svc := service.NewVehicleService(vehicleRepoStub{
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	got, err := svc.UpdateStatusByID(context.Background(), 7, model.VehicleStatusUpdateRequest{Status: model.IN_SERVICE})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Id != 7 || got.Status != model.IN_SERVICE {
		t.Fatalf("unexpected vehicle after status update: %+v", got)
	}
}
