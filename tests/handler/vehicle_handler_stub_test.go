package handler_test

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/olafser/ambulance-management-api-app/internal/handler"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

type vehicleServiceStub struct {
	listFn         func(ctx context.Context, status, station string) ([]model.Vehicle, error)
	createFn       func(ctx context.Context, req model.VehicleCreateRequest) (model.Vehicle, error)
	getByIDFn      func(ctx context.Context, vehicleID int64) (model.Vehicle, error)
	updateByIDFn   func(ctx context.Context, vehicleID int64, req model.VehicleUpdateRequest) (model.Vehicle, error)
	updateStatusFn func(ctx context.Context, vehicleID int64, req model.VehicleStatusUpdateRequest) (model.Vehicle, error)
	deleteByIDFn   func(ctx context.Context, vehicleID int64) error
}

func (s vehicleServiceStub) List(ctx context.Context, status, station string) ([]model.Vehicle, error) {
	if s.listFn != nil {
		return s.listFn(ctx, status, station)
	}
	return nil, nil
}

func (s vehicleServiceStub) Create(ctx context.Context, req model.VehicleCreateRequest) (model.Vehicle, error) {
	if s.createFn != nil {
		return s.createFn(ctx, req)
	}
	return model.Vehicle{}, nil
}

func (s vehicleServiceStub) GetByID(ctx context.Context, vehicleID int64) (model.Vehicle, error) {
	if s.getByIDFn != nil {
		return s.getByIDFn(ctx, vehicleID)
	}
	return model.Vehicle{}, nil
}

func (s vehicleServiceStub) UpdateByID(ctx context.Context, vehicleID int64, req model.VehicleUpdateRequest) (model.Vehicle, error) {
	if s.updateByIDFn != nil {
		return s.updateByIDFn(ctx, vehicleID, req)
	}
	return model.Vehicle{}, nil
}

func (s vehicleServiceStub) UpdateStatusByID(ctx context.Context, vehicleID int64, req model.VehicleStatusUpdateRequest) (model.Vehicle, error) {
	if s.updateStatusFn != nil {
		return s.updateStatusFn(ctx, vehicleID, req)
	}
	return model.Vehicle{}, nil
}

func (s vehicleServiceStub) DeleteByID(ctx context.Context, vehicleID int64) error {
	if s.deleteByIDFn != nil {
		return s.deleteByIDFn(ctx, vehicleID)
	}
	return nil
}

func buildVehicleRouter(svc service.VehicleService) *gin.Engine {
	gin.SetMode(gin.TestMode)
	engine := gin.New()
	vehicleHandler := handler.NewVehicleManagementAPI(svc)
	handler.NewRouterWithGinEngine(engine, handler.ApiHandleFunctions{VehicleManagementAPI: vehicleHandler})
	return engine
}
