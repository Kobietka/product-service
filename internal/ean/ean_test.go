package ean

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := []struct {
		Name  string
		Ean   string
		Valid bool
	}{
		{
			Name:  "ean-8",
			Ean:   "12345678",
			Valid: true,
		},
		{
			Name:  "ean-13",
			Ean:   "1234567890123",
			Valid: true,
		},
		{
			Name:  "upc-a code",
			Ean:   "123456789012",
			Valid: true,
		},
		{
			Name:  "empty",
			Ean:   "",
			Valid: false,
		},
		{
			Name:  "blank",
			Ean:   "    ",
			Valid: false,
		},
		{
			Name:  "invalid",
			Ean:   "1234",
			Valid: false,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			isValid := IsValid(test.Ean)
			assert.Equal(t, test.Valid, isValid)
		})
	}
}
