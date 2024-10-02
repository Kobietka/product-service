package types

import "context"

type Store interface {
	GetUnits(ctx context.Context) ([]string, error)
	GetNutrientTypes(ctx context.Context) ([]string, error)
	GetVitaminTypes(ctx context.Context) ([]string, error)
	GetMineralTypes(ctx context.Context) ([]string, error)
}
