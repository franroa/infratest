package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	t.Parallel()

	min := 0
	max := 100

	for i := 0; i < 100000; i++ {
		value := Random(min, max)
		assert.True(t, value >= min && value <= max)
	}
}

func TestUniqueId(t *testing.T) {
	t.Parallel()

	previouslySeen := map[string]bool{}

	for i := 0; i < 100; i++ {
		uniqueID := UniqueId()
		assert.Len(t, uniqueID, 6)
		assert.NotContains(t, previouslySeen, uniqueID)

		previouslySeen[uniqueID] = true
	}
}
