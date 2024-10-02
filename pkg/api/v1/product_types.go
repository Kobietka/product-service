package v1

type Product struct {
	Ean       string    `json:"ean"`
	Name      string    `json:"name"`
	Packaging Quantity  `json:"packaging"`
	Nutrition Nutrition `json:"nutrition"`
}

type Quantity struct {
	Value float32 `json:"value"`
	Unit  string  `json:"unit"`
}

type Nutrition struct {
	Per       Quantity   `json:"per"`
	Kcal      int32      `json:"kcal"`
	Nutrients []Nutrient `json:"nutrients"`
	Vitamins  []Vitamin  `json:"vitamins"`
	Minerals  []Mineral  `json:"minerals"`
}

type Nutrient struct {
	T        string   `json:"type"`
	Quantity Quantity `json:"quantity"`
}

type Mineral struct {
	T        string   `json:"type"`
	Quantity Quantity `json:"quantity"`
}

type Vitamin struct {
	T        string   `json:"type"`
	Quantity Quantity `json:"quantity"`
}
