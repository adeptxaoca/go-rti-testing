package main

const (
	PriceTypeCost              = "COST"
	PriceTypeDiscount          = "DISCOUNT"
	OperatorEqual              = "EQ"
	OperatorGreaterThanOrEqual = "GTE"
	OperatorLessThanOrEqual    = "LTE"
)

type RuleApplicability struct {
	CodeName string `json:"codeName"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type Price struct {
	Cost                float64             `json:"cost"`
	PriceType           string              `json:"priceType,omitempty"`
	RuleApplicabilities []RuleApplicability `json:"ruleApplicabilities,omitempty"`
}

type Component struct {
	Name   string  `json:"name"`
	IsMain bool    `json:"isMain,omitempty"`
	Prices []Price `json:"prices"`
}

type Product struct {
	Name       string      `json:"name"`
	Components []Component `json:"components"`
}

type Condition struct {
	RuleName string `json:"ruleName"`
	Value    string `json:"value"`
}

type Offer struct {
	Product
	TotalCost Price `json:"totalCost"`
}
