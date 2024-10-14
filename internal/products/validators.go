package products

import (
	"errors"
	"github.com/Kobietka/product-service/internal/ean"
	v1 "github.com/Kobietka/product-service/pkg/api/v1"
	"github.com/Kobietka/product-service/pkg/text"
)

var (
	ErrorProductEanMissing  = errors.New("PRODUCT_EAN_MISSING")
	ErrorProductEanInvalid  = errors.New("PRODUCT_EAN_INVALID")
	ErrorProductNameMissing = errors.New("PRODUCT_NAME_MISSING")
)

func validateProduct(product v1.Product) error {
	if text.IsBlankString(product.Ean) {
		return ErrorProductEanMissing
	}

	if !ean.IsValid(product.Ean) {
		return ErrorProductEanInvalid
	}

	if text.IsBlankString(product.Name) {
		return ErrorProductNameMissing
	}

	err := validateQuantity(product.Packaging)
	if err != nil {
		return err
	}

	return validateNutrition(product.Nutrition)
}

var (
	ErrorNutritionKcalInvalid = errors.New("NUTRITION_KCAL_INVALID")
)

func validateNutrition(nutrition v1.Nutrition) error {
	err := validateQuantity(nutrition.Per)
	if err != nil {
		return err
	}

	if nutrition.Kcal < 0 {
		return ErrorNutritionKcalInvalid
	}

	err = validateNutrients(nutrition.Nutrients)
	if err != nil {
		return err
	}

	err = validateVitamins(nutrition.Vitamins)
	if err != nil {
		return err
	}

	err = validateMinerals(nutrition.Minerals)
	if err != nil {
		return err
	}

	return nil
}

var (
	ErrorQuantityUnitMissing  = errors.New("QUANTITY_UNIT_MISSING")
	ErrorQuantityValueInvalid = errors.New("QUANTITY_VALUE_INVALID")
)

func validateQuantity(quantity v1.Quantity) error {
	if text.IsBlankString(quantity.Unit) {
		return ErrorQuantityUnitMissing
	}

	if quantity.Value < 0 {
		return ErrorQuantityValueInvalid
	}

	return nil
}

const (
	FatType           = "FAT"
	CarbohydratesType = "CARBOHYDRATES"
	ProteinType       = "PROTEIN"
)

var (
	ErrorNutrientTypeMissing  = errors.New("NUTRIENT_TYPE_MISSING")
	ErrorFatMissing           = errors.New("NUTRIENT_FAT_MISSING")
	ErrorCarbohydratesMissing = errors.New("NUTRIENT_CARBOHYDRATES_MISSING")
	ErrorProteinMissing       = errors.New("NUTRIENT_PROTEIN_MISSING")
)

func validateNutrients(nutrients []v1.Nutrient) error {
	containsFat := false
	containsCarbohydrates := false
	containsProteins := false
	for _, item := range nutrients {
		if text.IsBlankString(item.T) {
			return ErrorNutrientTypeMissing
		}

		err := validateQuantity(item.Quantity)
		if err != nil {
			return err
		}

		switch item.T {
		case FatType:
			containsFat = true
		case CarbohydratesType:
			containsCarbohydrates = true
		case ProteinType:
			containsProteins = true
		}
	}

	if !containsFat {
		return ErrorFatMissing
	}
	if !containsCarbohydrates {
		return ErrorCarbohydratesMissing
	}
	if !containsProteins {
		return ErrorProteinMissing
	}

	return nil
}

var (
	ErrorVitaminTypeMissing = errors.New("VITAMIN_TYPE_MISSING")
)

func validateVitamins(vitamins []v1.Vitamin) error {
	for _, vitamin := range vitamins {
		if text.IsBlankString(vitamin.T) {
			return ErrorVitaminTypeMissing
		}

		err := validateQuantity(vitamin.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}

var (
	ErrorMineralTypeMissing = errors.New("MINERAL_TYPE_MISSING")
)

func validateMinerals(minerals []v1.Mineral) error {
	for _, mineral := range minerals {
		if text.IsBlankString(mineral.T) {
			return ErrorMineralTypeMissing
		}

		err := validateQuantity(mineral.Quantity)
		if err != nil {
			return err
		}
	}
	return nil
}
