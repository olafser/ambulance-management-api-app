package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/repository"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestDispatchServiceCreate_MarksVehicleOnMissionWhenAvailable(t *testing.T) {
	var capturedVehicleID int64
	var capturedStatus string

	svc := service.NewDispatchService(
		dispatchRepoStub{
			createFn: func(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
				return entity.DispatchEntity{DispatchID: 101, AmbulanceCallSign: dispatch.AmbulanceCallSign, Status: dispatch.Status}, nil
			},
		},
		vehicleRepoStub{
			getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{VehicleID: 77, CallSign: callSign, Status: string(model.AVAILABLE)}, nil
			},
			updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
				capturedVehicleID = vehicleID
				capturedStatus = status
				return entity.VehicleEntity{VehicleID: vehicleID, Status: status}, nil
			},
		},
	)

	got, err := svc.Create(context.Background(), validDispatchCreateRequest())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Id != 101 {
		t.Fatalf("expected dispatch id 101, got %d", got.Id)
	}
	if capturedVehicleID != 77 || capturedStatus != string(model.ON_MISSION) {
		t.Fatalf("unexpected vehicle update: id=%d status=%s", capturedVehicleID, capturedStatus)
	}
}

func TestDispatchServiceUpdateStatusByID_ReleasesVehicleWhenCompleted(t *testing.T) {
	var capturedStatus string

	svc := service.NewDispatchService(
		dispatchRepoStub{
			updateStatusFn: func(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
				return entity.DispatchEntity{
					DispatchID:        dispatchID,
					AmbulanceCallSign: "AMB-101",
					Status:            status,
					UpdatedAt:         updatedAt,
				}, nil
			},
		},
		vehicleRepoStub{
			getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{VehicleID: 77, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
			},
			updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
				capturedStatus = status
				return entity.VehicleEntity{VehicleID: vehicleID, Status: status}, nil
			},
		},
	)

	got, err := svc.UpdateStatusByID(context.Background(), 7, model.DispatchStatusUpdateRequest{Status: model.COMPLETED})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Status != model.COMPLETED {
		t.Fatalf("expected completed dispatch, got %+v", got)
	}
	if capturedStatus != string(model.AVAILABLE) {
		t.Fatalf("expected vehicle to become AVAILABLE, got %s", capturedStatus)
	}
}

func TestDispatchServiceUpdateByID_ReleasesVehicleWhenCompleted(t *testing.T) {
	var capturedStatus string

	svc := service.NewDispatchService(
		dispatchRepoStub{
			updateByIDFn: func(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
				return entity.DispatchEntity{
					DispatchID:        dispatchID,
					AmbulanceCallSign: dispatch.AmbulanceCallSign,
					Status:            dispatch.Status,
					IncidentNumber:    dispatch.IncidentNumber,
					CallerName:        dispatch.CallerName,
					PatientName:       dispatch.PatientName,
					StreetAddress:     dispatch.StreetAddress,
					City:              dispatch.City,
					DispatchReason:    dispatch.DispatchReason,
					Priority:          dispatch.Priority,
					DispatcherName:    dispatch.DispatcherName,
					CreatedAt:         dispatch.CreatedAt,
					UpdatedAt:         dispatch.UpdatedAt,
					Notes:             dispatch.Notes,
				}, nil
			},
		},
		vehicleRepoStub{
			getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{VehicleID: 77, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
			},
			updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
				capturedStatus = status
				return entity.VehicleEntity{VehicleID: vehicleID, Status: status}, nil
			},
		},
	)

	req := model.DispatchUpdateRequest{
		IncidentNumber:      "DIS-2026-0002",
		CallerName:          "Caller",
		PatientName:         "Patient",
		StreetAddress:       "Main street 2",
		City:                "Trnava",
		DispatchReason:      "Injury",
		Priority:            model.MEDIUM,
		Status:              model.COMPLETED,
		AmbulanceCallSign:   "AMB-101",
		DestinationHospital: "City hospital",
		DispatcherName:      "Dispatcher",
		CreatedAt:           time.Date(2026, 4, 20, 9, 0, 0, 0, time.UTC),
		Notes:               "done",
	}

	got, err := svc.UpdateByID(context.Background(), 7, req)
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Status != model.COMPLETED {
		t.Fatalf("expected completed dispatch, got %+v", got)
	}
	if capturedStatus != string(model.AVAILABLE) {
		t.Fatalf("expected vehicle to become AVAILABLE, got %s", capturedStatus)
	}
}

func TestDispatchServiceCreate_MapsVehicleNotFound(t *testing.T) {
	svc := service.NewDispatchService(
		dispatchRepoStub{},
		vehicleRepoStub{
			getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
				return entity.VehicleEntity{}, repository.ErrVehicleNotFound
			},
		},
	)

	_, err := svc.Create(context.Background(), validDispatchCreateRequest())
	if !errors.Is(err, service.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}
