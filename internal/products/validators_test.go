package products

import (
	v1 "github.com/Kobietka/product-service/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestValidateProduct(t *testing.T) {
	correctNutrition := v1.Nutrition{
		Per: v1.Quantity{
			Value: 12,
			Unit:  "g",
		},
		Kcal: 123,
		Nutrients: []v1.Nutrient{
			{
				T: "PROTEIN",
				Quantity: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
			},
			{
				T: "CARBOHYDRATES",
				Quantity: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
			},
			{
				T: "FAT",
				Quantity: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
			},
		},
	}

	correctPackaging := v1.Quantity{
		Value: 12,
		Unit:  "g",
	}

	tests := []struct {
		Name        string
		Product     v1.Product
		ExpectedErr error
	}{
		{
			Name: "product with ean-8",
			Product: v1.Product{
				Ean:       "12345678",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: nil,
		},
		{
			Name: "product with ean-13",
			Product: v1.Product{
				Ean:       "1234567890123",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: nil,
		},
		{
			Name: "product with upc-a",
			Product: v1.Product{
				Ean:       "123456789012",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: nil,
		},
		{
			Name: "product with invalid ean",
			Product: v1.Product{
				Ean:       "12345",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorProductEanInvalid,
		},
		{
			Name: "product with empty ean",
			Product: v1.Product{
				Ean:       "",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorProductEanMissing,
		},
		{
			Name: "product with blank ean",
			Product: v1.Product{
				Ean:       "    ",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorProductEanMissing,
		},
		{
			Name: "product with empty name",
			Product: v1.Product{
				Ean:       "12345678",
				Name:      "",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorProductNameMissing,
		},
		{
			Name: "product with blank name",
			Product: v1.Product{
				Ean:       "12345678",
				Name:      "    ",
				Nutrition: correctNutrition,
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorProductNameMissing,
		},
		{
			Name: "product nutrition is checked",
			Product: v1.Product{
				Ean:  "12345678",
				Name: "product name",
				Nutrition: v1.Nutrition{
					Per: v1.Quantity{
						Value: 12,
						Unit:  "",
					},
					Kcal: -123,
				},
				Packaging: correctPackaging,
			},
			ExpectedErr: ErrorQuantityUnitMissing,
		},
		{
			Name: "product packaging is checked",
			Product: v1.Product{
				Ean:       "1234567890123",
				Name:      "product name",
				Nutrition: correctNutrition,
				Packaging: v1.Quantity{
					Value: -12,
					Unit:  "g",
				},
			},
			ExpectedErr: validateQuantity(v1.Quantity{Value: -12, Unit: "g"}),
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateProduct(test.Product)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestValidateQuantity(t *testing.T) {
	tests := []struct {
		Name        string
		Quantity    v1.Quantity
		ExpectedErr error
	}{
		{
			Name: "correct quantity",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "g",
			},
			ExpectedErr: nil,
		},
		{
			Name: "negative value",
			Quantity: v1.Quantity{
				Value: -12,
				Unit:  "g",
			},
			ExpectedErr: ErrorQuantityValueInvalid,
		},
		{
			Name: "value zero",
			Quantity: v1.Quantity{
				Value: 0,
				Unit:  "g",
			},
			ExpectedErr: nil,
		},
		{
			Name: "empty unit",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "",
			},
			ExpectedErr: ErrorQuantityUnitMissing,
		},
		{
			Name: "blank unit",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "   ",
			},
			ExpectedErr: ErrorQuantityUnitMissing,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateQuantity(test.Quantity)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestValidateNutrition(t *testing.T) {
	correctNutrients := []v1.Nutrient{
		{
			T: "PROTEIN",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "g",
			},
		},
		{
			T: "CARBOHYDRATES",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "g",
			},
		},
		{
			T: "FAT",
			Quantity: v1.Quantity{
				Value: 12,
				Unit:  "g",
			},
		},
	}

	tests := []struct {
		Name        string
		Nutrition   v1.Nutrition
		ExpectedErr error
	}{
		{
			Name: "nutrition correct",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal:      123,
				Nutrients: correctNutrients,
			},
			ExpectedErr: nil,
		},
		{
			Name: "checks quantity",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: -12,
					Unit:  "g",
				},
				Kcal:      123,
				Nutrients: correctNutrients,
			},
			ExpectedErr: ErrorQuantityValueInvalid,
		},
		{
			Name: "kcal invalid",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal:      -123,
				Nutrients: correctNutrients,
			},
			ExpectedErr: ErrorNutritionKcalInvalid,
		},
		{
			Name: "kcal zero",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal:      0,
				Nutrients: correctNutrients,
			},
			ExpectedErr: nil,
		},
		{
			Name: "checks nutrients",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal: 123,
				Nutrients: []v1.Nutrient{
					{
						T: "CARBOHYDRATES",
						Quantity: v1.Quantity{
							Value: 12,
							Unit:  "g",
						},
					},
					{
						T: "FAT",
						Quantity: v1.Quantity{
							Value: 12,
							Unit:  "g",
						},
					},
				},
			},
			ExpectedErr: ErrorProteinMissing,
		},
		{
			Name: "checks minerals",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal:      123,
				Nutrients: correctNutrients,
				Minerals: []v1.Mineral{
					{
						T: "",
						Quantity: v1.Quantity{
							Value: 12,
							Unit:  "g",
						},
					},
				},
			},
			ExpectedErr: ErrorMineralTypeMissing,
		},
		{
			Name: "checks vitamins",
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: 12,
					Unit:  "g",
				},
				Kcal:      123,
				Nutrients: correctNutrients,
				Vitamins: []v1.Vitamin{
					{
						T: "",
						Quantity: v1.Quantity{
							Value: 12,
							Unit:  "g",
						},
					},
				},
			},
			ExpectedErr: ErrorVitaminTypeMissing,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateNutrition(test.Nutrition)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestValidateNutrients(t *testing.T) {
	tests := []struct {
		Name        string
		Nutrients   []v1.Nutrient
		ExpectedErr error
	}{
		{
			Name: "nutrients valid",
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: nil,
		},
		{
			Name: "protein missing",
			Nutrients: []v1.Nutrient{
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorProteinMissing,
		},
		{
			Name: "carbohydrates missing",
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorCarbohydratesMissing,
		},
		{
			Name: "fat missing",
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorFatMissing,
		},
		{
			Name: "nutrient type empty",
			Nutrients: []v1.Nutrient{
				{
					T: "",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorNutrientTypeMissing,
		},
		{
			Name: "nutrient type blank",
			Nutrients: []v1.Nutrient{
				{
					T: "   ",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorNutrientTypeMissing,
		},
		{
			Name: "checks nutrient quantity",
			Nutrients: []v1.Nutrient{
				{
					T: "PROTEIN",
					Quantity: v1.Quantity{
						Value: -12,
						Unit:  "g",
					},
				},
				{
					T: "CARBOHYDRATES",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
				{
					T: "FAT",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorQuantityValueInvalid,
		},
		{
			Name:        "nutrients empty",
			Nutrients:   []v1.Nutrient{},
			ExpectedErr: ErrorFatMissing,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateNutrients(test.Nutrients)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestValidateVitamins(t *testing.T) {
	tests := []struct {
		Name        string
		Vitamins    []v1.Vitamin
		ExpectedErr error
	}{
		{
			Name: "vitamins correct",
			Vitamins: []v1.Vitamin{
				{
					T: "VITAMIN_A",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:        "vitamins empty",
			Vitamins:    []v1.Vitamin{},
			ExpectedErr: nil,
		},
		{
			Name: "vitamin type empty",
			Vitamins: []v1.Vitamin{
				{
					T: "",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorVitaminTypeMissing,
		},
		{
			Name: "vitamin type blank",
			Vitamins: []v1.Vitamin{
				{
					T: "    ",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorVitaminTypeMissing,
		},
		{
			Name: "checks vitamin quantity",
			Vitamins: []v1.Vitamin{
				{
					T: "VITAMIN_A",
					Quantity: v1.Quantity{
						Value: -12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorQuantityValueInvalid,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateVitamins(test.Vitamins)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}

func TestValidateMinerals(t *testing.T) {
	tests := []struct {
		Name        string
		Minerals    []v1.Mineral
		ExpectedErr error
	}{
		{
			Name: "minerals correct",
			Minerals: []v1.Mineral{
				{
					T: "CALCIUM",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: nil,
		},
		{
			Name:        "minerals empty",
			Minerals:    []v1.Mineral{},
			ExpectedErr: nil,
		},
		{
			Name: "mineral type empty",
			Minerals: []v1.Mineral{
				{
					T: "",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorMineralTypeMissing,
		},
		{
			Name: "mineral type blank",
			Minerals: []v1.Mineral{
				{
					T: "   ",
					Quantity: v1.Quantity{
						Value: 12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorMineralTypeMissing,
		},
		{
			Name: "checks mineral quantity",
			Minerals: []v1.Mineral{
				{
					T: "IRON",
					Quantity: v1.Quantity{
						Value: -12,
						Unit:  "g",
					},
				},
			},
			ExpectedErr: ErrorQuantityValueInvalid,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			err := validateMinerals(test.Minerals)
			assert.Equal(t, test.ExpectedErr, err)
		})
	}
}
