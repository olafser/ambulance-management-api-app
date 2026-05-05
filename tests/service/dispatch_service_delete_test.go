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
	return entity.DispatchEntity{DispatchID: dispatchID, Status: status}, nil
}

func (s dispatchRepoStub) DeleteByID(ctx context.Context, dispatchID int64) error {
	if s.deleteByIDFn != nil {
		return s.deleteByIDFn(ctx, dispatchID)
	}
	return nil
}

func TestDispatchServiceDeleteByID_AllowsCompletedDispatch(t *testing.T) {
	var deleteCalled bool

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.COMPLETED)}, nil
		},
		deleteByIDFn: func(ctx context.Context, dispatchID int64) error {
			deleteCalled = true
			return nil
		},
	})

	if err := svc.DeleteByID(context.Background(), 12); err != nil {
		t.Fatalf("expected nil error, got %v", err)
	}
	if !deleteCalled {
		t.Fatal("expected delete to be called for completed dispatch")
	}
}

func TestDispatchServiceDeleteByID_RejectsNonCompletedDispatch(t *testing.T) {
	var deleteCalled bool

	svc := service.NewDispatchService(dispatchRepoStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (entity.DispatchEntity, error) {
			return entity.DispatchEntity{DispatchID: dispatchID, Status: string(model.ON_ROUTE)}, nil
		},
		deleteByIDFn: func(ctx context.Context, dispatchID int64) error {
			deleteCalled = true
			return nil
		},
	})

	err := svc.DeleteByID(context.Background(), 12)
	if !errors.Is(err, service.ErrBadRequest) {
		t.Fatalf("expected ErrBadRequest, got %v", err)
	}
	if deleteCalled {
		t.Fatal("expected delete not to be called for non-completed dispatch")
	}
}
