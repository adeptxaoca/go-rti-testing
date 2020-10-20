package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
)

type CalculateRequest struct {
	Product    Product     `json:"product"`
	Conditions []Condition `json:"conditions"`
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", ping)
	mux.HandleFunc("/calculate", calculate)

	srv := http.Server{Addr: ":8080", Handler: mux}
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

func ping(w http.ResponseWriter, _ *http.Request) {
	_, _ = fmt.Fprint(w, "pong")
}

func calculate(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		_, _ = w.Write([]byte(http.StatusText(http.StatusMethodNotAllowed)))
		return
	}

	var calcReq CalculateRequest
	if err := decodeJson(w, req, &calcReq); err != nil {
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	offer, err := Calculate(&calcReq.Product, calcReq.Conditions)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	if offer != nil {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(offer); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}
}

func decodeJson(w http.ResponseWriter, req *http.Request, dst interface{}) error {
	if contentType := req.Header.Get("Content-Type"); contentType != "application/json" {
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return errors.New("content-type header is not application/json")
	}

	err := json.NewDecoder(req.Body).Decode(dst)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("invalid request body")
	}

	return nil
}

func Calculate(product *Product, conditions []Condition) (offer *Offer, err error) {
	if product == nil {
		return
	}

	totalCost, components, err := componentSearch(product.Components, conditions)
	if err != nil || totalCost == nil {
		return
	}

	offer = &Offer{TotalCost: *totalCost}
	offer.Product.Name = product.Name
	offer.Product.Components = components

	return offer, nil
}

// Function searches for suitable components and calculates total cost.
func componentSearch(components []Component, conditions []Condition) (*Price, []Component, error) {
	var totalCost float64
	var relevant []Component

	for _, component := range components {
		valid, cost, err := validateComponent(component, conditions)
		if err != nil {
			return new(Price), nil, err
		}

		if !valid {
			if component.IsMain {
				return nil, nil, nil
			}
			continue
		}

		relevant = append(relevant,
			Component{
				Name:   component.Name,
				IsMain: component.IsMain,
				Prices: []Price{{Cost: cost}},
			})
		totalCost += cost
	}

	return &Price{Cost: totalCost}, relevant, nil
}

// Function checks the component and returns the discounted cost.
func validateComponent(component Component, conditions []Condition) (bool, float64, error) {
	var cost, discount float64

	for _, price := range component.Prices {
		match, err := check(price.RuleApplicabilities, conditions)
		if err != nil {
			return false, 0, err
		}

		switch strings.ToUpper(price.PriceType) {
		case PriceTypeCost:
			if !match {
				continue
			}
			if cost > 0 {
				return false, 0, nil
			}
			cost = price.Cost
		case PriceTypeDiscount:
			if match && discount < price.Cost {
				discount = price.Cost
			}
		}
	}

	if cost == 0 {
		return false, 0, nil
	}

	return true, discountedCost(cost, discount), nil
}

// Function checks conditions according to selected rules.
// If there are no conditions or all conditions are met, then returns true.
func check(rules []RuleApplicability, conditions []Condition) (bool, error) {
	if len(conditions) == 0 {
		return true, nil
	}

	ruleMap := make(map[string]RuleApplicability)
	for _, rule := range rules {
		ruleMap[strings.ToLower(rule.CodeName)] = rule
	}

	for _, condition := range conditions {
		if rule, ok := ruleMap[strings.ToLower(condition.RuleName)]; ok {
			met, err := conditionCheckByRule(condition, rule)
			if err != nil || !met {
				return false, err
			}
		}
	}

	return true, nil
}

// Function performs the condition check by the rule.
// If the condition is fulfilled, then returns true.
func conditionCheckByRule(condition Condition, rule RuleApplicability) (bool, error) {
	switch rule.Operator {
	case OperatorEqual:
		return rule.Value == condition.Value, nil
	case OperatorLessThanOrEqual:
		return lessThanOrEqual(condition.Value, rule.Value)
	case OperatorGreaterThanOrEqual:
		return greaterThanOrEqual(condition.Value, rule.Value)
	default:
		return false, nil
	}
}

// v1 <= v2
func lessThanOrEqual(v1, v2 string) (bool, error) {
	v1f, v2f, err := parseValues(v1, v2)
	if err != nil {
		return false, err
	}

	return v1f <= v2f, nil
}

// v1 >= v2
func greaterThanOrEqual(v1, v2 string) (bool, error) {
	v1f, v2f, err := parseValues(v1, v2)
	if err != nil {
		return false, err
	}

	return v1f >= v2f, nil
}

// string to float64
func parseValues(a, b string) (float64, float64, error) {
	af, err := strconv.ParseFloat(a, 64)
	if err != nil {
		return 0, 0, err
	}

	bf, err := strconv.ParseFloat(b, 64)
	if err != nil {
		return 0, 0, err
	}

	return af, bf, nil
}

// Calculate discounted cost
func discountedCost(cost, discount float64) float64 {
	if discount > 100 {
		discount = 100
	}

	return math.Round(cost*(100-discount)) / 100
}
