package text

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsBlankString(t *testing.T) {
	tests := []struct {
		Name    string
		Value   string
		IsBlank bool
	}{
		{
			Name:    "correct",
			Value:   "value",
			IsBlank: false,
		},
		{
			Name:    "blank",
			Value:   "    ",
			IsBlank: true,
		},
		{
			Name:    "empty",
			Value:   "",
			IsBlank: true,
		},
	}
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			isBlank := IsBlankString(test.Value)
			assert.Equal(t, test.IsBlank, isBlank)
		})
	}
}
