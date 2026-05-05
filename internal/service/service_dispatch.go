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

type DispatchService interface {
	List(ctx context.Context, status, city string) ([]model.Dispatch, error)
	Create(ctx context.Context, req model.DispatchCreateRequest) (model.Dispatch, error)
	GetByID(ctx context.Context, dispatchID int64) (model.Dispatch, error)
	UpdateByID(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error)
	UpdateStatusByID(ctx context.Context, dispatchID int64, req model.DispatchStatusUpdateRequest) (model.Dispatch, error)
	DeleteByID(ctx context.Context, dispatchID int64) error
}

type serviceDispatch struct {
	repo repository.DispatchRepository
}

func NewDispatchService(repo repository.DispatchRepository) DispatchService {
	return &serviceDispatch{repo: repo}
}

func (s *serviceDispatch) List(ctx context.Context, status, city string) ([]model.Dispatch, error) {
	if status != "" && !isValidDispatchStatus(model.DispatchStatus(status)) {
		return nil, fmt.Errorf("%w: invalid status filter", ErrBadRequest)
	}

	items, err := s.repo.List(ctx, status, city)
	if err != nil {
		return nil, err
	}

	return mapper.ToDispatchModels(items), nil
}

func (s *serviceDispatch) Create(ctx context.Context, req model.DispatchCreateRequest) (model.Dispatch, error) {
	if err := validateDispatchCreateRequest(req); err != nil {
		return model.Dispatch{}, err
	}

	entity := mapper.ToDispatchEntityFromCreate(0, req)
	created, err := s.repo.Create(ctx, entity)
	if err != nil {
		return model.Dispatch{}, translateDispatchRepoErr(err)
	}

	return mapper.ToDispatchModel(created), nil
}

func (s *serviceDispatch) GetByID(ctx context.Context, dispatchID int64) (model.Dispatch, error) {
	if dispatchID <= 0 {
		return model.Dispatch{}, fmt.Errorf("%w: dispatchId must be positive", ErrBadRequest)
	}

	item, err := s.repo.GetByID(ctx, dispatchID)
	if err != nil {
		return model.Dispatch{}, translateDispatchRepoErr(err)
	}

	return mapper.ToDispatchModel(item), nil
}

func (s *serviceDispatch) UpdateByID(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error) {
	if dispatchID <= 0 {
		return model.Dispatch{}, fmt.Errorf("%w: dispatchId must be positive", ErrBadRequest)
	}
	if err := validateDispatchUpdateRequest(req); err != nil {
		return model.Dispatch{}, err
	}

	entity := mapper.ToDispatchEntityFromUpdate(dispatchID, req, time.Now().UTC())
	updated, err := s.repo.UpdateByID(ctx, dispatchID, entity)
	if err != nil {
		return model.Dispatch{}, translateDispatchRepoErr(err)
	}

	return mapper.ToDispatchModel(updated), nil
}

func (s *serviceDispatch) UpdateStatusByID(ctx context.Context, dispatchID int64, req model.DispatchStatusUpdateRequest) (model.Dispatch, error) {
	if dispatchID <= 0 {
		return model.Dispatch{}, fmt.Errorf("%w: dispatchId must be positive", ErrBadRequest)
	}
	if !isValidDispatchStatus(req.Status) {
		return model.Dispatch{}, fmt.Errorf("%w: invalid status value", ErrBadRequest)
	}

	updated, err := s.repo.UpdateStatusByID(ctx, dispatchID, string(req.Status), time.Now().UTC())
	if err != nil {
		return model.Dispatch{}, translateDispatchRepoErr(err)
	}

	return mapper.ToDispatchModel(updated), nil
}

func (s *serviceDispatch) DeleteByID(ctx context.Context, dispatchID int64) error {
	if dispatchID <= 0 {
		return fmt.Errorf("%w: dispatchId must be positive", ErrBadRequest)
	}

	if err := s.repo.DeleteByID(ctx, dispatchID); err != nil {
		return translateDispatchRepoErr(err)
	}
	return nil
}

func validateDispatchCreateRequest(req model.DispatchCreateRequest) error {
	if strings.TrimSpace(req.IncidentNumber) == "" || strings.TrimSpace(req.CallerName) == "" ||
		strings.TrimSpace(req.PatientName) == "" || strings.TrimSpace(req.StreetAddress) == "" ||
		strings.TrimSpace(req.City) == "" || strings.TrimSpace(req.DispatchReason) == "" ||
		strings.TrimSpace(req.AmbulanceCallSign) == "" || strings.TrimSpace(req.DestinationHospital) == "" ||
		strings.TrimSpace(req.DispatcherName) == "" {
		return fmt.Errorf("%w: required dispatch fields are missing", ErrBadRequest)
	}
	if !isValidDispatchPriority(req.Priority) {
		return fmt.Errorf("%w: invalid priority value", ErrBadRequest)
	}
	if !isValidDispatchStatus(req.Status) {
		return fmt.Errorf("%w: invalid status value", ErrBadRequest)
	}
	if req.CreatedAt.IsZero() {
		return fmt.Errorf("%w: createdAt is required", ErrBadRequest)
	}
	return nil
}

func validateDispatchUpdateRequest(req model.DispatchUpdateRequest) error {
	createEquivalent := model.DispatchCreateRequest{
		IncidentNumber:      req.IncidentNumber,
		CallerName:          req.CallerName,
		PatientName:         req.PatientName,
		StreetAddress:       req.StreetAddress,
		City:                req.City,
		DispatchReason:      req.DispatchReason,
		Priority:            req.Priority,
		Status:              req.Status,
		AmbulanceCallSign:   req.AmbulanceCallSign,
		DestinationHospital: req.DestinationHospital,
		DispatcherName:      req.DispatcherName,
		CreatedAt:           req.CreatedAt,
		Notes:               req.Notes,
	}
	return validateDispatchCreateRequest(createEquivalent)
}

func isValidDispatchStatus(status model.DispatchStatus) bool {
	switch status {
	case model.ACCEPTED, model.ON_ROUTE, model.ON_SCENE, model.TRANSPORTING_TO_HOSPITAL, model.COMPLETED:
		return true
	default:
		return false
	}
}

func isValidDispatchPriority(priority model.DispatchPriority) bool {
	switch priority {
	case model.LOW, model.MEDIUM, model.HIGH, model.CRITICAL:
		return true
	default:
		return false
	}
}

func translateDispatchRepoErr(err error) error {
	if errors.Is(err, repository.ErrDispatchNotFound) {
		return fmt.Errorf("%w: dispatch not found", ErrNotFound)
	}
	if errors.Is(err, repository.ErrDispatchConflict) {
		return fmt.Errorf("%w: dispatch already exists", ErrConflict)
	}
	return err
}
