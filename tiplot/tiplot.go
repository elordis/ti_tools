package tiplot

import (
	"image/color"

	"github.com/elordis/ti_drive_plot/ti"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

//nolint:gomnd
func DriveStyle(d *ti.DriveTemplate) (draw.LineStyle, draw.GlyphStyle) {
	l := draw.LineStyle{ //nolint:exhaustruct
		Color: color.Black,
		Width: 1,
	}

	name, number := d.DataName[:len(d.DataName)-2], d.DataName[len(d.DataName)-2:]
	switch number {
	case "x1":
		l.Dashes = []vg.Length{18}
	case "x2":
		l.Dashes = []vg.Length{15}
	case "x3":
		l.Dashes = []vg.Length{12}
	case "x4":
		l.Dashes = []vg.Length{9}
	case "x5":
		l.Dashes = []vg.Length{6}
	case "x6":
		l.Dashes = []vg.Length{3}
	}

	g := draw.GlyphStyle{
		Color:  color.Black,
		Radius: 4, //nolint:gomnd
		Shape:  draw.CrossGlyph{},
	}

	var _color color.RGBA

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
			g.Shape = draw.SquareGlyph{}
		case "FissionSpinnerDrive":
			g.Shape = draw.TriangleGlyph{}
		case "PegasusDrive":
			g.Shape = draw.RingGlyph{}
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

	l.Color = _color
	g.Color = _color

	return l, g
}
