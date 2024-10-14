package setup

import (
	_ "embed"
	"github.com/Kobietka/product-service/pkg/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed schema.sql
var schema string

//go:embed seed.sql
var seed string

func NewSeeder(pool *pgxpool.Pool) postgres.Seeder {
	return postgres.NewSeeder(pool, schema, seed)
}
