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

func httpError(w http.ResponseWriter, err error) {
	switch errors.GetType(err) {
	case errors.UnsupportedMediaType:
		w.WriteHeader(http.StatusUnsupportedMediaType)
		_, _ = fmt.Fprintf(w, errors.MsgUnsupportedMediaType, err.Error())
	case errors.MethodNotAllowed:
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
	case errors.BadRequest:
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
	default:
		log.Printf("ERROR %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(http.StatusText(http.StatusInternalServerError)))
	}
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
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(offer); err != nil {
			httpError(w, errors.Internal.Wrap(err, "json.NewEncoder(w).Encode(offer)"))
			return
		}
	}
}
