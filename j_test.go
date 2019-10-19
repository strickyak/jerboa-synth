package jerboa

import (
	"log"
	"testing"
)

// This is hardly a test.
func TestOne(t *testing.T) {
	vfo := &VFO{
		Knobs: Knobs{
			Gain: 1,
			A:    5,
		},
	}
	*FlagRate = 10
	for i := 0; i < 100; i++ {
		vfo.Step()
		log.Printf("## %4d: %10.6f\n", i, vfo.Z)
	}
}
