package mapper

import (
	"github.com/olafser/ambulance-management-api-app/internal/entity"
	"github.com/olafser/ambulance-management-api-app/internal/model"
)

func ToVehicleModel(e entity.VehicleEntity) model.Vehicle {
	return model.Vehicle{
		Id:              e.VehicleID,
		CallSign:        e.CallSign,
		VehicleType:     e.VehicleType,
		PlateNumber:     e.PlateNumber,
		Station:         e.Station,
		AssignedCrew:    e.AssignedCrew,
		Status:          model.VehicleStatus(e.Status),
		MileageKm:       e.MileageKm,
		LastServiceDate: e.LastServiceDate,
		Notes:           e.Notes,
	}
}

func ToVehicleModels(items []entity.VehicleEntity) []model.Vehicle {
	result := make([]model.Vehicle, 0, len(items))
	for _, item := range items {
		result = append(result, ToVehicleModel(item))
	}
	return result
}

func ToVehicleEntityFromCreate(id int64, req model.VehicleCreateRequest) entity.VehicleEntity {
	crew := req.AssignedCrew
	if crew == "" {
		crew = "Unassigned"
	}

	return entity.VehicleEntity{
		VehicleID:       id,
		CallSign:        req.CallSign,
		VehicleType:     req.VehicleType,
		PlateNumber:     req.PlateNumber,
		Station:         req.Station,
		AssignedCrew:    crew,
		Status:          string(req.Status),
		MileageKm:       req.MileageKm,
		LastServiceDate: req.LastServiceDate,
		Notes:           req.Notes,
	}
}

func ToVehicleEntityFromUpdate(id int64, req model.VehicleUpdateRequest) entity.VehicleEntity {
	crew := req.AssignedCrew
	if crew == "" {
		crew = "Unassigned"
	}

	return entity.VehicleEntity{
		VehicleID:       id,
		CallSign:        req.CallSign,
		VehicleType:     req.VehicleType,
		PlateNumber:     req.PlateNumber,
		Station:         req.Station,
		AssignedCrew:    crew,
		Status:          string(req.Status),
		MileageKm:       req.MileageKm,
		LastServiceDate: req.LastServiceDate,
		Notes:           req.Notes,
	}
}
