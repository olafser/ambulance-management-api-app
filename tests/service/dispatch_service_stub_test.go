package service_test

import (
	"context"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
)

type dispatchRepoStub struct {
	listFn         func(ctx context.Context, status, city string) ([]entity.DispatchEntity, error)
	createFn       func(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error)
	getByIDFn      func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error)
	updateByIDFn   func(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error)
	updateStatusFn func(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error)
	deleteByIDFn   func(ctx context.Context, dispatchID int64) error
}

func (s dispatchRepoStub) List(ctx context.Context, status, city string) ([]entity.DispatchEntity, error) {
	if s.listFn != nil {
		return s.listFn(ctx, status, city)
	}
	return nil, nil
}

func (s dispatchRepoStub) Create(ctx context.Context, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
	if s.createFn != nil {
		return s.createFn(ctx, dispatch)
	}
	return dispatch, nil
}

func (s dispatchRepoStub) GetByID(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, dispatchID)
	}
	return entity.DispatchEntity{}, nil
}

func (s dispatchRepoStub) UpdateByID(ctx context.Context, dispatchID int64, dispatch entity.DispatchEntity) (entity.DispatchEntity, error) {
	if s.updateByIDFn != nil {
		return s.updateByIDFn(ctx, dispatchID, dispatch)
	}
	return dispatch, nil
}

func (s dispatchRepoStub) UpdateStatusByID(ctx context.Context, dispatchID int64, status string, updatedAt time.Time) (entity.DispatchEntity, error) {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(ctx, dispatchID, status, updatedAt)
	}
	return entity.DispatchEntity{DispatchID: dispatchID, Status: status, UpdatedAt: updatedAt}, nil
}

func (s dispatchRepoStub) DeleteByID(ctx context.Context, dispatchID int64) error {
	if s.deleteByIDFn != nil {
		return s.deleteByIDFn(ctx, dispatchID)
	}
	return nil
}

func validDispatchCreateRequest() model.DispatchCreateRequest {
	return model.DispatchCreateRequest{
		IncidentNumber:      "DIS-2026-0001",
		CallerName:          "Caller",
		PatientName:         "Patient",
		StreetAddress:       "Main street 1",
		City:                "Bratislava",
		DispatchReason:      "Chest pain",
		Priority:            model.HIGH,
		Status:              model.ACCEPTED,
		AmbulanceCallSign:   "AMB-101",
		DestinationHospital: "City hospital",
		DispatcherName:      "Dispatcher",
		CreatedAt:           time.Date(2026, 4, 20, 8, 0, 0, 0, time.UTC),
		Notes:               "urgent",
	}
}

var _ = dispatchRepoStub{}
var _ = validDispatchCreateRequest
