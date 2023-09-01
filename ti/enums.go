package ti

import (
	"errors"
	"fmt"
	"strings"
)

var ErrBadValue = errors.New("bad value")

//go:generate stringer -type=PowerPlantClass
type PowerPlantClass int

const (
	PPCAnyGeneral PowerPlantClass = iota
	PPCFuelCell
	PPCSolidCoreFission
	PPCLiquidCoreFission
	PPCGasCoreFission
	PPCSaltWaterCore
	PPCZPinchFusion
	PPCAnyMagneticConfinementFusion
	PPCElectrostaticConfinementFusion
	PPCHybridConfinementFusion
	PPCToroidMagneticConfinementFusion
	PPCMirroredMagneticConfinementFusion
	PPCInertialConfinementFusion
	PPCAntimatterSolidCore
	PPCAntimatterGasCore
	PPCAntimatterPlasmaCore
	PPCAntimatterBeamCore
)

func (p *PowerPlantClass) UnmarshalJSON(data []byte) error {
	s := string(data)
	if !(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		return fmt.Errorf("trying to unmarshall unexpected PowerPlantClass value: %s: %w", data, ErrBadValue)
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
	switch s {
	case "Any_General":
		*p = PPCAnyGeneral
	case "Fuel_Cell":
		*p = PPCFuelCell
	case "Solid_Core_Fission":
		*p = PPCSolidCoreFission
	case "Liquid_Core_Fission":
		*p = PPCLiquidCoreFission
	case "Gas_Core_Fission":
		*p = PPCGasCoreFission
	case "Salt_Water_Core":
		*p = PPCSaltWaterCore
	case "Z_Pinch_Fusion":
		*p = PPCZPinchFusion
	case "Any_Magnetic_Confinement_Fusion":
		*p = PPCAnyMagneticConfinementFusion
	case "Electrostatic_Confinement_Fusion":
		*p = PPCElectrostaticConfinementFusion
	case "Hybrid_Confinement_Fusion":
		*p = PPCHybridConfinementFusion
	case "Toroid_Magnetic_Confinement_Fusion":
		*p = PPCToroidMagneticConfinementFusion
	case "Mirrored_Magnetic_Confinement_Fusion":
		*p = PPCMirroredMagneticConfinementFusion
	case "Inertial_Confinement_Fusion":
		*p = PPCInertialConfinementFusion
	case "Antimatter_Solid_Core":
		*p = PPCAntimatterSolidCore
	case "Antimatter_Gas_Core":
		*p = PPCAntimatterGasCore
	case "Antimatter_Plasma_Core":
		*p = PPCAntimatterPlasmaCore
	case "Antimatter_Beam_Core":
		*p = PPCAntimatterBeamCore
	default:
		return fmt.Errorf("trying to unmarshall unexpected PowerPlantClass value: %s: %w ", data, ErrBadValue)
	}

	return nil
}

//go:generate stringer -type=DriveClassification
type DriveClassification int

const (
	DCChemical DriveClassification = iota
	DCElectrothermal
	DCElectromagnetic
	DCElectrostatic
	DCFissionThermal
	DCFissionPulse
	DCFusionThermal
	DCFusionPulse
	DCAntimatter
)

func (d *DriveClassification) UnmarshalJSON(data []byte) error {
	s := string(data)
	if !(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		return fmt.Errorf("trying to unmarshall unexpected DriveClassification value: %s: %w", data, ErrBadValue)
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
	switch s {
	case "Chemical":
		*d = DCChemical
	case "Electrothermal":
		*d = DCElectrothermal
	case "Electromagnetic":
		*d = DCElectromagnetic
	case "Electrostatic":
		*d = DCElectrostatic
	case "Fission_Thermal":
		*d = DCFissionThermal
	case "Fission_Pulse":
		*d = DCFissionPulse
	case "Fusion_Thermal":
		*d = DCFusionThermal
	case "Fusion_Pulse":
		*d = DCFusionPulse
	case "Antimatter":
		*d = DCAntimatter
	default:
		return fmt.Errorf("trying to unmarshall unexpected DriveClassification value: %s: %w ", data, ErrBadValue)
	}

	return nil
}

//go:generate stringer -type=CoolingCycle
type CoolingCycle int

const (
	CCCalc CoolingCycle = iota
	CCOpen
	CCClosed
)

func (c *CoolingCycle) UnmarshalJSON(data []byte) error {
	s := string(data)
	if !(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		return fmt.Errorf("trying to unmarshall unexpected Propellant value: %s: %w", data, ErrBadValue)
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
	switch s {
	case "Calc":
		*c = CCCalc
	case "Open":
		*c = CCOpen
	case "Closed":
		*c = CCClosed
	default:
		return fmt.Errorf("trying to unmarshall unexpected Propellant value: %s: %w ", data, ErrBadValue)
	}

	return nil
}

//go:generate stringer -type=Propellant
type Propellant int

const (
	PropReactionProducts Propellant = iota
	PropAnything
	PropHydrogen
	PropWater
	PropNobleGases
	PropVolatiles
	PropMetals
)

func (p *Propellant) UnmarshalJSON(data []byte) error {
	s := string(data)
	if !(strings.HasPrefix(s, `"`) && strings.HasSuffix(s, `"`)) {
		return fmt.Errorf("trying to unmarshall unexpected Propellant value: %s: %w", data, ErrBadValue)
	}

	s = strings.TrimSuffix(strings.TrimPrefix(s, `"`), `"`)
	switch s {
	case "ReactionProducts":
		*p = PropReactionProducts
	case "Anything":
		*p = PropAnything
	case "Hydrogen":
		*p = PropHydrogen
	case "Water":
		*p = PropWater
	case "NobleGases":
		*p = PropNobleGases
	case "Volatiles":
		*p = PropVolatiles
	case "Metals":
		*p = PropMetals
	default:
		return fmt.Errorf("trying to unmarshall unexpected Propellant value: %s: %w ", data, ErrBadValue)
	}

	return nil
}
