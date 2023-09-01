package ti

import "fmt"

type PartTemplate struct {
	GenericTemplate
	RequiredProjectName    string                `json:"requiredProjectName,omitempty"`
	WeightedBuildMaterials BuildMaterialsWeights `json:"weightedBuildMaterials,omitempty"`
	requiredProject        *TechTemplate
}

func (p *PartTemplate) makeRequiredProject(m map[string]*TechTemplate) error {
	if p.RequiredProjectName == "" {
		return nil
	}

	p.requiredProject = m[p.RequiredProjectName]
	if p.requiredProject == nil {
		return fmt.Errorf("unknown RequiredProjectName %s in %s: %w", p.RequiredProjectName, p.DataName, ErrBadValue)
	}

	return nil
}

func (p *PartTemplate) IsAlien() bool {
	return p.RequiredProjectName == "Project_AlienMasterProject"
}

func SumPartsEffectiveRPCost(parts ...*PartTemplate) int {
	techs := make([]*TechTemplate, 0, len(parts))

	for _, part := range parts {
		if part != nil && part.requiredProject != nil {
			techs = append(techs, part.requiredProject)
		}
	}

	return SumTechEffectiveRPCost(techs...)
}

type DriveTemplate struct {
	PartTemplate
	Thrusters                  int                   `json:"thrusters,omitempty"`
	DriveClassification        DriveClassification   `json:"driveClassification,omitempty"`
	ThrustN                    float64               `json:"thrust_N,omitempty"`
	EVKps                      float64               `json:"EV_kps,omitempty"`
	SpecificPowerKGMW          float64               `json:"specificPower_kgMW,omitempty"`
	Efficiency                 float64               `json:"efficiency,omitempty"`
	FlatMassT                  float64               `json:"flatMass_tons,omitempty"`
	RequiredPowerPlant         PowerPlantClass       `json:"requiredPowerPlant,omitempty"`
	ThrustCap                  float64               `json:"thrustCap,omitempty"`
	Cooling                    CoolingCycle          `json:"cooling,omitempty"`
	Propellant                 Propellant            `json:"propellant,omitempty"`
	PerTankPropellantMaterials BuildMaterialsWeights `json:"perTankPropellantMaterials,omitempty"`
}

func (d *DriveTemplate) IsSelfPowered() bool {
	switch d.DriveClassification { //nolint:exhaustive
	case DCChemical, DCFissionPulse, DCFusionPulse:
		return true
	default:
		return false
	}
}

func (d *DriveTemplate) IsPulsed() bool {
	switch d.DriveClassification { //nolint:exhaustive
	case DCFissionPulse, DCFusionPulse:
		return true
	default:
		return false
	}
}

func (d *DriveTemplate) IsNuclear() bool {
	diff := d.RequiredPowerPlant - PPCSolidCoreFission

	return 0 <= diff && diff <= 14
}

func (d *DriveTemplate) IsFusion() bool {
	diff := d.RequiredPowerPlant - PPCZPinchFusion

	return 0 <= diff && diff <= 6
}

func (d *DriveTemplate) ThrustPowerGW() float64 {
	return d.ThrustN * d.EVKps / 2 / 1000000 //nolint:gomnd
}

func (d *DriveTemplate) MassFlowKg() float64 {
	return d.ThrustN / (float64(d.Thrusters) * d.EVKps * 1000) //nolint:gomnd
}

func (d *DriveTemplate) IsOpenCycle() bool {
	switch d.Cooling { //nolint:exhaustive
	case CCOpen:
		return true
	case CCCalc:
		return d.IsPulsed() || (d.MassFlowKg() >= 3) //nolint:gomnd
	default:
		return false
	}
}

func (d *DriveTemplate) PowerRequired() float64 {
	if d.IsSelfPowered() {
		return 0
	}

	return d.ThrustPowerGW() / d.Efficiency
}

func (d *DriveTemplate) MassT() float64 {
	return d.FlatMassT + d.ThrustPowerGW()*d.SpecificPowerKGMW
}

func (d *DriveTemplate) BuildMaterials() *BuildMaterials {
	return d.WeightedBuildMaterials.NewBuildMaterials(d.MassT())
}

type PowerPlantTemplate struct {
	PartTemplate
	MaxOutputGW      float64         `json:"maxOutput_GW,omitempty"`
	SpecificPowerTGW float64         `json:"specificPower_tGW,omitempty"`
	PowerPlantClass  PowerPlantClass `json:"powerPlantClass,omitempty"`
	Efficiency       float64         `json:"efficiency,omitempty"`
	Crew             int             `json:"crew,omitempty"`
}

func (p *PowerPlantTemplate) CompatibleWith(d *DriveTemplate) bool {
	var classCompatible bool

	switch d.RequiredPowerPlant { //nolint:exhaustive
	case PPCAnyGeneral:
		classCompatible = true
	case PPCAnyMagneticConfinementFusion:
		diff := p.PowerPlantClass - PPCHybridConfinementFusion
		classCompatible = (0 <= diff && diff <= 2)
	default:
		classCompatible = (p.PowerPlantClass == d.RequiredPowerPlant)
	}

	powerCompatible := (p.MaxOutputGW >= d.PowerRequired())

	return classCompatible && powerCompatible
}

type UtilityModuleTemplate struct {
	PartTemplate
	MassT                      float64 `json:"mass_tons,omitempty"`
	ThrsutMult                 float64 `json:"thrustMultiplier,omitempty"`
	EVMult                     float64 `json:"EVMultiplier,omitempty"`
	RequiresHydrogenPropellant bool    `json:"requiresHydrogenPropellant,omitempty"`
	RequiresNuclearDrive       bool    `json:"requiresNuclearDrive,omitempty"`
	RequiresFusionDrive        bool    `json:"requiresFusionDrive,omitempty"`
	Crew                       int     `json:"crew,omitempty"`
}

func (u *UtilityModuleTemplate) CompatibleWith(d *DriveTemplate) bool {
	if u == nil {
		return true
	}

	fusionCompatible := d.IsFusion() || !u.RequiresFusionDrive
	nuclearCompatible := d.IsNuclear() || !u.RequiresNuclearDrive
	hydrogenCompatible := (d.Propellant == PropHydrogen) || !u.RequiresHydrogenPropellant

	return fusionCompatible && nuclearCompatible && hydrogenCompatible
}

func (u *UtilityModuleTemplate) BuildMaterials() *BuildMaterials {
	if u == nil {
		return &BuildMaterials{} //nolint:exhaustruct
	}

	return u.WeightedBuildMaterials.NewBuildMaterials(u.MassT).Add(NewCrewMaterials(u.Crew))
}

type RadiatorTemplate struct {
	PartTemplate
	SpecificPower float64 `json:"specificPower_2s_KWkg,omitempty"`
	Crew          int     `json:"crew,omitempty"`
}
