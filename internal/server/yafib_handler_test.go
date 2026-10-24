package server_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/jyork/yafib/v2/internal/engine"
	"github.com/jyork/yafib/v2/internal/server"
)

type generatorFunc func(int) ([]uint64, error)

func (f generatorFunc) CalculateFibonacciSequence(count int) ([]uint64, error) {
	return f(count)
}

func TestYafibHandler(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		query          string
		wantStatus     int
		wantCount      int
		wantCall       bool
		wantNumbers    []uint64
		wantError      string
		generatorError error
	}{
		{name: "GET", method: http.MethodGet, query: "?count=3", wantStatus: 200, wantCount: 3, wantCall: true, wantNumbers: []uint64{0, 1, 1}},
		{name: "POST rejected", method: http.MethodPost, query: "?count=3", wantStatus: 405, wantError: "Method not allowed. Use GET."},
		{name: "zero", method: http.MethodGet, query: "?count=0", wantStatus: 200, wantCount: 0, wantCall: true, wantNumbers: []uint64{}},
		{name: "missing count", method: http.MethodGet, wantStatus: 400, wantError: "Missing required query parameter 'count'"},
		{name: "empty count", method: http.MethodGet, query: "?count=", wantStatus: 400, wantError: "Missing required query parameter 'count'"},
		{name: "malformed count", method: http.MethodGet, query: "?count=abc", wantStatus: 400, wantError: "Parameter 'count' must be a valid integer"},
		{name: "negative count", method: http.MethodGet, query: "?count=-1", wantStatus: 400, wantError: "Parameter 'count' cannot be negative"},
		{name: "unsupported method", method: http.MethodPut, query: "?count=3", wantStatus: 405, wantError: "Method not allowed. Use GET."},
		{name: "wrapped validation error", method: http.MethodGet, query: "?count=95", wantStatus: 400, wantCount: 95, wantCall: true, generatorError: fmt.Errorf("out of range: %w", engine.ErrInvalidCount), wantError: "out of range: invalid Fibonacci count"},
		{name: "unexpected error", method: http.MethodGet, query: "?count=4", wantStatus: 500, wantCount: 4, wantCall: true, generatorError: errors.New("calculation failed"), wantError: "Failed to generate Fibonacci sequence"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			generator := generatorFunc(func(count int) ([]uint64, error) {
				calls++
				if count != tt.wantCount {
					t.Errorf("generator count = %d, want %d", count, tt.wantCount)
				}
				return tt.wantNumbers, tt.generatorError
			})
			rr := httptest.NewRecorder()
			server.NewServer(generator).YafibHandler(rr, httptest.NewRequest(tt.method, "/api/yafib"+tt.query, nil))
			if rr.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", rr.Code, tt.wantStatus, rr.Body.String())
			}
			wantCalls := 0
			if tt.wantCall {
				wantCalls = 1
			}
			if calls != wantCalls {
				t.Errorf("generator calls = %d, want %d", calls, wantCalls)
			}
			if got := rr.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q", got)
			}
			if tt.wantStatus == http.StatusMethodNotAllowed && rr.Header().Get("Allow") != "GET" {
				t.Errorf("Allow = %q", rr.Header().Get("Allow"))
			}
			if tt.wantStatus == http.StatusOK {
				var body server.SuccessResponse
				if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(body.Numbers, tt.wantNumbers) {
					t.Errorf("numbers = %v, want %v", body.Numbers, tt.wantNumbers)
				}
			} else {
				var body server.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if body.Error != tt.wantError {
					t.Errorf("error = %q, want %q", body.Error, tt.wantError)
				}
			}
		})
	}
}

func TestYafibHandlerWithEngine(t *testing.T) {
	handler := server.NewServer(engine.NewFibonacciEngine())
	for _, count := range []int{0, 94, 95} {
		t.Run(fmt.Sprint(count), func(t *testing.T) {
			rr := httptest.NewRecorder()
			handler.YafibHandler(rr, httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/yafib?count=%d", count), nil))
			if count == 95 {
				if rr.Code != http.StatusBadRequest {
					t.Fatalf("status = %d, want 400", rr.Code)
				}
				var body server.ErrorResponse
				if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
					t.Fatal(err)
				}
				if !strings.Contains(body.Error, "maximum allowed count of 94") {
					t.Errorf("unexpected error: %q", body.Error)
				}
				return
			}
			if rr.Code != http.StatusOK {
				t.Fatalf("status = %d, want 200", rr.Code)
			}
			var body server.SuccessResponse
			if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
				t.Fatal(err)
			}
			if body.Numbers == nil || len(body.Numbers) != count {
				t.Fatalf("numbers = %v, want %d terms", body.Numbers, count)
			}
			if count == 94 && body.Numbers[93] != 12200160415121876738 {
				t.Errorf("incorrect final value: %d", body.Numbers[93])
			}
		})
	}
}
