package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/model"
)

func TestVehiclesVehicleIdStatusPatch_UsesPathIDAndReturnsUpdatedVehicle(t *testing.T) {
	var capturedID int64
	var capturedStatus model.VehicleStatus

	router := buildVehicleRouter(vehicleServiceStub{
		updateStatusFn: func(ctx context.Context, vehicleID int64, req model.VehicleStatusUpdateRequest) (model.Vehicle, error) {
			capturedID = vehicleID
			capturedStatus = req.Status
			return model.Vehicle{Id: vehicleID, Status: req.Status}, nil
		},
	})

	payload := []byte(`{"status":"IN_SERVICE"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/vehicles/7/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if capturedID != 7 || capturedStatus != model.IN_SERVICE {
		t.Fatalf("unexpected service call args: id=%d status=%s", capturedID, capturedStatus)
	}
}
