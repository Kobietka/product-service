package database

import "github.com/jackc/pgx/v5"

func addProductQueries(batch *pgx.Batch, ean string) {
	productQuery := `
		SELECT ean, name 
		FROM product
		WHERE ean = $1
	`
	batch.Queue(productQuery, ean)

	packagingQuery := `
		SELECT
		packaging.ean,
		packaging.value,
		unit.value
		FROM packaging
		JOIN unit ON unit.id = unit_id
		WHERE packaging.ean = $1
	`
	batch.Queue(packagingQuery, ean)

	nutritionQuery := `
		SELECT
		ean,
		kcal
		FROM nutrition
		WHERE ean = $1
    `
	batch.Queue(nutritionQuery, ean)

	nutritionQuantityQuery := `
		SELECT
		nutrition_quantity.ean,
		nutrition_quantity.value,
		unit.value
		FROM nutrition_quantity
		JOIN unit ON unit.id = unit_id
		WHERE nutrition_quantity.ean = $1
	`
	batch.Queue(nutritionQuantityQuery, ean)

	nutrientsQuery := `
		SELECT
		nutrient.ean,
		nutrient_type.type,
		nutrient.value,
		unit.value
		FROM nutrient
		JOIN unit ON unit.id = unit_id
		JOIN nutrient_type ON nutrient.type_id = nutrient_type.id
		WHERE nutrient.ean = $1
    `
	batch.Queue(nutrientsQuery, ean)

	vitaminsQuery := `
		SELECT
		vitamin.ean,
		vitamin_type.type,
		vitamin.value,
		unit.value
		FROM vitamin
		JOIN unit ON unit.id = unit_id
		JOIN vitamin_type ON vitamin.type_id = vitamin_type.id
		WHERE vitamin.ean = $1
	`
	batch.Queue(vitaminsQuery, ean)

	mineralsQuery := `
		SELECT
		mineral.ean,
		mineral_type.type,
		mineral.value,
		unit.value
		FROM mineral
		JOIN unit ON unit.id = unit_id
		JOIN mineral_type ON mineral.type_id = mineral_type.id
		WHERE mineral.ean = $1
	`
	batch.Queue(mineralsQuery, ean)
}
