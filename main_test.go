package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestComputeCombination(t *testing.T) {
	tcs := []struct {
		input, solution Word
		expected CombinationArray
	}{
		{"AAAAA", "AABAA", CombinationArray{2, 2, 0, 2, 2}},
		{"CAAAA", "AABAA", CombinationArray{0, 2, 1, 2, 2}},
	}

	for _, tc := range tcs {
		res := computeCombination(tc.input, tc.solution)
		assert.Equal(t, tc.expected, res)
	}
}

