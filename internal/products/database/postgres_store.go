package database

import (
	"context"
	"errors"
	v1 "github.com/Kobietka/product-service/pkg/api/v1"
	"github.com/Kobietka/product-service/pkg/array"
	"github.com/Kobietka/product-service/pkg/postgres"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) PostgresStore {
	return PostgresStore{pool: pool}
}

func (s PostgresStore) GetProduct(ctx context.Context, ean string) (v1.Product, error) {
	batch := pgx.Batch{}
	addProductQueries(&batch, ean)
	batchResults := s.pool.SendBatch(ctx, &batch)

	productE, err := postgres.CollectOneRow[productEntity](batchResults)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, v1.ErrorDataNotFound
		}
		return v1.Product{}, err
	}

	packagingE, err := postgres.CollectOneRow[packagingEntity](batchResults)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, v1.ErrorDataNotFound
		}
		return v1.Product{}, err
	}

	nutritionE, err := postgres.CollectOneRow[nutritionEntity](batchResults)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, v1.ErrorDataNotFound
		}
		return v1.Product{}, err
	}

	nutritionQuantityE, err := postgres.CollectOneRow[nutritionQuantityEntity](batchResults)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, v1.ErrorDataNotFound
		}
		return v1.Product{}, err
	}

	nutrientEntities, err := postgres.CollectRows[nutrientEntity](batchResults)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, err
		}
	}

	vitaminEntities, err := postgres.CollectRows[vitaminEntity](batchResults)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, err
		}
	}

	mineralEntities, err := postgres.CollectRows[mineralEntity](batchResults)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return v1.Product{}, err
		}
	}

	err = batchResults.Close()
	if err != nil {
		return v1.Product{}, err
	}

	return v1.Product{
		Ean:  productE.Ean,
		Name: productE.Name,
		Packaging: v1.Quantity{
			Value: packagingE.Value,
			Unit:  packagingE.Unit,
		},
		Nutrition: v1.Nutrition{
			Per: v1.Quantity{
				Value: nutritionQuantityE.Value,
				Unit:  nutritionQuantityE.Unit,
			},
			Kcal: nutritionE.Kcal,
			Nutrients: array.MapArray(nutrientEntities, func(entity nutrientEntity) v1.Nutrient {
				return v1.Nutrient{
					T: entity.Type,
					Quantity: v1.Quantity{
						Value: entity.Value,
						Unit:  entity.Unit,
					},
				}
			}),
			Vitamins: array.MapArray(vitaminEntities, func(entity vitaminEntity) v1.Vitamin {
				return v1.Vitamin{
					T: entity.Type,
					Quantity: v1.Quantity{
						Value: entity.Value,
						Unit:  entity.Unit,
					},
				}
			}),
			Minerals: array.MapArray(mineralEntities, func(entity mineralEntity) v1.Mineral {
				return v1.Mineral{
					T: entity.Type,
					Quantity: v1.Quantity{
						Value: entity.Value,
						Unit:  entity.Unit,
					},
				}
			}),
		},
	}, nil
}

func (s PostgresStore) SearchProducts(ctx context.Context, query string, limit int8) ([]v1.Product, error) {
	conn, err := s.pool.Acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	productsSearchQuery := `
		SELECT 
		ean, 
		name 
		FROM product 
		WHERE LOWER(name) LIKE LOWER($1) 
		LIMIT $2
	`
	rows, err := conn.Query(ctx, productsSearchQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, err
	}

	productEntities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[productEntity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return make([]v1.Product, 0), nil
		}
		return nil, err
	}

	batch := pgx.Batch{}
	for _, entity := range productEntities {
		addProductQueries(&batch, entity.Ean)
	}

	batchResults := conn.SendBatch(ctx, &batch)

	products := make([]v1.Product, 0)
	for range productEntities {
		productE, entityErr := postgres.CollectOneRow[productEntity](batchResults)
		packagingE, entityErr := postgres.CollectOneRow[packagingEntity](batchResults)
		nutritionE, entityErr := postgres.CollectOneRow[nutritionEntity](batchResults)
		nutritionQuantityE, entityErr := postgres.CollectOneRow[nutritionQuantityEntity](batchResults)
		nutrientEntities, entityErr := postgres.CollectRows[nutrientEntity](batchResults)
		vitaminEntities, entityErr := postgres.CollectRows[vitaminEntity](batchResults)
		mineralEntities, entityErr := postgres.CollectRows[mineralEntity](batchResults)

		if entityErr != nil {
			if !errors.Is(entityErr, pgx.ErrNoRows) {
				continue
			}
		}

		product := v1.Product{
			Ean:  productE.Ean,
			Name: productE.Name,
			Packaging: v1.Quantity{
				Value: packagingE.Value,
				Unit:  packagingE.Unit,
			},
			Nutrition: v1.Nutrition{
				Per: v1.Quantity{
					Value: nutritionQuantityE.Value,
					Unit:  nutritionQuantityE.Unit,
				},
				Kcal: nutritionE.Kcal,
				Nutrients: array.MapArray(nutrientEntities, func(entity nutrientEntity) v1.Nutrient {
					return v1.Nutrient{
						T: entity.Type,
						Quantity: v1.Quantity{
							Value: entity.Value,
							Unit:  entity.Unit,
						},
					}
				}),
				Vitamins: array.MapArray(vitaminEntities, func(entity vitaminEntity) v1.Vitamin {
					return v1.Vitamin{
						T: entity.Type,
						Quantity: v1.Quantity{
							Value: entity.Value,
							Unit:  entity.Unit,
						},
					}
				}),
				Minerals: array.MapArray(mineralEntities, func(entity mineralEntity) v1.Mineral {
					return v1.Mineral{
						T: entity.Type,
						Quantity: v1.Quantity{
							Value: entity.Value,
							Unit:  entity.Unit,
						},
					}
				}),
			},
		}
		products = append(products, product)
	}

	err = batchResults.Close()
	if err != nil {
		return nil, err
	}

	return products, nil
}

