package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestDispatchServiceCreate_AssignsAvailableVehicleAndMarksItOnMission(t *testing.T) {
	var updatedStatus string

	svc := service.NewDispatchService(dispatchRepoStub{
		createFn: func(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
			dispatch.DispatchID = 21
			return dispatch, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.AVAILABLE)}, nil
		},
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			updatedStatus = status
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	got, err := svc.Create(context.Background(), validDispatchCreateRequest())
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if got.Id != 21 {
		t.Fatalf("expected dispatch id 21, got %d", got.Id)
	}
	if updatedStatus != string(model.ON_MISSION) {
		t.Fatalf("expected vehicle to be set to ON_MISSION, got %s", updatedStatus)
	}
}

func TestDispatchServiceCreate_RejectsUnavailableVehicle(t *testing.T) {
	var createCalled bool

	svc := service.NewDispatchService(dispatchRepoStub{
		createFn: func(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
			createCalled = true
			return entity.DispatchEntity{}, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
		},
	})

	_, err := svc.Create(context.Background(), validDispatchCreateRequest())
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if createCalled {
		t.Fatal("expected dispatch not to be created when vehicle is unavailable")
	}
}

func TestDispatchServiceUpdateByID_ReleasesVehicleWhenCompleted(t *testing.T) {
	var releasedStatus string

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.ON_ROUTE), AmbulanceCallSign: "A-101"}, nil
		},
		updateByIDFn: func(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
			dispatch.DispatchID = dispatchID
			return dispatch, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
		},
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			releasedStatus = status
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	_, err := svc.UpdateByID(context.Background(), 7, dispatchUpdateRequest(model.COMPLETED, "A-101"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if releasedStatus != string(model.AVAILABLE) {
		t.Fatalf("expected vehicle to be released to AVAILABLE, got %s", releasedStatus)
	}
}

func TestDispatchServiceUpdateStatusByID_ReleasesVehicleWhenCompleted(t *testing.T) {
	var releasedStatus string

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.ON_SCENE), AmbulanceCallSign: "A-101"}, nil
		},
		updateStatusFn: func(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: status, AmbulanceCallSign: "A-101"}, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
		},
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			releasedStatus = status
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	_, err := svc.UpdateStatusByID(context.Background(), 7, model.DispatchStatusUpdateRequest{Status: model.COMPLETED})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if releasedStatus != string(model.AVAILABLE) {
		t.Fatalf("expected vehicle to be released to AVAILABLE, got %s", releasedStatus)
	}
}

func TestDispatchServiceUpdateByID_ReopensCompletedDispatchAndMarksVehicleOnMission(t *testing.T) {
	var statuses []string

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.COMPLETED), AmbulanceCallSign: "A-101"}, nil
		},
		updateByIDFn: func(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
			dispatch.DispatchID = dispatchID
			return dispatch, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.AVAILABLE)}, nil
		},
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			statuses = append(statuses, status)
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	_, err := svc.UpdateByID(context.Background(), 7, dispatchUpdateRequest(model.ON_ROUTE, "A-101"))
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(statuses) != 1 || statuses[0] != string(model.ON_MISSION) {
		t.Fatalf("expected vehicle to be set to ON_MISSION, got %v", statuses)
	}
}

func TestDispatchServiceUpdateByID_RejectsReopeningCompletedDispatchWhenVehicleUnavailable(t *testing.T) {
	var updateCalled bool

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.COMPLETED), AmbulanceCallSign: "A-101"}, nil
		},
		updateByIDFn: func(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
			updateCalled = true
			return dispatch, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.ON_MISSION)}, nil
		},
	})

	_, err := svc.UpdateByID(context.Background(), 7, dispatchUpdateRequest(model.ON_ROUTE, "A-101"))
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if updateCalled {
		t.Fatal("expected dispatch update not to be called when vehicle is unavailable")
	}
}

func TestDispatchServiceUpdateStatusByID_RejectsReopeningCompletedDispatchWhenVehicleUnavailable(t *testing.T) {
	var updateCalled bool

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.COMPLETED), AmbulanceCallSign: "A-101"}, nil
		},
		updateStatusFn: func(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
			updateCalled = true
			return entity.DispatchEntity{DispatchID: dispatchID, Status: status, AmbulanceCallSign: "A-101"}, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.IN_SERVICE)}, nil
		},
	})

	_, err := svc.UpdateStatusByID(context.Background(), 7, model.DispatchStatusUpdateRequest{Status: model.ON_ROUTE})
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if updateCalled {
		t.Fatal("expected dispatch status update not to be called when vehicle is unavailable")
	}
}

func TestDispatchServiceUpdateStatusByID_ReopensCompletedDispatchAndMarksVehicleOnMission(t *testing.T) {
	var statuses []string

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.COMPLETED), AmbulanceCallSign: "A-101"}, nil
		},
		updateStatusFn: func(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: status, AmbulanceCallSign: "A-101"}, nil
		},
	}, vehicleRepoStub{
		getByCallSignFn: func(ctx context.Context, callSign string) (entity.VehicleEntity, error) {
			return entity.VehicleEntity{VehicleID: 9, CallSign: callSign, Status: string(model.AVAILABLE)}, nil
		},
		updateStatusFn: func(ctx context.Context, vehicleID int64, status string) (entity.VehicleEntity, error) {
			statuses = append(statuses, status)
			return entity.VehicleEntity{VehicleID: vehicleID, Status: status, CallSign: "A-101"}, nil
		},
	})

	_, err := svc.UpdateStatusByID(context.Background(), 7, model.DispatchStatusUpdateRequest{Status: model.ON_ROUTE})
	if err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if len(statuses) != 1 || statuses[0] != string(model.ON_MISSION) {
		t.Fatalf("expected vehicle to be set to ON_MISSION, got %v", statuses)
	}
}

func dispatchUpdateRequest(status model.DispatchStatus, callSign string) model.DispatchUpdateRequest {
	return model.DispatchUpdateRequest{
		IncidentNumber:      "INC-001",
		CallerName:          "John Doe",
		PatientName:         "Jane Doe",
		StreetAddress:       "Main Street 1",
		City:                "Bratislava",
		DispatchReason:      "Chest pain",
		Priority:            model.HIGH,
		Status:              status,
		AmbulanceCallSign:   callSign,
		DestinationHospital: "University Hospital",
		DispatcherName:      "Dispatcher One",
		CreatedAt:           time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
		Notes:               "ready",
	}
}

func validDispatchCreateRequest() model.DispatchCreateRequest {
	return model.DispatchCreateRequest{
		IncidentNumber:      "INC-001",
		CallerName:          "John Doe",
		PatientName:         "Jane Doe",
		StreetAddress:       "Main Street 1",
		City:                "Bratislava",
		DispatchReason:      "Chest pain",
		Priority:            model.HIGH,
		Status:              model.ACCEPTED,
		AmbulanceCallSign:   "A-101",
		DestinationHospital: "University Hospital",
		DispatcherName:      "Dispatcher One",
		CreatedAt:           time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
		Notes:               "ready",
	}
}
