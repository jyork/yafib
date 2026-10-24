package server

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/jyork/yafib/v2/internal/engine"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

type SuccessResponse struct {
	Numbers []uint64 `json:"numbers"`
}

// SequenceGenerator supplies Fibonacci sequences to the HTTP handler.
type SequenceGenerator interface {
	CalculateFibonacciSequence(count int) ([]uint64, error)
}

type Server struct {
	generator SequenceGenerator
	logger    *slog.Logger
}

func NewServer(g SequenceGenerator, logger *slog.Logger) *Server {
	return &Server{generator: g, logger: logger}
}

// YafibHandler accepts GET with a count query parameter.
func (s *Server) YafibHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", "GET")
		s.respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed. Use GET.")
		return
	}

	queryParam := r.URL.Query().Get("count")
	if queryParam == "" {
		s.respondWithError(w, http.StatusBadRequest, "Missing required query parameter 'count'")
		return
	}

	count, err := strconv.Atoi(queryParam)
	if err != nil {
		s.respondWithError(w, http.StatusBadRequest, "Parameter 'count' must be a valid integer")
		return
	}
	if count < 0 {
		s.respondWithError(w, http.StatusBadRequest, "Parameter 'count' cannot be negative")
		return
	}

	numbers, err := s.generator.CalculateFibonacciSequence(count)
	if err != nil {
		if errors.Is(err, engine.ErrInvalidCount) {
			s.respondWithError(w, http.StatusBadRequest, err.Error())
			return
		}
		s.logger.ErrorContext(r.Context(), "calculation failed", "method", r.Method, "path", r.URL.Path, "count", count, "error", err)
		s.respondWithError(w, http.StatusInternalServerError, "Failed to generate Fibonacci sequence")
		return
	}

	s.writeJSON(w, http.StatusOK, SuccessResponse{Numbers: numbers})
}

func (s *Server) respondWithError(w http.ResponseWriter, statusCode int, message string) {
	s.writeJSON(w, statusCode, ErrorResponse{Error: message})
}

func (s *Server) writeJSON(w http.ResponseWriter, statusCode int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		// Headers have already been sent, so a second HTTP error cannot be written.
		s.logger.Error("response write failed", "status", statusCode, "error", err)
	}
}
