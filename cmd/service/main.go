package main

import (
	"context"
	"fmt"
	"github.com/charmbracelet/log"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kobietka/product-service/internal/config"
	dbsetup "github.com/kobietka/product-service/internal/database/setup"
	"github.com/kobietka/product-service/internal/products"
	productdb "github.com/kobietka/product-service/internal/products/database"
	"github.com/kobietka/product-service/internal/types"
	typesdb "github.com/kobietka/product-service/internal/types/database"
	"github.com/kobietka/product-service/pkg/logger"
	"github.com/labstack/echo/v4"
)

func main() {
	configStore := config.NewConfigStore()

	c, err := configStore.GetConfig()
	if err != nil {
		panic(err)
	}

	poolConfig, err := pgxpool.ParseConfig(c.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		panic(err)
	}

	seeder := dbsetup.NewSeeder(pool)
	err = seeder.CreateSchema(context.Background())
	if err != nil {
		panic(err)
	}
	err = seeder.Seed(context.Background())
	if err != nil {
		panic(err)
	}

	productStore := productdb.NewPostgresStore(pool)
	unitStore := typesdb.NewPostgresStore(pool)
	productServer := products.NewServer(productStore)
	typeServer := types.NewServer(unitStore)

	e := echo.New()
	e.Use(logger.NewBasicRequestLogger())

	productServer.Routes(e)
	typeServer.Routes(e)

	log.Fatal(e.Start(fmt.Sprintf(":%s", c.Port)))
}
