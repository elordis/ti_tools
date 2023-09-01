package ti

import (
	"fmt"
	"math"
)

type HabSiteTemplate struct {
	GenericTemplate
	ParentBodyName    string `json:"parentBodyName,omitempty"`
	MiningProfileName string `json:"miningProfileName,omitempty"`
	miningProfile     *MiningProfileTemplate
}

func (h *HabSiteTemplate) makeMiningProfile(m map[string]*MiningProfileTemplate) error {
	h.miningProfile = m[h.MiningProfileName]
	if h.miningProfile == nil {
		return fmt.Errorf("unknown MiningProfileName %s in %s: %w", h.MiningProfileName, h.DataName, ErrBadValue)
	}

	return nil
}

type MiningProfileTemplate struct {
	GenericTemplate
	WaterMean        float64 `json:"water_mean,omitempty"`
	WaterWidth       float64 `json:"water_width,omitempty"`
	WaterMin         float64 `json:"water_min,omitempty"`
	WaterJump        float64 `json:"water_jump,omitempty"`
	VolatilesMean    float64 `json:"volatiles_mean,omitempty"`
	VolatilesWidth   float64 `json:"volatiles_width,omitempty"`
	VolatilesMin     float64 `json:"volatiles_min,omitempty"`
	VolatilesJump    float64 `json:"volatiles_jump,omitempty"`
	MetalsMean       float64 `json:"metals_mean,omitempty"`
	MetalsWidth      float64 `json:"metals_width,omitempty"`
	MetalsMin        float64 `json:"metals_min,omitempty"`
	MetalsJump       float64 `json:"metals_jump,omitempty"`
	NobleMetalsMean  float64 `json:"nobles_mean,omitempty"`
	NobleMetalsWidth float64 `json:"nobles_width,omitempty"`
	NobleMetalsMin   float64 `json:"nobles_min,omitempty"`
	NobleMetalsJump  float64 `json:"nobles_jump,omitempty"`
	FissilesMean     float64 `json:"fissiles_mean,omitempty"`
	FissilesWidth    float64 `json:"fissiles_width,omitempty"`
	FissilesMin      float64 `json:"fissiles_min,omitempty"`
	FissilesJump     float64 `json:"fissiles_jump,omitempty"`
}

func makeYield(mean, width, _min, jump float64) float64 {
	if mean <= 0 && _min <= 0 && width <= 0 {
		return 0
	}

	posAdjust := width * jump / (1 - jump)
	negAadjust := 0.0

	t := math.Ceil((mean - _min) / width)
	if t > 0.0 {
		for i := 0.0; i < t; i++ {
			negAadjust += width * i * math.Pow(jump, i)
		}

		negAadjust /= t
	}

	p := math.Pow(jump, t)

	average := 0.5*(mean+posAdjust) + 0.5*(p*_min+(1-p)*(mean-negAadjust)) //nolint:gomnd
	if average < 0.01 {                                                    //nolint:gomnd
		return 0
	}

	return 0.5*(mean+posAdjust) + 0.5*(p*_min+(1-p)*(mean-negAadjust))
}

// MakeAverageYields calculates average yield of a mining site. Note that it gives correct results only for mining
// profiles with "modifyBySize" set to false.
func (m *MiningProfileTemplate) MakeAverageYields() *BuildMaterials {
	materials := BuildMaterials{ //nolint:exhaustruct
		Water:       makeYield(m.WaterMean, m.WaterWidth, m.WaterMin, m.WaterJump),
		Volatiles:   makeYield(m.VolatilesMean, m.VolatilesWidth, m.VolatilesMin, m.VolatilesJump),
		Metals:      makeYield(m.MetalsMean, m.MetalsWidth, m.MetalsMin, m.MetalsJump),
		NobleMetals: makeYield(m.NobleMetalsMean, m.NobleMetalsWidth, m.NobleMetalsMin, m.NobleMetalsJump),
		Fissiles:    makeYield(m.FissilesMean, m.FissilesWidth, m.FissilesMin, m.FissilesJump),
	}
	materials.NobleMetals = min(materials.NobleMetals, materials.Metals/3) //nolint:gomnd

	return &materials
}
