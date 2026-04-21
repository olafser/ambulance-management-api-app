package handler_test

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/olafser/ambulance-management-api-app/internal/handler"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

type dispatchServiceStub struct {
	listFn         func(ctx context.Context, status, city string) ([]model.Dispatch, error)
	createFn       func(ctx context.Context, req model.DispatchCreateRequest) (model.Dispatch, error)
	getByIDFn      func(ctx context.Context, dispatchID int64) (model.Dispatch, error)
	updateByIDFn   func(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error)
	updateStatusFn func(ctx context.Context, dispatchID int64, req model.DispatchStatusUpdateRequest) (model.Dispatch, error)
	deleteByIDFn   func(ctx context.Context, dispatchID int64) error
}

func (s dispatchServiceStub) List(ctx context.Context, status, city string) ([]model.Dispatch, error) {
	if s.listFn != nil {
		return s.listFn(ctx, status, city)
	}
	return nil, nil
}

func (s dispatchServiceStub) Create(ctx context.Context, req model.DispatchCreateRequest) (model.Dispatch, error) {
	if s.createFn != nil {
		return s.createFn(ctx, req)
	}
	return model.Dispatch{}, nil
}

func (s dispatchServiceStub) GetByID(ctx context.Context, dispatchID int64) (model.Dispatch, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, dispatchID)
	}
	return model.Dispatch{}, nil
}

func (s dispatchServiceStub) UpdateByID(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error) {
	if s.updateByIDFn != nil {
		return s.updateByIDFn(ctx, dispatchID, req)
	}
	return model.Dispatch{}, nil
}

func (s dispatchServiceStub) UpdateStatusByID(ctx context.Context, dispatchID int64, req model.DispatchStatusUpdateRequest) (model.Dispatch, error) {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(ctx, dispatchID, req)
	}
	return model.Dispatch{}, nil
}

func (s dispatchServiceStub) DeleteByID(ctx context.Context, dispatchID int64) error {
	if s.deleteByIDFn != nil {
		return s.deleteByIDFn(ctx, dispatchID)
	}
	return nil
}

func buildDispatchRouter(svc service.DispatchService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	dispatchHandler := handler.NewDispatchManagementAPI(svc)
	vehicleHandler := handler.NewVehicleManagementAPI(vehicleServiceStub{})
	handler.NewRouterWithGinEngine(engine, handler.ApiHandleFunctions{
		DispatchManagementAPI: dispatchHandler,
		VehicleManagementAPI:  vehicleHandler,
	})
	return engine
}
