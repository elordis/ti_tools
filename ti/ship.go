package ti

import (
	"fmt"
	"math"
)

type DriveAssembly struct {
	Drive         *DriveTemplate
	PowerPlant    *PowerPlantTemplate
	Radiator      *RadiatorTemplate
	EVUtility     *UtilityModuleTemplate
	ThrustUtility *UtilityModuleTemplate
	Engine        *GameEngine
}

func (d *DriveAssembly) String() string {
	if d == nil {
		return NameForMissing
	}

	driveName := NameForMissing
	if d.Drive != nil {
		driveName = d.Drive.FriendlyName
	}

	powerPlantName := NameForMissing
	if d.PowerPlant != nil {
		powerPlantName = d.PowerPlant.FriendlyName
	}

	radiatorName := NameForMissing
	if d.Radiator != nil {
		radiatorName = d.Radiator.FriendlyName
	}

	evName := NameForMissing
	if d.EVUtility != nil {
		evName = d.EVUtility.FriendlyName
	}

	thrustName := NameForMissing
	if d.ThrustUtility != nil {
		thrustName = d.ThrustUtility.FriendlyName
	}

	return fmt.Sprintf(
		"%s %s %s %s %s",
		abbreviate(driveName, AbbrLetters, true),
		abbreviate(powerPlantName, AbbrLetters, true),
		abbreviate(radiatorName, AbbrLetters, false),
		abbreviate(evName, AbbrLetters, false),
		abbreviate(thrustName, AbbrLetters, false),
	)
}

func (d *DriveAssembly) IsValid() bool {
	if !d.PowerPlant.CompatibleWith(d.Drive) {
		return false
	}

	if d.EVUtility != nil && !d.EVUtility.CompatibleWith(d.Drive) {
		return false
	}

	if d.ThrustUtility != nil && !d.ThrustUtility.CompatibleWith(d.Drive) {
		return false
	}

	return true
}

func (d *DriveAssembly) ModifiedEVKps() float64 {
	if d.EVUtility == nil {
		return d.Drive.EVKps
	}

	return d.Drive.EVKps * d.EVUtility.EVMult
}

func (d *DriveAssembly) ModifiedThrustN() float64 {
	if d.ThrustUtility == nil {
		return d.Drive.ThrustN
	}

	return d.Drive.ThrustN * d.ThrustUtility.ThrsutMult
}

func (d *DriveAssembly) TotalCrew() int {
	crew := d.PowerPlant.Crew + d.Radiator.Crew
	if d.EVUtility != nil {
		crew += d.EVUtility.Crew
	}

	if d.ThrustUtility != nil {
		crew += d.ThrustUtility.Crew
	}

	return crew
}

func (d *DriveAssembly) PowerPlantMassT() float64 {
	return d.PowerPlant.SpecificPowerTGW * d.Drive.PowerRequired()
}

func (d *DriveAssembly) PowerPlantBuildMaterials() *BuildMaterials {
	b := d.Radiator.WeightedBuildMaterials.NewBuildMaterials(d.PowerPlantMassT())
	b.Add(NewCrewMaterials(d.PowerPlant.Crew))

	return b
}

func (d *DriveAssembly) WasteHeatGW() float64 {
	crewHeat := 3.75 * float64(d.TotalCrew()) / 1000 //nolint:gomnd
	if d.Drive.IsOpenCycle() {
		return crewHeat
	}

	driveHeat := d.Drive.PowerRequired() * (1 - d.PowerPlant.Efficiency)

	return max(crewHeat, driveHeat)
}

func (d *DriveAssembly) RadiatorMassT() float64 {
	return d.WasteHeatGW() * 1000 / d.Radiator.SpecificPower
}

func (d *DriveAssembly) RadiatorBuildMaterials() *BuildMaterials {
	return d.Radiator.WeightedBuildMaterials.NewBuildMaterials(d.RadiatorMassT()).Add(NewCrewMaterials(d.Radiator.Crew))
}

func (d *DriveAssembly) MassT() float64 {
	mass := d.Drive.MassT() + d.PowerPlantMassT() + d.RadiatorMassT()
	if d.EVUtility != nil {
		mass += d.EVUtility.MassT
	}

	if d.ThrustUtility != nil {
		mass += d.ThrustUtility.MassT
	}

	return mass
}

