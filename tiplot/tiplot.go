package tiplot

import (
	"image/color"
	"math"
	"strconv"

	"github.com/elordis/ti_tools/ti"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

const FourtyFiveDegree = math.Pi / 4

//nolint:gomnd
func DriveGlyph(d *ti.DriveTemplate) draw.GlyphStyle { //nolint:gocyclo,maintidx,gocognit
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

			switch name {
			case "ApexSolidRocket":
				shape = draw.PlusGlyph{}
			case "MeteorLiquidRocket":
				shape = draw.SquareGlyph{}
			case "NeutronLiquidRocket":
				shape = draw.RingGlyph{}
			case "VentureLiquidRocket":
				shape = draw.TriangleGlyph{}
			case "DianaSuperheavyRocket":
				shape = draw.BoxGlyph{}
			case "NovaLiquidRocket":
				shape = draw.CircleGlyph{}
			case "SuperKronosLiquidRocket":
				shape = draw.PyramidGlyph{}
			}
		case ti.DCElectrothermal, ti.DCElectromagnetic, ti.DCElectrostatic:
			_color = color.RGBA{R: 74, G: 102, B: 107, A: 255}

			switch name {
			// EM
			case "Resistojet":
				shape = draw.CircleGlyph{}
			case "TungstenResistojet":
				shape = draw.CircleGlyph{}
			case "ArcjetDrive":
				shape = draw.CircleGlyph{}
			case "E-BeamDrive":
				shape = draw.CircleGlyph{}
			case "AmplitronDrive":
				shape = draw.CircleGlyph{}
			// ET
			case "PlasmaWaveDrive":
				shape = draw.PyramidGlyph{}
			case "LorentzDrive":
				shape = draw.PyramidGlyph{}
			case "HeliconDrive":
				shape = draw.PyramidGlyph{}
			case "VASIMR":
				shape = draw.PyramidGlyph{}
			case "PonderomotiveVASIMR":
				shape = draw.PyramidGlyph{}
			case "PulsedPlasmoidDrive":
				shape = draw.PyramidGlyph{}
			case "MassDriver":
				shape = draw.PyramidGlyph{}
			case "SuperconductingMassDriver":
				shape = draw.PyramidGlyph{}
			// ES
			case "HallDrive":
				shape = draw.BoxGlyph{}
			case "IonDrive":
				shape = draw.BoxGlyph{}
			case "GridDrive":
				shape = draw.BoxGlyph{}
			case "ColloidDrive":
				shape = draw.BoxGlyph{}
			}
		case ti.DCFissionPulse, ti.DCFusionPulse:
			_color = color.RGBA{R: 222, G: 57, B: 0, A: 255}

			switch name {
			case "Z-pinchMicrofissionDrive":
				shape = draw.PlusGlyph{}
			case "NeutroniumMicrofissionDrive":
				shape = draw.SquareGlyph{}
			case "AntimatterMicrofissionDrive":
				shape = draw.BoxGlyph{}
			case "MinimagOrion":
				shape = draw.TriangleGlyph{}
			case "AdvancedMinimagOrion":
				shape = draw.PyramidGlyph{}
			case "OrionDrive":
				shape = draw.RingGlyph{}
			case "AdvancedOrionDrive":
				shape = draw.CircleGlyph{}
			}
		}
	case ti.PPCSolidCoreFission:
		_color = color.RGBA{R: 38, G: 62, B: 66, A: 255}

		switch name {
		case "KiwiDrive":
			shape = draw.CrossGlyph{}
		case "NervaDrive":
			shape = draw.CrossGlyph{}
		case "SnareDrive":
			shape = draw.CrossGlyph{}
		case "RoverDrive":
			shape = draw.CrossGlyph{}
		case "CermetNerva":
			shape = draw.CrossGlyph{}
		case "AdvancedNervaDrive":
			shape = draw.PlusGlyph{}
		case "Dumbo":
			shape = draw.RingGlyph{}
		case "AdvancedCermetNerva":
			shape = draw.PlusGlyph{}
		case "HeavyDumbo":
			shape = draw.CircleGlyph{}
		case "PulsarDrive":
			shape = draw.SquareGlyph{}
		case "AdvancedPulsarDrive":
			shape = draw.BoxGlyph{}
		case "PebbleDrive":
			shape = draw.TriangleGlyph{}
		case "AdvancedPebbleDrive":
			shape = draw.PyramidGlyph{}
		}
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

		switch name {
		case "VortexDrive":
			shape = draw.CrossGlyph{}
		case "AdvancedVortexDrive":
			shape = draw.PlusGlyph{}
		case "CavityDrive":
			shape = draw.CrossGlyph{}
		case "AdvancedCavityDrive":
			shape = draw.SquareGlyph{}
		case "QuartzDrive":
			shape = draw.CrossGlyph{}
		case "LightbulbDrive":
			shape = draw.CrossGlyph{}
		case "PharosDrive":
			shape = draw.RingGlyph{}
		case "FissionLantern":
			shape = draw.BoxGlyph{}
		case "FissionFragDrive":
			shape = draw.CrossGlyph{}
		case "DustyPlasmaDrive":
			shape = draw.CrossGlyph{}
		case "BurnerDrive":
			shape = draw.RingGlyph{}
		case "FlareDrive":
			shape = draw.PyramidGlyph{}
		case "FirestarDrive":
			shape = draw.CircleGlyph{}
		}
	case ti.PPCSaltWaterCore:
		_color = color.RGBA{R: 74, G: 167, B: 239, A: 255}

		switch name {
		case "NeutronFluxDrive":
			shape = draw.TriangleGlyph{}
		case "NeutronFluxTorch":
			shape = draw.PyramidGlyph{}
		}
	case ti.PPCZPinchFusion:
		_color = color.RGBA{R: 173, G: 211, B: 220, A: 255}

		switch name {
		case "TritonPulseDrive":
			shape = draw.RingGlyph{}
		case "FireflyTorch":
			shape = draw.SquareGlyph{}
		case "ZetaHelionDrive":
			shape = draw.CircleGlyph{}
		case "ZetaBoronFusionDrive":
			shape = draw.PyramidGlyph{}
		}
	case ti.PPCAnyMagneticConfinementFusion, ti.PPCElectrostaticConfinementFusion:
		_color = color.RGBA{R: 93, G: 80, B: 123, A: 255}

		switch name {
		case "LithiumFusionLantern":
			shape = draw.PyramidGlyph{}
		case "MagProtiumFusionDrive":
			shape = draw.CircleGlyph{}
		}
	case ti.PPCHybridConfinementFusion:
		_color = color.RGBA{R: 16, G: 17, B: 96, A: 255}

		switch name {
		case "HybridFusionDrive":
			shape = draw.RingGlyph{}
		case "IcarusDrive":
			shape = draw.TriangleGlyph{}
		case "IcarusTorch":
			shape = draw.PyramidGlyph{}
		case "AlienFusionLantern":
			shape = draw.SquareGlyph{}
		case "AlienFusionTorch":
			shape = draw.BoxGlyph{}
		}
	case ti.PPCToroidMagneticConfinementFusion:
		_color = color.RGBA{R: 170, G: 17, B: 95, A: 255}

		switch name {
		case "TritonTorusDrive":
			shape = draw.SquareGlyph{}
		case "HelionTorusDrive":
			shape = draw.RingGlyph{}
		case "AdvancedHelionTorusDrive":
			shape = draw.CircleGlyph{}
		}
	case ti.PPCMirroredMagneticConfinementFusion:
		_color = color.RGBA{R: 24, G: 116, B: 11, A: 255}

		switch name {
		case "TritonReflexDrive":
			shape = draw.SquareGlyph{}
		case "HelionReflexDrive":
			shape = draw.RingGlyph{}
		case "AdvancedHelionReflexDrive":
			shape = draw.CircleGlyph{}
		}
	case ti.PPCInertialConfinementFusion:
		_color = color.RGBA{R: 121, G: 155, B: 236, A: 255}

		switch name {
		case "TritonHopeInertialDrive":
			shape = draw.CrossGlyph{}
		case "TritonVistaInertialDrive":
			shape = draw.RingGlyph{}
		case "HelionInertialDrive":
			shape = draw.CircleGlyph{}
		case "DaedelusTorch":
			shape = draw.PlusGlyph{}
		case "BoronInertialDrive":
			shape = draw.TriangleGlyph{}
		case "BoronInertialTorch":
			shape = draw.PyramidGlyph{}
		case "ProtiumInertialTorch":
			shape = draw.SquareGlyph{}
		case "ProtiumConverterTorch":
			shape = draw.BoxGlyph{}
		}
	case ti.PPCAntimatterPlasmaCore, ti.PPCAntimatterBeamCore:
		_color = color.RGBA{R: 101, G: 199, B: 211, A: 255}

		switch name {
		case "AntimatterPulsedPlasmaCoreDrive":
			shape = draw.SquareGlyph{}
		case "AntimatterPlasmaCoreDrive":
			shape = draw.RingGlyph{}
		case "AdvancedAntimatterPlasmaCoreDrive":
			shape = draw.CircleGlyph{}
		case "PionTorch":
			shape = draw.PyramidGlyph{}
		}
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
