package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

type apiDispatchManagementImpl struct {
	dispatchService service.DispatchService
}

func NewDispatchManagementAPI(dispatchService service.DispatchService) DispatchManagementAPI {
	return &apiDispatchManagementImpl{dispatchService: dispatchService}
}

func (h *apiDispatchManagementImpl) DispatchesGet(c *gin.Context) {
	status := c.Query("status")
	city := c.Query("city")

	items, err := h.dispatchService.List(c.Request.Context(), status, city)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *apiDispatchManagementImpl) DispatchesPost(c *gin.Context) {
	var req model.DispatchCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid dispatch payload")
		return
	}

	item, err := h.dispatchService.Create(c.Request.Context(), req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *apiDispatchManagementImpl) DispatchesDispatchIdDelete(c *gin.Context) {
	id, err := parseDispatchID(c.Param("dispatchId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	if err := h.dispatchService.DeleteByID(c.Request.Context(), id); err != nil {
		respondWithError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *apiDispatchManagementImpl) DispatchesDispatchIdGet(c *gin.Context) {
	id, err := parseDispatchID(c.Param("dispatchId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	item, err := h.dispatchService.GetByID(c.Request.Context(), id)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *apiDispatchManagementImpl) DispatchesDispatchIdPut(c *gin.Context) {
	id, err := parseDispatchID(c.Param("dispatchId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	var req model.DispatchUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid dispatch payload")
		return
	}

	item, err := h.dispatchService.UpdateByID(c.Request.Context(), id, req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *apiDispatchManagementImpl) DispatchesDispatchIdStatusPatch(c *gin.Context) {
	id, err := parseDispatchID(c.Param("dispatchId"))
	if err != nil {
		respondWithBadRequest(c, err.Error())
		return
	}

	var req model.DispatchStatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondWithBadRequest(c, "invalid dispatch status update payload")
		return
	}

	item, err := h.dispatchService.UpdateStatusByID(c.Request.Context(), id, req)
	if err != nil {
		respondWithError(c, err)
		return
	}

	c.JSON(http.StatusOK, item)
}

func parseDispatchID(raw string) (int64, error) {
	id, err := strconv.ParseInt(raw, 10, 64)
	if err != nil || id <= 0 {
		return 0, fmt.Errorf("invalid dispatchId")
	}
	return id, nil
}
