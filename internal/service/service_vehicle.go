package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/mapper"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/repository"
)

var ErrBadRequest = errors.New("bad request")
var ErrNotFound = errors.New("not found")
var ErrConflict = errors.New("conflict")

type VehicleService interface {
	List(ctx context.Context, status, station string) ([]model.Vehicle, error)
	Create(ctx context.Context, req model.VehicleCreateRequest) (model.Vehicle, error)
	GetByID(ctx context.Context, vehicleID int64) (model.Vehicle, error)
	UpdateByID(ctx context.Context, vehicleID int64, req model.VehicleUpdateRequest) (model.Vehicle, error)
	UpdateStatusByID(ctx context.Context, vehicleID int64, req model.VehicleStatusUpdateRequest) (model.Vehicle, error)
	DeleteByID(ctx context.Context, vehicleID int64) error
}

type serviceVehicle struct {
	repo repository.VehicleRepository
}

func NewVehicleService(repo repository.VehicleRepository) VehicleService {
	return &serviceVehicle{repo: repo}
}

func (s *serviceVehicle) List(ctx context.Context, status, station string) ([]model.Vehicle, error) {
	if status != "" && !isValidStatus(model.VehicleStatus(status)) {
		return nil, fmt.Errorf("%w: invalid status filter", ErrBadRequest)
	}

	items, err := s.repo.List(ctx, status, station)
	if err != nil {
		return nil, err
	}

	return mapper.ToVehicleModels(items), nil
}

func (s *serviceVehicle) Create(ctx context.Context, req model.VehicleCreateRequest) (model.Vehicle, error) {
	if err := validateCreateRequest(req); err != nil {
		return model.Vehicle{}, err
	}

	entity := mapper.ToVehicleEntityFromCreate(0, req)
	created, err := s.repo.Create(ctx, entity)
	if err != nil {
		return model.Vehicle{}, translateRepoErr(err)
	}

	return mapper.ToVehicleModel(created), nil
}

func (s *serviceVehicle) GetByID(ctx context.Context, vehicleID int64) (model.Vehicle, error) {
	if vehicleID <= 0 {
		return model.Vehicle{}, fmt.Errorf("%w: vehicleId must be positive", ErrBadRequest)
	}

	item, err := s.repo.GetByID(ctx, vehicleID)
	if err != nil {
		return model.Vehicle{}, translateRepoErr(err)
	}

	return mapper.ToVehicleModel(item), nil
}

func (s *serviceVehicle) UpdateByID(ctx context.Context, vehicleID int64, req model.VehicleUpdateRequest) (model.Vehicle, error) {
	if vehicleID <= 0 {
		return model.Vehicle{}, fmt.Errorf("%w: vehicleId must be positive", ErrBadRequest)
	}
	if err := validateUpdateRequest(req); err != nil {
		return model.Vehicle{}, err
	}

	entity := mapper.ToVehicleEntityFromUpdate(vehicleID, req)
	updated, err := s.repo.UpdateByID(ctx, vehicleID, entity)
	if err != nil {
		return model.Vehicle{}, translateRepoErr(err)
	}

	return mapper.ToVehicleModel(updated), nil
}

func (s *serviceVehicle) UpdateStatusByID(ctx context.Context, vehicleID int64, req model.VehicleStatusUpdateRequest) (model.Vehicle, error) {
	if vehicleID <= 0 {
		return model.Vehicle{}, fmt.Errorf("%w: vehicleId must be positive", ErrBadRequest)
	}
	if !isValidStatus(req.Status) {
		return model.Vehicle{}, fmt.Errorf("%w: invalid status value", ErrBadRequest)
	}

	updated, err := s.repo.UpdateStatusByID(ctx, vehicleID, string(req.Status))
	if err != nil {
		return model.Vehicle{}, translateRepoErr(err)
	}

	return mapper.ToVehicleModel(updated), nil
}

func (s *serviceVehicle) DeleteByID(ctx context.Context, vehicleID int64) error {
	if vehicleID <= 0 {
		return fmt.Errorf("%w: vehicleId must be positive", ErrBadRequest)
	}

	if err := s.repo.DeleteByID(ctx, vehicleID); err != nil {
		return translateRepoErr(err)
	}
	return nil
}

func validateCreateRequest(req model.VehicleCreateRequest) error {
	if strings.TrimSpace(req.CallSign) == "" || strings.TrimSpace(req.VehicleType) == "" ||
		strings.TrimSpace(req.PlateNumber) == "" || strings.TrimSpace(req.Station) == "" {
		return fmt.Errorf("%w: callSign, vehicleType, plateNumber and station are required", ErrBadRequest)
	}
	if !isValidStatus(req.Status) {
		return fmt.Errorf("%w: invalid status value", ErrBadRequest)
	}
	if req.MileageKm < 0 {
		return fmt.Errorf("%w: mileageKm must be >= 0", ErrBadRequest)
	}
	if _, err := time.Parse("2006-01-02", req.LastServiceDate); err != nil {
		return fmt.Errorf("%w: lastServiceDate must use format YYYY-MM-DD", ErrBadRequest)
	}
	return nil
}

func validateUpdateRequest(req model.VehicleUpdateRequest) error {
	createEquivalent := model.VehicleCreateRequest{
		CallSign:        req.CallSign,
		VehicleType:     req.VehicleType,
		PlateNumber:     req.PlateNumber,
		Station:         req.Station,
		AssignedCrew:    req.AssignedCrew,
		Status:          req.Status,
		MileageKm:       req.MileageKm,
		LastServiceDate: req.LastServiceDate,
		Notes:           req.Notes,
	}
	return validateCreateRequest(createEquivalent)
}

func isValidStatus(status model.VehicleStatus) bool {
	switch status {
	case model.AVAILABLE, model.ON_MISSION, model.OUT_OF_SERVICE, model.IN_SERVICE:
		return true
	default:
		return false
	}
}

func translateRepoErr(err error) error {
	if errors.Is(err, repository.ErrVehicleNotFound) {
		return fmt.Errorf("%w: vehicle not found", ErrNotFound)
	}
	if errors.Is(err, repository.ErrVehicleConflict) {
		return fmt.Errorf("%w: vehicle already exists", ErrConflict)
	}
	return err
}
