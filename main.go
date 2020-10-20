package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/felixge/httpsnoop"

	"go-rti-testing/pkg/errors"
)

type CalculateRequest struct {
	Product    Product     `json:"product"`
	Conditions []Condition `json:"conditions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/calculate", calculate)

	srv := http.Server{Addr: ":8080", Handler: logRequest(mux)}
	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	log.Println("Starting http server on :8080...")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m := httpsnoop.CaptureMetrics(handler, w, r)
		log.Printf(
			"method=%s path=%s status=%d duration=%s written=%d",
			r.Method, r.URL, m.Code, m.Duration, m.Written,
		)
	})
}

func decodeJson(req *http.Request, dst interface{}) error {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		return errors.UnsupportedMediaType.New("application/json")
	}

	err := json.NewDecoder(req.Body).Decode(dst)
	if err != nil {
		return errors.BadRequest.New("invalid request body")
	}

	return nil
}

func encodeJson(w http.ResponseWriter, v interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		return errors.Internal.Wrap(err, "json.NewEncoder(w).Encode(offer)")
	}
	return nil
}

func newErrorResponse(message string) *ErrorResponse {
	return &ErrorResponse{
		Error: message,
	}
}

func httpError(w http.ResponseWriter, err error) {
	var status int
	switch errors.GetType(err) {
	case errors.UnsupportedMediaType:
		status = http.StatusUnsupportedMediaType
		_ = encodeJson(w, newErrorResponse(fmt.Sprintf(errors.MsgUnsupportedMediaType, err.Error())))
	case errors.MethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case errors.BadRequest:
		status = http.StatusBadRequest
		_ = encodeJson(w, newErrorResponse(err.Error()))
	default:
		status = http.StatusInternalServerError
		log.Printf("ERROR %s", err)
		_ = encodeJson(w, newErrorResponse(http.StatusText(http.StatusInternalServerError)))
	}
	w.WriteHeader(status)
}

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "pong")
}

func calculate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		httpError(w, errors.MethodNotAllowed.New(""))
		return
	}

	var calcReq CalculateRequest
	if err := decodeJson(req, &calcReq); err != nil {
		httpError(w, err)
		return
	}

	offer, err := Calculate(&calcReq.Product, calcReq.Conditions)
	if err != nil {
		httpError(w, err)
		return
	}

	if offer != nil {
		if err := encodeJson(w, offer); err != nil {
			httpError(w, err)
			return
		}
	}
}
