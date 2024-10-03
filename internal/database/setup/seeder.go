package setup

import (
	_ "embed"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kobietka/product-service/pkg/postgres"
)

//go:embed schema.sql
var schema string

//go:embed seed.sql
var seed string

func NewSeeder(pool *pgxpool.Pool) postgres.Seeder {
	return postgres.NewSeeder(pool, schema, seed)
}
