package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestVehiclesPost_ReturnsCreatedVehicle(t *testing.T) {
	router := buildVehicleRouter(vehicleServiceStub{
		createFn: func(ctx context.Context, req model.VehicleCreateRequest) (model.Vehicle, error) {
			return model.Vehicle{Id: 11, CallSign: req.CallSign, Status: req.Status}, nil
		},
	})

	payload := []byte(`{"callSign":"A-11","vehicleType":"Type B","plateNumber":"BA123AA","station":"North","status":"AVAILABLE","mileageKm":10,"lastServiceDate":"2026-04-01"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/vehicles", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var got model.Vehicle
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Id != 11 || got.CallSign != "A-11" {
		t.Fatalf("unexpected response body: %+v", got)
	}
}

func TestVehiclesVehicleIdGet_ReturnsNotFoundWhenServiceReturnsNotFound(t *testing.T) {
	router := buildVehicleRouter(vehicleServiceStub{
		getByIDFn: func(ctx context.Context, vehicleID int64) (model.Vehicle, error) {
			return model.Vehicle{}, fmt.Errorf("%w: vehicle not found", service.ErrNotFound)
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/vehicles/404", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}

	var got model.Error
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if got.Message == "" {
		t.Fatalf("expected non-empty error message")
	}
}
