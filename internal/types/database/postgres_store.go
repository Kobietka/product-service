package database

import (
	"context"
	"errors"
	"github.com/Kobietka/product-service/pkg/array"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresStore struct {
	pool *pgxpool.Pool
}

func NewPostgresStore(pool *pgxpool.Pool) PostgresStore {
	return PostgresStore{pool}
}

func (s PostgresStore) GetUnits(ctx context.Context) ([]string, error) {
	query := `SELECT id, value FROM unit`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[unitEntity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return array.MapArray(entities, func(entity unitEntity) string {
		return entity.Value
	}), nil
}

func (s PostgresStore) GetNutrientTypes(ctx context.Context) ([]string, error) {
	query := `SELECT id, type FROM nutrient_type`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[nutrientTypeEntity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return array.MapArray(entities, func(entity nutrientTypeEntity) string {
		return entity.Type
	}), nil
}

func (s PostgresStore) GetVitaminTypes(ctx context.Context) ([]string, error) {
	query := `SELECT id, type FROM vitamin_type`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[vitaminTypeEntity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return array.MapArray(entities, func(entity vitaminTypeEntity) string {
		return entity.Type
	}), nil
}

func (s PostgresStore) GetMineralTypes(ctx context.Context) ([]string, error) {
	query := `SELECT id, type FROM mineral_type`

	rows, err := s.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities, err := pgx.CollectRows(rows, pgx.RowToStructByPos[mineralTypeEntity])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return []string{}, nil
		}
		return nil, err
	}

	return array.MapArray(entities, func(entity mineralTypeEntity) string {
		return entity.Type
	}), nil
}
