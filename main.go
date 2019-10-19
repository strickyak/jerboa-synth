// +build main

// Simulate basic synthesizer modules.
//
// $ go run main.go -db=0 -r=48000 -v=0 |  paplay --rate=48000 --channels=1 --format=s16le --raw /dev/stdin
package main

import (
	. "github.com/strickyak/jerboa-synth"

	"bufio"
	"flag"
	"io"
	"log"
	"math"
	"os"
)

var _ = log.Printf

var flagGain = flag.Float64("db", 0, "output gain in db")

func main() {
	flag.Parse()

	w := bufio.NewWriterSize(os.Stdout, 1024)
	defer func() { w.Flush() }()

	L1 := &VFO{
		Knobs: Knobs{
			Name:        "L1",
			Gain:        1,
			Sensitivity: 1.1,
			Base:        10,
		},
	}

	L2 := &VFO{
		Knobs: Knobs{
			Name:        "L2",
			Gain:        1,
			Sensitivity: 1.2,
			Base:        13,
		},
	}

	L3 := &VFO{
		Knobs: Knobs{
			Name:        "L3",
			Gain:        1,
			Sensitivity: 1.3,
			Base:        3,
		},
	}

	V1 := &VFO{
		Knobs: Knobs{
			Name:        "V1",
			Gain:        1,
			Sensitivity: 1,
			Base:        555,
		},
	}

	V2 := &VFO{
		Knobs: Knobs{
			Name:        "V2",
			Gain:        1,
			Sensitivity: 1,
			Base:        858,
		},
	}

	V3 := &VFO{
		Knobs: Knobs{
			Name:        "V3",
			Gain:        0.2,
			Sensitivity: 1,
			Base:        1000,
		},
	}

	r := &Rack{
		Devs: []Device{
			L1, L2, L3,
			V1, V2, V3,
		},
		Conns: []Conn{
			&Z2A{L1, L2}, &Z2A{L2, L3}, &Z2A{L3, L1},

			&Z2A{L1, V1}, &Z2A{L2, V2}, &Z2A{L3, V3},
		},
	}

	gain := math.Pow(10, *flagGain/10) * 10000 // flagGain in db

	for i := 0; i < 1000000000; i++ {
		// log.Printf("############# %d", i)
		r.Step()
		mix := V1.GetZ() + V2.GetZ() // + V3.GetZ()
		norm := mix / 3              // -1 to 1
		// log.Printf("############# norm = %10.6f", norm)
		if *FlagV {
			log.Printf("############# [%5d] %10.6f %10.6f %10.6f %10.6f %10.6f %10.6f", i, L1.GetZ(), L2.GetZ(), L3.GetZ(), V1.GetZ(), V2.GetZ(), V3.GetZ())
		}

		word := int16(norm * gain)
		Put16LE(w, word)
	}
}

func Put16LE(w io.Writer, word int16) {
	u := uint16(word)
	lo := byte(u & 0xFF)
	hi := byte((u >> 8) & 0xFF)
	cc, err := w.Write([]byte{lo, hi})
	if err != nil {
		log.Panicf("Put16LE: cannot write: %v", err)
	}
	if cc != 2 {
		log.Panicf("Put16LE: short write: %d", cc)
	}
}
