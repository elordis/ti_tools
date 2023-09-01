package ti

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type GenericTemplate struct {
	DataName     string `json:"dataName,omitempty"`
	FriendlyName string `json:"friendlyName,omitempty"`
}

func (t *GenericTemplate) GetDataName() string {
	return t.DataName
}

type BuildMaterialsWeights struct {
	Water       UnmarshallableFloat64 `json:"water,omitempty"`
	Volatiles   UnmarshallableFloat64 `json:"volatiles,omitempty"`
	Metals      UnmarshallableFloat64 `json:"metals,omitempty"`
	NobleMetals UnmarshallableFloat64 `json:"nobleMetals,omitempty"`
	Fissiles    UnmarshallableFloat64 `json:"fissiles,omitempty"`
	Antimatter  UnmarshallableFloat64 `json:"antimatter,omitempty"`
	Exotics     UnmarshallableFloat64 `json:"exotics,omitempty"`
}

func (b *BuildMaterialsWeights) NewBuildMaterials(m float64) *BuildMaterials {
	return &BuildMaterials{
		Water:       MaterialsToTons * m * float64(b.Water),
		Volatiles:   MaterialsToTons * m * float64(b.Volatiles),
		Metals:      MaterialsToTons * m * float64(b.Metals),
		NobleMetals: MaterialsToTons * m * float64(b.NobleMetals),
		Fissiles:    MaterialsToTons * m * float64(b.Fissiles),
		Antimatter:  MaterialsToTons * m * float64(b.Antimatter),
		Exotics:     MaterialsToTons * m * float64(b.Exotics),
	}
}

type BuildMaterials struct {
	Water       float64 `json:"water,omitempty"`
	Volatiles   float64 `json:"volatiles,omitempty"`
	Metals      float64 `json:"metals,omitempty"`
	NobleMetals float64 `json:"nobleMetals,omitempty"`
	Fissiles    float64 `json:"fissiles,omitempty"`
	Antimatter  float64 `json:"antimatter,omitempty"`
	Exotics     float64 `json:"exotics,omitempty"`
}

// Add sums build materials amounts. Adding to self id undefined behavior.
func (b *BuildMaterials) Add(a ...*BuildMaterials) *BuildMaterials {
	for _, bm := range a {
		if bm == nil {
			continue
		}

		b.Water += bm.Water
		b.Volatiles += bm.Volatiles
		b.Metals += bm.Metals
		b.NobleMetals += bm.NobleMetals
		b.Fissiles += bm.Fissiles
		b.Antimatter += bm.Antimatter
		b.Exotics += bm.Exotics
	}

	return b
}

func (b *BuildMaterials) Mul(m float64) *BuildMaterials {
	b.Water *= m
	b.Volatiles *= m
	b.Metals *= m
	b.NobleMetals *= m
	b.Fissiles *= m
	b.Antimatter *= m
	b.Exotics *= m

	return b
}
func (b *BuildMaterials) Div(m float64) *BuildMaterials {
	b.Water /= m
	b.Volatiles /= m
	b.Metals /= m
	b.NobleMetals /= m
	b.Fissiles /= m
	b.Antimatter /= m
	b.Exotics /= m

	return b
}

func (b *BuildMaterials) Ceil() *BuildMaterials {
	b.Water = math.Ceil(b.Water)
	b.Volatiles = math.Ceil(b.Volatiles)
	b.Metals = math.Ceil(b.Metals)
	b.NobleMetals = math.Ceil(b.NobleMetals)
	b.Fissiles = math.Ceil(b.Fissiles)
	b.Antimatter = math.Ceil(b.Antimatter)
	b.Exotics = math.Ceil(b.Exotics)

	return b
}

func (b *BuildMaterials) String() string {
	return fmt.Sprintf(
		"Wa: %.0f Vo: %.0f Me: %.0f No: %.0f Fi: %.0f AM: %.0f Ex: %.0f",
		b.Water,
		b.Volatiles,
		b.Metals,
		b.NobleMetals,
		b.Fissiles,
		b.Antimatter,
		b.Exotics,
	)
}

func NewCrewMaterials(crew int) *BuildMaterials {
	return &BuildMaterials{ //nolint:exhaustruct
		Water:     MaterialsToTons * float64(crew) * CrewWaterCost,
		Volatiles: MaterialsToTons * float64(crew) * CrewVolatilesCost,
	}
}

type DataNamer interface {
	GetDataName() string
}

func LoadTemplate[T DataNamer](fileName string) ([]T, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("opening %s: %w", fileName, err)
	}

	var slice []T

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&slice)
	if err != nil {
		return nil, fmt.Errorf("unmarshalling %s: %w", fileName, err)
	}

	return slice, nil
}

func MakeMap[T DataNamer](slice []T) map[string]T {
	m := make(map[string]T, len(slice))
	for _, template := range slice {
		m[template.GetDataName()] = template
	}

	return m
}
