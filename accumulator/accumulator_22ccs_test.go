package accumulator

import (
	"math/big"
	"testing"
)

func testSetup(t *testing.T) {
	maxCapacity := big.NewInt(20000)
	accumulator := NewAccumulator(maxCapacity)
	accumulator.SetUp()

}
