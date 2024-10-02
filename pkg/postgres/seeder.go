package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Seeder struct {
	pool   *pgxpool.Pool
	schema string
	seed   string
}

func NewSeeder(pool *pgxpool.Pool, schema, seed string) Seeder {
	return Seeder{pool: pool, schema: schema, seed: seed}
}

func (s Seeder) CreateSchema(ctx context.Context) error {
	_, execErr := s.pool.Exec(ctx, s.schema)
	if execErr != nil {
		return execErr
	}

	return nil
}

func (s Seeder) Seed(ctx context.Context) error {
	_, execErr := s.pool.Exec(ctx, s.seed)
	if execErr != nil {
		return execErr
	}

	return nil
}
