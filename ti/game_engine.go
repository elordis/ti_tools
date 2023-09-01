package ti

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path/filepath"
)

// Constants used in Terra Invicta source code.
const (
	G                 = 9.80665
	MaterialsToTons   = 0.1
	CrewWaterCost     = 2
	CrewVolatilesCost = 2
	DaysPerMonth      = 30.436874
	FuelTPerTank      = 100
)

type GameEngine struct {
	DriveTemplates         map[string]*DriveTemplate
	PowerPlantTemplates    map[string]*PowerPlantTemplate
	RadiatorTemplates      map[string]*RadiatorTemplate
	UtilityModuleTemplates map[string]*UtilityModuleTemplate
	TechTemplates          map[string]*TechTemplate
	ProjectTemplates       map[string]*TechTemplate
	HabSiteTemplates       map[string]*HabSiteTemplate
	MiningProfileTemplates map[string]*MiningProfileTemplate
	SimulationConstraints  SimulationConstraints
	driveTemplates         []*DriveTemplate
	powerPlantTemplates    []*PowerPlantTemplate
	radiatorTemplates      []*RadiatorTemplate
	utilityModuleTemplates []*UtilityModuleTemplate
	techTemplates          []*TechTemplate
	projectTemplates       []*TechTemplate
	habSiteTemplates       []*HabSiteTemplate
	miningProfileTemplates []*MiningProfileTemplate
	evUtilities            []*UtilityModuleTemplate
	thrustUtilities        []*UtilityModuleTemplate
	miningAmount           *BuildMaterials
}