func (s PostgresStore) CreateProduct(ctx context.Context, product v1.Product) error {
	batch := pgx.Batch{}

	productQuery := `
		INSERT INTO product(ean, name)
		VALUES ($1, $2);
	`
	batch.Queue(productQuery, product.Ean, product.Name)

	packagingQuery := `
		INSERT INTO packaging(ean, value, unit_id)
		VALUES ($1, $2, (SELECT id FROM unit WHERE value = $3));
	`
	batch.Queue(packagingQuery, product.Ean, product.Packaging.Value, product.Packaging.Unit)

	nutritionQuery := `
		INSERT INTO nutrition(ean, kcal)
		VALUES ($1, $2);
	`
	batch.Queue(nutritionQuery, product.Ean, product.Nutrition.Kcal)

	nutritionQuantityQuery := `
		INSERT INTO nutrition_quantity(ean, value, unit_id)
		VALUES ($1, $2, (SELECT id FROM unit WHERE value = $3));
	`
	batch.Queue(nutritionQuantityQuery, product.Ean, product.Nutrition.Per.Value, product.Nutrition.Per.Unit)

	nutrientQuery := `
		INSERT INTO nutrient(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM nutrient_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, nutrient := range product.Nutrition.Nutrients {
		batch.Queue(nutrientQuery, product.Ean, nutrient.T, nutrient.Quantity.Value, nutrient.Quantity.Unit)
	}

	vitaminQuery := `
		INSERT INTO vitamin(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM vitamin_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, vitamin := range product.Nutrition.Vitamins {
		batch.Queue(vitaminQuery, product.Ean, vitamin.T, vitamin.Quantity.Value, vitamin.Quantity.Unit)
	}

	mineralQuery := `
		INSERT INTO mineral(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM mineral_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, mineral := range product.Nutrition.Minerals {
		batch.Queue(mineralQuery, product.Ean, mineral.T, mineral.Quantity.Value, mineral.Quantity.Unit)
	}

	err := s.pool.SendBatch(ctx, &batch).Close()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return v1.ErrorInvalidData
		}
		return err
	}

	return err
}

func (s PostgresStore) UpdateProduct(ctx context.Context, product v1.Product) error {
	batch := pgx.Batch{}

	productQuery := `
		UPDATE product
		SET name = $2
		WHERE ean = $1
	`
	batch.Queue(productQuery, product.Ean, product.Name)

	packagingQuery := `
		UPDATE packaging
		SET value = $2, unit_id = (SELECT id FROM unit WHERE value = $3)
		WHERE ean = $1
	`
	batch.Queue(packagingQuery, product.Ean, product.Packaging.Value, product.Packaging.Unit)

	nutritionQuery := `
		UPDATE nutrition
		SET kcal = $2
		WHERE ean = $1
	`
	batch.Queue(nutritionQuery, product.Ean, product.Nutrition.Kcal)

	nutritionQuantityQuery := `
		UPDATE nutrition_quantity
		SET value = $2, unit_id = (SELECT id FROM unit WHERE value = $3)
		WHERE ean = $1
	`
	batch.Queue(nutritionQuantityQuery, product.Ean, product.Nutrition.Per.Value, product.Nutrition.Per.Unit)

	deleteNutrientsQuery := `DELETE FROM nutrient WHERE ean = $1`
	batch.Queue(deleteNutrientsQuery, product.Ean)

	nutrientQuery := `
		INSERT INTO nutrient(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM nutrient_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, nutrient := range product.Nutrition.Nutrients {
		batch.Queue(nutrientQuery, product.Ean, nutrient.T, nutrient.Quantity.Value, nutrient.Quantity.Unit)
	}

	deleteVitaminsQuery := `DELETE FROM vitamin WHERE ean = $1`
	batch.Queue(deleteVitaminsQuery, product.Ean)

	vitaminQuery := `
		INSERT INTO vitamin(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM vitamin_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, vitamin := range product.Nutrition.Vitamins {
		batch.Queue(vitaminQuery, product.Ean, vitamin.T, vitamin.Quantity.Value, vitamin.Quantity.Unit)
	}

	deleteMineralsQuery := `DELETE FROM mineral WHERE ean = $1`
	batch.Queue(deleteMineralsQuery, product.Ean)

	mineralQuery := `
		INSERT INTO mineral(ean, type_id, value, unit_id)
		VALUES ($1, (SELECT id FROM mineral_type WHERE type = $2), $3, (SELECT id FROM unit WHERE value = $4));
	`
	for _, mineral := range product.Nutrition.Minerals {
		batch.Queue(mineralQuery, product.Ean, mineral.T, mineral.Quantity.Value, mineral.Quantity.Unit)
	}

	batchResults := s.pool.SendBatch(ctx, &batch)

	tag, err := batchResults.Exec()
	if err != nil {
		return err
	}

	if tag.RowsAffected() == 0 {
		return v1.ErrorProductDoesNotExist
	}

	err = batchResults.Close()
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			return v1.ErrorInvalidData
		}
		return err
	}

	return nil
}

func (s PostgresStore) DeleteProduct(ctx context.Context, ean string) error {
	query := `DELETE FROM product WHERE ean = $1`

	tag, execErr := s.pool.Exec(ctx, query, ean)
	if execErr != nil {
		return execErr
	}

	if tag.RowsAffected() == 0 {
		return v1.ErrorProductDoesNotExist
	}

	return nil
}
