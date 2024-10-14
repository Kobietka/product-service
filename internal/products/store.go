package products

import (
	"context"
	v1 "github.com/Kobietka/product-service/pkg/api/v1"
)

type Store interface {
	GetProduct(ctx context.Context, ean string) (v1.Product, error)
	SearchProducts(ctx context.Context, query string, limit int8) ([]v1.Product, error)
	CreateProduct(ctx context.Context, product v1.Product) error
	UpdateProduct(ctx context.Context, product v1.Product) error
	DeleteProduct(ctx context.Context, ean string) error
}
