package postgres

import (
	"github.com/jackc/pgx/v5"
)

func CollectOneRow[T any](results pgx.BatchResults) (T, error) {
	rows, queryErr := results.Query()
	if queryErr != nil {
		return *new(T), queryErr
	}
	defer rows.Close()

	data, collectErr := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[T])
	if collectErr != nil {
		return *new(T), collectErr
	}

	return data, nil
}

func CollectRows[T any](results pgx.BatchResults) ([]T, error) {
	rows, queryErr := results.Query()
	if queryErr != nil {
		return make([]T, 0), queryErr
	}
	defer rows.Close()

	data, collectErr := pgx.CollectRows(rows, pgx.RowToStructByPos[T])
	if collectErr != nil {
		return make([]T, 0), collectErr
	}

	return data, nil
}
