package tiplot

import (
	"image/color"
	"math"
	"strconv"

	"github.com/elordis/ti_drive_plot/ti"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const FourtyFiveDegree = math.Pi / 4

//nolint:gomnd
func DriveGlyph(d *ti.DriveTemplate) draw.GlyphStyle {
	name := d.DataName[:len(d.DataName)-2]
	_color := color.RGBA{R: 0, G: 0, B: 0, A: 255}
	radius := vg.Length(4) //nolint:gomnd
	shape := draw.GlyphDrawer(draw.CrossGlyph{})

	// SquareGlyph
	// TriangleGlyph
	// RingGlyph
	// BoxGlyph
	// PyramidGlyph
	// CircleGlyph
	// PlusGlyph
	// CrossGlyph
	switch d.RequiredPowerPlant { //nolint:exhaustive
	case ti.PPCAnyGeneral:
		switch d.DriveClassification { //nolint:exhaustive
		case ti.DCChemical:
			_color = color.RGBA{R: 55, G: 147, B: 67, A: 255}
		case ti.DCElectrothermal, ti.DCElectromagnetic, ti.DCElectrostatic:
			_color = color.RGBA{R: 74, G: 102, B: 107, A: 255}
		case ti.DCFissionPulse, ti.DCFusionPulse:
			_color = color.RGBA{R: 222, G: 57, B: 0, A: 255}
		}
	case ti.PPCSolidCoreFission:
		_color = color.RGBA{R: 38, G: 62, B: 66, A: 255}
	case ti.PPCLiquidCoreFission:
		_color = color.RGBA{R: 46, G: 147, B: 88, A: 255}

		switch name {
		case "LarsDrive":
			shape = draw.SquareGlyph{}
		case "FissionSpinnerDrive":
			shape = draw.TriangleGlyph{}
		case "PegasusDrive":
			shape = draw.RingGlyph{}
		}
	case ti.PPCGasCoreFission:
		_color = color.RGBA{R: 137, G: 176, B: 38, A: 255}
	case ti.PPCSaltWaterCore:
		_color = color.RGBA{R: 74, G: 167, B: 239, A: 255}
	case ti.PPCZPinchFusion:
		_color = color.RGBA{R: 173, G: 211, B: 220, A: 255}
	case ti.PPCAnyMagneticConfinementFusion, ti.PPCElectrostaticConfinementFusion:
		_color = color.RGBA{R: 93, G: 80, B: 123, A: 255}
	case ti.PPCHybridConfinementFusion:
		_color = color.RGBA{R: 16, G: 17, B: 96, A: 255}
	case ti.PPCToroidMagneticConfinementFusion:
		_color = color.RGBA{R: 170, G: 17, B: 95, A: 255}
	case ti.PPCMirroredMagneticConfinementFusion:
		_color = color.RGBA{R: 24, G: 116, B: 11, A: 255}
	case ti.PPCInertialConfinementFusion:
		_color = color.RGBA{R: 121, G: 155, B: 236, A: 255}
	case ti.PPCAntimatterPlasmaCore, ti.PPCAntimatterBeamCore:
		_color = color.RGBA{R: 101, G: 199, B: 211, A: 255}
	}

	return draw.GlyphStyle{
		Color:  _color,
		Radius: radius,
		Shape:  shape,
	}
}

type MajorLogTicks struct {
}

func (t MajorLogTicks) Ticks(min, max float64) []plot.Tick {
	if min <= 0 || max <= 0 {
		panic("Values must be greater than 0 for a log scale.")
	}

	var ticks []plot.Tick

	val := math.Pow10(int(math.Log10(min)))
	max = math.Pow10(int(math.Ceil(math.Log10(max))))

	for val < max {
		for i := 1; i < 10; i++ {
			v := val * float64(i)
			tick := plot.Tick{
				Value: v,
				Label: strconv.FormatFloat(v, 'f', -1, 64),
			}
			ticks = append(ticks, tick)
		}

		val *= 10
	}

	ticks = append(ticks, plot.Tick{Value: val, Label: strconv.FormatFloat(val, 'f', -1, 64)})

	return ticks
}
