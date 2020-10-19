package main

import "testing"

var product = Product{
	Name: "Игровой",
	Components: []Component{
		{
			IsMain: true,
			Name:   "Интернет",
			Prices: []Price{
				{
					Cost:      100,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "adsl",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "10",
						},
					},
				},
				{
					Cost:      150,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "adsl",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "15",
						},
					},
				},
				{
					Cost:      500,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "xpon",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "100",
						},
					},
				},
				{
					Cost:      900,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "xpon",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "200",
						},
					},
				},
				{
					Cost:      200,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "fttb",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "30",
						},
					},
				},
				{
					Cost:      400,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "fttb",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "50",
						},
					},
				},
				{
					Cost:      600,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "fttb",
						},
						{
							CodeName: "internetSpeed",
							Operator: OperatorEqual,
							Value:    "200",
						},
					},
				},
				{
					Cost:      10,
					PriceType: PriceTypeDiscount,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "internetSpeed",
							Operator: OperatorGreaterThanOrEqual,
							Value:    "50",
						},
					},
				},
				{
					Cost:      15,
					PriceType: PriceTypeDiscount,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "internetSpeed",
							Operator: OperatorGreaterThanOrEqual,
							Value:    "100",
						},
					},
				},
			},
		},
		{
			Name: "ADSL Модем",
			Prices: []Price{
				{
					Cost:      300,
					PriceType: PriceTypeCost,
					RuleApplicabilities: []RuleApplicability{
						{
							CodeName: "technology",
							Operator: OperatorEqual,
							Value:    "adsl",
						},
					},
				},
			},
		},
	},
}

func TestCalculateNil(t *testing.T) {
	r, err := Calculate(nil, nil)
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r != nil {
		t.Error(r)
	}
}

func TestCalculateNotNil(t *testing.T) {
	r, err := Calculate(&Product{}, []Condition{})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r == nil {
		t.Error(r)
	}
}

func TestCalculate2Price(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "internetSpeed",
			Value:    "200",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r != nil {
		for _, component := range r.Components {
			if len(component.Prices) > 1 {
				t.Error("У компонента должна быть только 1 цена")
			}
		}
		t.Error("По данным условиям невозможно продать продукт")
	}
}

func TestCalculateADSL(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "adsl",
		},
		{
			RuleName: "internetSpeed",
			Value:    "10",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
	}
	if r == nil {
		t.Error("Неверно расчитанно предложение")
		return
	}
	if len(r.Components) != 2 {
		t.Error("Должно быть 2 компонента")
	}
	if r.TotalCost.Cost != 400 {
		t.Error("Неверно расчитана сумма")
	}

	for _, component := range r.Components {
		if len(component.Prices) > 1 {
			t.Error("У компонента должна быть только 1 цена")
		}
		if component.Name == "Интернет" && component.Prices[0].Cost != 100 {
			t.Error("Неверно расчитана цена компонента Интернет")
		}
	}
}

func TestCalculateADSL2(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "adsl",
		},
		{
			RuleName: "internetSpeed",
			Value:    "15",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
	}
	if r == nil {
		t.Error("Неверно расчитанно предложение")
		return
	}
	if len(r.Components) != 2 {
		t.Error("Должно быть 2 компонента")
	}
	if r.TotalCost.Cost != 450 {
		t.Error("Неверно расчитана сумма")
	}

	for _, component := range r.Components {
		if len(component.Prices) > 1 {
			t.Error("У компонента должна быть только 1 цена")
		}
		if component.Name == "Интернет" && component.Prices[0].Cost != 150 {
			t.Error("Неверно расчитана цена компонента Интернет")
		}
	}
}

func TestCalculateIsMain(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "adsl",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}

	if r != nil {
		t.Error("Должен исключиться обязательный компонент")
	}
}

func TestCalculateFTTB(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "fttb",
		},
		{
			RuleName: "internetSpeed",
			Value:    "30",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r == nil {
		t.Error("Неверно расчитанно предложение")
		return
	}
	if r.TotalCost.Cost != 200 {
		t.Error("Неверно расчитана сумма")
	}

	for _, component := range r.Components {
		if len(component.Prices) > 1 {
			t.Error("У компонента должна быть только 1 цена")
		}
		if component.Name == "Интернет" && component.Prices[0].Cost != 200 {
			t.Error("Неверно расчитана цена компонента Интернет")
		}
	}
}

func TestCalculateDiscount(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "xpon",
		},
		{
			RuleName: "internetSpeed",
			Value:    "200",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r == nil {
		t.Error("Неверно расчитанно предложение")
		return
	}
	if r.TotalCost.Cost != 765 {
		t.Error("Неверно расчитана сумма с учетом скидки")
	}

	for _, component := range r.Components {
		if len(component.Prices) > 1 {
			t.Error("У компонента должна быть только 1 цена")
		}
		if component.Name == "Интернет" && component.Prices[0].Cost != 765 {
			t.Error("Неверно расчитана цена компонента Интернет с учетом скидки")
		}
	}
}

func TestCalculateDiscount2(t *testing.T) {
	r, err := Calculate(&product, []Condition{
		{
			RuleName: "technology",
			Value:    "fttb",
		},
		{
			RuleName: "internetSpeed",
			Value:    "50",
		},
	})
	if err != nil {
		t.Error("Error calculating", err)
		return
	}
	if r == nil {
		t.Error("Неверно расчитанно предложение")
		return
	}
	if r.TotalCost.Cost != 360 {
		t.Error("Неверно расчитана сумма с учетом скидки")
	}

	for _, component := range r.Components {
		if len(component.Prices) > 1 {
			t.Error("У компонента должна быть только 1 цена")
		}
		if component.Name == "Интернет" && component.Prices[0].Cost != 360 {
			t.Error("Неверно расчитана цена компонента Интернет с учетом скидки")
		}
	}
}