func (e *GameEngine) initShipParts() error {
	for _, t := range e.driveTemplates {
		err := t.makeRequiredProject(e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	for _, t := range e.powerPlantTemplates {
		err := t.makeRequiredProject(e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	for _, t := range e.utilityModuleTemplates {
		err := t.makeRequiredProject(e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	e.evUtilities = make([]*UtilityModuleTemplate, 0, len(e.UtilityModuleTemplates))
	e.evUtilities = append(e.evUtilities, nil)
	e.thrustUtilities = make([]*UtilityModuleTemplate, 0, len(e.UtilityModuleTemplates))
	e.thrustUtilities = append(e.thrustUtilities, nil)

	for _, t := range e.utilityModuleTemplates {
		if t.EVMult > 1 {
			e.evUtilities = append(e.evUtilities, t)
		}

		if t.ThrsutMult > 1 {
			e.thrustUtilities = append(e.thrustUtilities, t)
		}
	}

	for _, t := range e.radiatorTemplates {
		err := t.makeRequiredProject(e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	return nil
}

func (e *GameEngine) initTechs() error {
	for _, t := range e.techTemplates {
		err := t.buildImmediatePrereqs(e.TechTemplates, e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	for _, t := range e.techTemplates {
		t.buildAllPrereqs()
	}

	for _, t := range e.projectTemplates {
		err := t.buildImmediatePrereqs(e.TechTemplates, e.ProjectTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	for _, t := range e.projectTemplates {
		t.buildAllPrereqs()
	}

	return nil
}

func (e *GameEngine) initMining() error {
	for _, t := range e.habSiteTemplates {
		err := t.makeMiningProfile(e.MiningProfileTemplates)
		if err != nil {
			return fmt.Errorf("building engine cache: %w", err)
		}
	}

	return nil
}

func NewGameEngine(templatePath string) (*GameEngine, error) {
	var e GameEngine

	var err error

	e.driveTemplates, err = LoadTemplate[*DriveTemplate](filepath.Join(templatePath, "TIDriveTemplate.json"))
	if err != nil {
		return nil, fmt.Errorf("loading drives: %w", err)
	}

	e.powerPlantTemplates, err = LoadTemplate[*PowerPlantTemplate](filepath.Join(templatePath, "TIPowerPlantTemplate.json")) //nolint:lll
	if err != nil {
		return nil, fmt.Errorf("loading power plants: %w", err)
	}

	e.radiatorTemplates, err = LoadTemplate[*RadiatorTemplate](filepath.Join(templatePath, "TIRadiatorTemplate.json"))
	if err != nil {
		return nil, fmt.Errorf("loading radiators: %w", err)
	}

	e.utilityModuleTemplates, err = LoadTemplate[*UtilityModuleTemplate](filepath.Join(templatePath, "TIUtilityModuleTemplate.json")) //nolint:lll
	if err != nil {
		return nil, fmt.Errorf("loading utility modules: %w", err)
	}

	e.techTemplates, err = LoadTemplate[*TechTemplate](filepath.Join(templatePath, "TITechTemplate.json"))
	if err != nil {
		return nil, fmt.Errorf("loading techs: %w", err)
	}

	e.projectTemplates, err = LoadTemplate[*TechTemplate](filepath.Join(templatePath, "TIProjectTemplate.json"))
	if err != nil {
		return nil, fmt.Errorf("loading projects: %w", err)
	}

	e.habSiteTemplates, err = LoadTemplate[*HabSiteTemplate](filepath.Join(templatePath, "TIHabSiteTemplate.json"))
	if err != nil {
		return nil, fmt.Errorf("loading hab sites: %w", err)
	}

	e.miningProfileTemplates, err = LoadTemplate[*MiningProfileTemplate](filepath.Join(templatePath, "TIMiningProfileTemplate.json")) //nolint:lll
	if err != nil {
		return nil, fmt.Errorf("loading mining profiles: %w", err)
	}

	e.DriveTemplates = MakeMap(e.driveTemplates)
	e.PowerPlantTemplates = MakeMap(e.powerPlantTemplates)
	e.RadiatorTemplates = MakeMap(e.radiatorTemplates)
	e.UtilityModuleTemplates = MakeMap(e.utilityModuleTemplates)
	e.TechTemplates = MakeMap(e.techTemplates)
	e.ProjectTemplates = MakeMap(e.projectTemplates)
	e.HabSiteTemplates = MakeMap(e.habSiteTemplates)
	e.MiningProfileTemplates = MakeMap(e.miningProfileTemplates)

	err = e.initShipParts()
	if err != nil {
		return nil, err
	}

	err = e.initTechs()
	if err != nil {
		return nil, err
	}

	err = e.initMining()
	if err != nil {
		return nil, err
	}

	return &e, nil
}

func (e *GameEngine) ApplyConstraints(c SimulationConstraints) {
	e.SimulationConstraints = c

	for _, t := range e.techTemplates {
		t.effectiveCost = int(float64(t.ResearchCost) / (1 + e.SimulationConstraints.TechBonus))
	}

	for _, t := range e.projectTemplates {
		t.effectiveCost = int(float64(t.ResearchCost) / (1 + e.SimulationConstraints.ProjectBonus))
	}

	bodiesMap := make(map[string]struct{}, len(e.SimulationConstraints.MiningBodies))
	for _, body := range e.SimulationConstraints.MiningBodies {
		bodiesMap[body] = struct{}{}
	}

	e.miningAmount = &BuildMaterials{} //nolint:exhaustruct
	for _, habSite := range e.habSiteTemplates {
		if _, ok := bodiesMap[habSite.ParentBodyName]; ok {
			e.miningAmount.Add(habSite.miningProfile.MakeAverageYields())
		}
	}

	e.miningAmount.Mul(1 + e.SimulationConstraints.MiningBonus)
	e.miningAmount.Add(e.SimulationConstraints.MiningFlat)
}

func (e *GameEngine) BuildMaterialsMiningDays(bm *BuildMaterials) float64 {
	var water, volatiles, metals, nobleMetals, fissiles, antimatter, exotics float64
	if e.miningAmount.Water != 0.0 {
		water = bm.Water / e.miningAmount.Water
	} else if bm.Water != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.Volatiles != 0.0 {
		volatiles = bm.Volatiles / e.miningAmount.Volatiles
	} else if bm.Volatiles != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.Metals != 0.0 {
		metals = bm.Metals / e.miningAmount.Metals
	} else if bm.Metals != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.NobleMetals != 0.0 {
		nobleMetals = bm.NobleMetals / e.miningAmount.NobleMetals
	} else if bm.NobleMetals != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.Fissiles != 0.0 {
		fissiles = bm.Fissiles / e.miningAmount.Fissiles
	} else if bm.Water != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.Antimatter != 0.0 {
		antimatter = bm.Antimatter / e.miningAmount.Antimatter
	} else if bm.Antimatter != 0.0 {
		return math.Inf(1)
	}

	if e.miningAmount.Exotics != 0.0 {
		exotics = bm.Exotics / e.miningAmount.Exotics
	} else if bm.Exotics != 0.0 {
		return math.Inf(1)
	}

	return DaysPerMonth * max(water, volatiles, metals, nobleMetals, fissiles, antimatter, exotics)
}

func (e *GameEngine) BuildMaterialsPerDay(bm *BuildMaterials) float64 {
	var water, volatiles, metals, nobleMetals, fissiles, antimatter, exotics float64
	if bm.Water == 0 {
		water = math.Inf(1)
	} else {
		water = e.miningAmount.Water / bm.Water
	}

	if bm.Volatiles == 0 {
		volatiles = math.Inf(1)
	} else {
		volatiles = e.miningAmount.Volatiles / bm.Volatiles
	}

	if bm.Metals == 0 {
		metals = math.Inf(1)
	} else {
		metals = e.miningAmount.Metals / bm.Metals
	}

	if bm.NobleMetals == 0 {
		nobleMetals = math.Inf(1)
	} else {
		nobleMetals = e.miningAmount.NobleMetals / bm.NobleMetals
	}

	if bm.Fissiles == 0 {
		fissiles = math.Inf(1)
	} else {
		fissiles = e.miningAmount.Fissiles / bm.Fissiles
	}

	if bm.Antimatter == 0 {
		antimatter = math.Inf(1)
	} else {
		antimatter = e.miningAmount.Antimatter / bm.Antimatter
	}

	if bm.Exotics == 0 {
		exotics = math.Inf(1)
	} else {
		exotics = e.miningAmount.Exotics / bm.Exotics
	}

	return min(water, volatiles, metals, nobleMetals, fissiles, antimatter, exotics) / DaysPerMonth
}

func (e *GameEngine) ForAllDriveAssemblies(includeAlien bool, f func(*DriveAssembly)) { //nolint:gocognit
	for _, drive := range e.driveTemplates {
		if !includeAlien && drive.IsAlien() {
			continue
		}

		for _, powerPlant := range e.powerPlantTemplates {
			if !includeAlien && powerPlant.IsAlien() || !powerPlant.CompatibleWith(drive) {
				continue
			}

			for _, radiator := range e.radiatorTemplates {
				if !includeAlien && radiator.IsAlien() {
					continue
				}

				for _, evUtility := range e.evUtilities {
					if evUtility != nil && (!includeAlien && evUtility.IsAlien() || !evUtility.CompatibleWith(drive)) {
						continue
					}

					for _, thrustUtility := range e.thrustUtilities {
						if thrustUtility != nil && (!includeAlien && thrustUtility.IsAlien() || !thrustUtility.CompatibleWith(drive)) {
							continue
						}

						d := DriveAssembly{
							Drive:         drive,
							PowerPlant:    powerPlant,
							Radiator:      radiator,
							EVUtility:     evUtility,
							ThrustUtility: thrustUtility,
							Engine:        e,
						}
						if d.IsValid() {
							f(&d)
						}
					}
				}
			}
		}
	}
}

func (e *GameEngine) MiningAmount() *BuildMaterials {
	b := *e.miningAmount

	return &b
}

// NewDriveAssembly returns DriveAssembly struct from provided parts dataNames. If requested assembly is invalid or
// parts are not found it returns nil.
func (e *GameEngine) NewDriveAssembly(drive, powerPlant, radiator, eVUtility, thrustUtility string) *DriveAssembly {
	d := DriveAssembly{
		Drive:         e.DriveTemplates[drive],
		PowerPlant:    e.PowerPlantTemplates[powerPlant],
		Radiator:      e.RadiatorTemplates[radiator],
		EVUtility:     e.UtilityModuleTemplates[eVUtility],
		ThrustUtility: e.UtilityModuleTemplates[thrustUtility],
		Engine:        e,
	}
	if !d.IsValid() {
		return nil
	}

	return &d
}

type SimulationConstraints struct {
	MiningBodies []string        `json:"mining_bodies,omitempty"`
	MiningBonus  float64         `json:"mining_bonus,omitempty"`
	MiningFlat   *BuildMaterials `json:"mining_flat,omitempty"`
	TechBonus    float64         `json:"tech_bonus,omitempty"`
	ProjectBonus float64         `json:"project_bonus,omitempty"`
}

func NewConstraints(fileName string) (SimulationConstraints, error) {
	var c SimulationConstraints

	file, err := os.Open(fileName)
	if err != nil {
		return c, fmt.Errorf("opening '%s': %w", fileName, err)
	}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&c)
	if err != nil {
		return c, fmt.Errorf("unmarshalling '%s': %w", fileName, err)
	}

	return c, nil
}
