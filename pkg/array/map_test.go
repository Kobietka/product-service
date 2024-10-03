package array

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMapArray(t *testing.T) {
	type InputStruct struct {
		Value int
	}

	tests := []struct {
		Name           string
		Input          []InputStruct
		ExpectedOutput []int
	}{
		{
			Name: "maps correctly",
			Input: []InputStruct{
				{Value: 1},
				{Value: 2},
				{Value: 3},
			},
			ExpectedOutput: []int{1, 2, 3},
		},
		{
			Name:           "empty input",
			Input:          []InputStruct{},
			ExpectedOutput: []int{},
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			output := MapArray(test.Input, func(i InputStruct) int {
				return i.Value
			})
			assert.Equal(t, test.ExpectedOutput, output)
		})
	}
}
