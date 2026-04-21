package mapper

import (
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
)

func ToDispatchModel(e entity.DispatchEntity) model.Dispatch {
	return model.Dispatch{
		Id:                  e.DispatchID,
		IncidentNumber:      e.IncidentNumber,
		CallerName:          e.CallerName,
		PatientName:         e.PatientName,
		StreetAddress:       e.StreetAddress,
		City:                e.City,
		DispatchReason:      e.DispatchReason,
		Priority:            model.DispatchPriority(e.Priority),
		Status:              model.DispatchStatus(e.Status),
		AmbulanceCallSign:   e.AmbulanceCallSign,
		DestinationHospital: e.DestinationHospital,
		DispatcherName:      e.DispatcherName,
		CreatedAt:           e.CreatedAt,
		UpdatedAt:           e.UpdatedAt,
		Notes:               e.Notes,
	}
}

func ToDispatchModels(items []entity.DispatchEntity) []model.Dispatch {
	result := make([]model.Dispatch, 0, len(items))
	for _, item := range items {
		result = append(result, ToDispatchModel(item))
	}
	return result
}

func ToDispatchEntityFromCreate(id int64, req model.DispatchCreateRequest) entity.DispatchEntity {
	return entity.DispatchEntity{
		DispatchID:          id,
		IncidentNumber:      req.IncidentNumber,
		CallerName:          req.CallerName,
		PatientName:         req.PatientName,
		StreetAddress:       req.StreetAddress,
		City:                req.City,
		DispatchReason:      req.DispatchReason,
		Priority:            string(req.Priority),
		Status:              string(req.Status),
		AmbulanceCallSign:   req.AmbulanceCallSign,
		DestinationHospital: req.DestinationHospital,
		DispatcherName:      req.DispatcherName,
		CreatedAt:           req.CreatedAt,
		UpdatedAt:           req.CreatedAt,
		Notes:               req.Notes,
	}
}

func ToDispatchEntityFromUpdate(id int64, req model.DispatchUpdateRequest, updatedAt time.Time) entity.DispatchEntity {
	return entity.DispatchEntity{
		DispatchID:          id,
		IncidentNumber:      req.IncidentNumber,
		CallerName:          req.CallerName,
		PatientName:         req.PatientName,
		StreetAddress:       req.StreetAddress,
		City:                req.City,
		DispatchReason:      req.DispatchReason,
		Priority:            string(req.Priority),
		Status:              string(req.Status),
		AmbulanceCallSign:   req.AmbulanceCallSign,
		DestinationHospital: req.DestinationHospital,
		DispatcherName:      req.DispatcherName,
		CreatedAt:           req.CreatedAt,
		UpdatedAt:           updatedAt,
		Notes:               req.Notes,
	}
}
