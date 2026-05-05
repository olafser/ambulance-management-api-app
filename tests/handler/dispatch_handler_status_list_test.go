package handler_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/olafser/ambulance-management-api-app/internal/model"
)

func TestDispatchesGet_PassesFiltersToService(t *testing.T) {
	var capturedStatus string
	var capturedCity string

	router := buildDispatchRouter(dispatchServiceStub{
		listFn: func(ctx context.Context, status, city string) ([]model.Dispatch, error) {
			capturedStatus = status
			capturedCity = city
			return []model.Dispatch{}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/dispatches?status=ON_ROUTE&city=Nitra", nil)
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if capturedStatus != "ON_ROUTE" || capturedCity != "Nitra" {
		t.Fatalf("unexpected query propagation: status=%s city=%s", capturedStatus, capturedCity)
	}
}

func TestDispatchesDispatchIdStatusPatch_UsesPathIDAndReturnsUpdatedDispatch(t *testing.T) {
	var capturedID int64
	var capturedStatus model.DispatchStatus

	router := buildDispatchRouter(dispatchServiceStub{
		updateStatusFn: func(ctx context.Context, dispatchID int64, req model.DispatchStatusUpdateRequest) (model.Dispatch, error) {
			capturedID = dispatchID
			capturedStatus = req.Status
			return model.Dispatch{Id: dispatchID, Status: req.Status}, nil
		},
	})

	payload := []byte(`{"status":"COMPLETED"}`)
	req := httptest.NewRequest(http.MethodPatch, "/api/dispatches/7/status", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, res.Code)
	}
	if capturedID != 7 || capturedStatus != model.COMPLETED {
		t.Fatalf("unexpected service call args: id=%d status=%s", capturedID, capturedStatus)
	}
}

func TestDispatchesDispatchIdStatusPatch_ReturnsBadRequestOnInvalidPayload(t *testing.T) {
	router := buildDispatchRouter(dispatchServiceStub{})

	req := httptest.NewRequest(http.MethodPatch, "/api/dispatches/4/status", bytes.NewReader([]byte(`{"status":`)))
	req.Header.Set("Content-Type", "application/json")
	res := httptest.NewRecorder()

	router.ServeHTTP(res, req)

	if res.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, res.Code)
	}
}
