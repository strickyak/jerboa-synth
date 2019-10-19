package jerboa

import (
	"flag"
	"log"
	"math"
)

var _ = log.Printf

const VCC = 5.0
const GND = 0.0

const MIN = -2.5
const MAX = 2.5

func Clamp(a float64) float64 {
	switch {
	case a < MIN:
		return MIN
	case a > MAX:
		return MAX
	default:
		return a
	}
}

type Rack struct {
	Devs  []Device
	Conns []Conn
}

func (o *Rack) Step() {
	for _, e := range o.Devs {
		e.Step()
	}
	for _, e := range o.Conns {
		e.Step()
	}
}

type Conn interface {
	Step()
}

type Z2A struct {
	FromZ Device
	ToA   Device
}

func (o *Z2A) Step() {
	z := o.FromZ.GetZ()
	o.ToA.SetA(z)
	// log.Printf("Z2A: Z%q to A%q: %10.6f\n", o.FromZ, o.ToA, z)
}

type Device interface {
	OutBias(float64)
	OutGain(float64)

	InBase(float64)
	InSensitivity(float64)

	SetA(float64)
	SetB(float64)
	GetZ() float64

	Step()
	String() string
}

type Knobs struct {
	Name        string
	Bias        float64
	Gain        float64
	Base        float64
	Sensitivity float64

	A float64
	B float64
	Z float64
}

func (o *Knobs) OutBias(a float64) {
	o.Bias = a
}
func (o *Knobs) OutGain(a float64) {
	o.Gain = a
}

func (o *Knobs) InBase(a float64) {
	o.Base = a
}
func (o *Knobs) InSensitivity(a float64) {
	o.Sensitivity = a
}

func (o *Knobs) SetA(a float64) {
	o.A = Clamp(a)
}
func (o *Knobs) SetB(a float64) {
	o.B = Clamp(a)
}
func (o Knobs) GetZ() float64 {
	return Clamp(o.Z)
}
func (o Knobs) String() string {
	return o.Name
}

type VFO struct {
	Knobs
	Theta float64
}

var FlagRate = flag.Float64("r", 8000, "samples per second")
var FlagV = flag.Bool("v", false, "verbose")

func (o *VFO) Step() {
	freq := o.Base * math.Pow(2, o.A*o.Sensitivity)
	deltaTime := 1 / *FlagRate
	deltaPhase := deltaTime * freq

	o.Theta += deltaPhase
	if o.Theta > math.Pi {
		o.Theta -= 2 * math.Pi
	}

	o.Z = o.Bias + o.Gain*math.Sin(o.Theta)
	if *FlagV {
		log.Printf("## VFO.Step: %q: A %8.6f freq %8.6f dT %8.6f dP %8.6f T %8.6f Z %8.6f", o.Name, o.A, freq, deltaTime, deltaPhase, o.Theta, o.Z)
	}
}
