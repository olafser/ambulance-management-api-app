package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

type apiVehicleManagementImpl struct {
	vehicleService service.VehicleService
}

func NewVehicleManagementAPI(vehicleService service.VehicleService) VehicleManagementAPI {
	return &apiVehicleManagementImpl{vehicleService: vehicleService}
}

func (h *apiVehicleManagementImpl) VehiclesGet(c *gin.Context) {
	status := c.Query("status")
	station := c.Query("station")

	items, err := h.vehicleService.List(c.Request.Context(), status, station)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *apiVehicleManagementImpl) VehiclesPost(c *gin.Context) {
	var req model.VehicleCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid vehicle payload")
		return
	}

	item, err := h.vehicleService.Create(c.Request.Context(), req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *apiVehicleManagementImpl) VehiclesVehicleIdDelete(c *gin.Context) {
	id, err := parseVehicleID(c.Param("vehicleId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	if err := h.vehicleService.DeleteByID(c.Request.Context(), id); err != nil {
		respondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *apiVehicleManagementImpl) VehiclesVehicleIdGet(c *gin.Context) {
	id, err := parseVehicleID(c.Param("vehicleId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	item, err := h.vehicleService.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *apiVehicleManagementImpl) VehiclesVehicleIdPut(c *gin.Context) {
	id, err := parseVehicleID(c.Param("vehicleId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	var req model.VehicleUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid vehicle payload")
		return
	}

	item, err := h.vehicleService.UpdateByID(c.Request.Context(), id, req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *apiVehicleManagementImpl) VehiclesVehicleIdStatusPatch(c *gin.Context) {
	id, err := parseVehicleID(c.Param("vehiceId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	var req model.VehicleStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid status update payload")
		return
	}

	item, err := h.vehicleService.UpdateStatusByID(c.Request.Context(), id, req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func parseVehicleID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid vehicleId")
	}
	return id, nil
}

func respondWithError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, service.ErrBadRequest):
		c.JSON(http.StatusBadRequest, model.Error{Message: err.Error()})
	case errors.Is(err, service.ErrNotFound):
		c.JSON(http.StatusNotFound, model.Error{Message: err.Error()})
	case errors.Is(err, service.ErrConflict):
		c.JSON(http.StatusConflict, model.Error{Message: err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, model.Error{Message: "internal server error"})
	}
}

func respondWithBadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, model.Error{Message: message})
}