func (d *DriveAssembly) BuildMaterials() *BuildMaterials {
	return (&BuildMaterials{}).Add( //nolint:exhaustruct
		d.Drive.BuildMaterials(),
		d.PowerPlantBuildMaterials(),
		d.RadiatorBuildMaterials(),
		d.EVUtility.BuildMaterials(),
		d.ThrustUtility.BuildMaterials(),
	)
}

func (d *DriveAssembly) FuelTForPayloadDV(payload float64, dv float64) float64 {
	return ((payload + d.MassT()) * (math.Pow(math.E, dv/d.ModifiedEVKps()) - 1))
}

func (d *DriveAssembly) TanksForPayloadDV(payload float64, dv float64) int {
	return int(math.Ceil(d.FuelTForPayloadDV(payload, dv) / FuelTPerTank))
}

func (d *DriveAssembly) EffectiveRPCost() int {
	var evUtil, thrustUtil *PartTemplate
	if d.EVUtility != nil {
		evUtil = &d.EVUtility.PartTemplate
	}

	if d.ThrustUtility != nil {
		thrustUtil = &d.ThrustUtility.PartTemplate
	}

	return SumPartsEffectiveRPCost(
		&d.Drive.PartTemplate,
		&d.PowerPlant.PartTemplate,
		&d.Radiator.PartTemplate,
		evUtil,
		thrustUtil,
	)
}

func (d *DriveAssembly) MaxPayloadConstrained(dv, cruiseAccel, combatAccel float64) float64 {
	cruiseAccelMps2 := cruiseAccel * G
	cruiseTotalMassT := d.ModifiedThrustN() / cruiseAccelMps2 / 1000 //nolint:gomnd
	cruiseFuelT := cruiseTotalMassT * (1 - math.Pow(math.E, -dv/d.ModifiedEVKps()))
	cruisePayload := cruiseTotalMassT - d.MassT() - FuelTPerTank*max(1, math.Ceil(cruiseFuelT/FuelTPerTank))

	combatAccelMps2 := combatAccel * G
	combatTotalMassT := d.Drive.ThrustCap * d.ModifiedThrustN() / combatAccelMps2 / 1000 //nolint:gomnd
	combatFuelT := combatTotalMassT * (1 - math.Pow(math.E, -dv/d.ModifiedEVKps()))
	combatPayload := combatTotalMassT - d.MassT() - FuelTPerTank*max(1, math.Ceil(combatFuelT/FuelTPerTank))

	return max(0, min(cruisePayload, combatPayload))
}

type Ship struct {
	DriveAssembly *DriveAssembly
	Tanks         int
	PayloadMassT  float64
}

func (s *Ship) DryMassT() float64 {
	return s.PayloadMassT + s.DriveAssembly.MassT()
}

func (s *Ship) FuelMassT() float64 {
	return FuelTPerTank * float64(s.Tanks)
}

func (s *Ship) WetMassT() float64 {
	return s.DryMassT() + s.FuelMassT()
}

func (s *Ship) DVKps() float64 {
	return s.DriveAssembly.ModifiedEVKps() * math.Log(s.WetMassT()/s.DryMassT())
}

func (s *Ship) CruiseAccelerationG() float64 {
	a := s.DriveAssembly.ModifiedThrustN() / (1000 * s.WetMassT()) / G //nolint:gomnd

	return min(a, 2) //nolint:gomnd
}

func (s *Ship) CombatAccelerationG() float64 {
	a := s.DriveAssembly.Drive.ThrustCap * s.DriveAssembly.ModifiedThrustN() / (1000 * s.WetMassT()) / G //nolint:gomnd

	return min(a, 4) //nolint:gomnd
}

func (s *Ship) FuelBuildMaterials() *BuildMaterials {
	return s.DriveAssembly.Drive.PerTankPropellantMaterials.NewBuildMaterials(s.FuelMassT())
}

func (s *Ship) BuildMaterials() *BuildMaterials {
	return (&BuildMaterials{}).Add( //nolint:exhaustruct
		s.DriveAssembly.BuildMaterials(),
		s.FuelBuildMaterials(),
	)
}

func (s *Ship) PayloadTPerMiningDay() float64 {
	return s.PayloadMassT * s.DriveAssembly.Engine.BuildMaterialsPerDay(s.BuildMaterials())
}
