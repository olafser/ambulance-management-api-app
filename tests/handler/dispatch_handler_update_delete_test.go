package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/olafser/ambulance-management-api-app/internal/model"
)

func TestDispatchesDispatchIdPut_UsesPathIDAndReturnsUpdatedDispatch(t *testing.T) {
	var capturedID int64

	router := buildDispatchRouter(dispatchServiceStub{
		updateByIDFn: func(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error) {
			capturedID = dispatchID
			return model.Dispatch{Id: dispatchID, IncidentNumber: req.IncidentNumber, Status: req.Status}, nil
		},
	})

	payload := []byte(`{"incidentNumber":"DIS-2026-0002","callerName":"Caller","patientName":"Patient","streetAddress":"Main 2","city":"Nitra","dispatchReason":"Injury","priority":"MEDIUM","status":"ON_ROUTE","ambulanceCallSign":"AMB-02","destinationHospital":"Hospital","dispatcherName":"Operator","createdAt":"2026-04-20T10:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/dispatches/8", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if capturedID != 8 {
		t.Fatalf("expected service to receive id=8, got id=%d", capturedID)
	}
}

func TestDispatchesDispatchIdDelete_ReturnsNoContent(t *testing.T) {
	var capturedID int64
	router := buildDispatchRouter(dispatchServiceStub{
		deleteByIDFn: func(ctx context.Context, dispatchID int64) error {
			capturedID = dispatchID
			return nil
		},
	})

	req := httptest.NewRequest(http.MethodDelete, "/api/dispatches/12", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, res.Code)
	}
	if capturedID != 12 {
		t.Fatalf("expected service to receive id=12, got id=%d", capturedID)
	}
}

func TestDispatchesDispatchIdPut_ReturnsBadRequestOnInvalidID(t *testing.T) {
	router := buildDispatchRouter(dispatchServiceStub{})
	payload := []byte(`{"incidentNumber":"DIS-3","callerName":"A","patientName":"B","streetAddress":"Main 3","city":"Kosice","dispatchReason":"Pain","priority":"LOW","status":"ACCEPTED","ambulanceCallSign":"AMB-03","destinationHospital":"Hospital","dispatcherName":"Operator","createdAt":"2026-04-20T12:00:00Z"}`)

	req := httptest.NewRequest(http.MethodPut, "/api/dispatches/0", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}

func TestDispatchesDispatchIdPut_PreservesHTTPFlow(t *testing.T) {
	now := time.Now().UTC()
	router := buildDispatchRouter(dispatchServiceStub{
		updateByIDFn: func(ctx context.Context, dispatchID int64, req model.DispatchUpdateRequest) (model.Dispatch, error) {
			return model.Dispatch{Id: dispatchID, IncidentNumber: req.IncidentNumber, Status: req.Status, CreatedAt: now, UpdatedAt: now}, nil
		},
	})

	payload := []byte(`{"incidentNumber":"DIS-2026-1234","callerName":"A","patientName":"B","streetAddress":"Main 4","city":"Trnava","dispatchReason":"Syncope","priority":"HIGH","status":"ON_SCENE","ambulanceCallSign":"AMB-99","destinationHospital":"Hospital","dispatcherName":"Operator","createdAt":"2026-04-20T12:00:00Z"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/dispatches/99", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
}
