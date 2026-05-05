package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/model"
	"github.com/olafser/ambulance-management-api-app/internal/service"
)

func TestDispatchesPost_ReturnsCreatedDispatch(t *testing.T) {
	router := buildDispatchRouter(dispatchServiceStub{
		createFn: func(ctx context.Context, req model.DispatchCreateRequest) (model.Dispatch, error) {
			return model.Dispatch{Id: 21, IncidentNumber: req.IncidentNumber, Status: req.Status}, nil
		},
	})

	payload := []byte(`{"incidentNumber":"DIS-2026-0001","callerName":"Caller","patientName":"Patient","streetAddress":"Main 1","city":"Bratislava","dispatchReason":"Chest pain","priority":"HIGH","status":"ACCEPTED","ambulanceCallSign":"AMB-01","destinationHospital":"Hospital","dispatcherName":"Operator","createdAt":"2026-04-20T09:05:00Z"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/dispatches", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, res.Code)
	}

	var got model.Dispatch
	if err := json.Unmarshal(res.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if got.Id != 21 || got.IncidentNumber != "DIS-2026-0001" {
		t.Fatalf("unexpected response payload: %+v", got)
	}
}

func TestDispatchesPost_ReturnsBadRequestOnInvalidPayload(t *testing.T) {
	router := buildDispatchRouter(dispatchServiceStub{})

	req := httptest.NewRequest(http.MethodPost, "/api/dispatches", bytes.NewReader([]byte(`{"incidentNumber":`)))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestDispatchesDispatchIdGet_ReturnsNotFoundWhenServiceReturnsNotFound(t *testing.T) {
	router := buildDispatchRouter(dispatchServiceStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (model.Dispatch, error) {
			return model.Dispatch{}, fmt.Errorf("%w: missing", service.ErrNotFound)
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/dispatches/999", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, res.Code)
	}
}

func TestDispatchesDispatchIdGet_ReturnsDispatch(t *testing.T) {
	now := time.Now().UTC()
	router := buildDispatchRouter(dispatchServiceStub{
		getByIDFn: func(ctx context.Context, dispatchID int64) (model.Dispatch, error) {
			return model.Dispatch{Id: dispatchID, IncidentNumber: "DIS-1", Status: model.ON_ROUTE, CreatedAt: now, UpdatedAt: now}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/dispatches/5", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}
