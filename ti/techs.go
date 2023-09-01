package ti

import (
	"fmt"
	"maps"
)

type TechTemplate struct {
	GenericTemplate
	ResearchCost     int      `json:"researchCost,omitempty"`
	Prereqs          []string `json:"prereqs,omitempty"`
	immediatePrereqs map[*TechTemplate]struct{}
	allPrereqs       map[*TechTemplate]struct{}
	effectiveCost    int
}

func (t *TechTemplate) buildImmediatePrereqs(tMap map[string]*TechTemplate, pMap map[string]*TechTemplate) error {
	t.immediatePrereqs = make(map[*TechTemplate]struct{}, len(t.Prereqs))
	for _, prereq := range t.Prereqs {
		prereqPtr := tMap[prereq]
		if prereqPtr != nil {
			t.immediatePrereqs[prereqPtr] = struct{}{}

			continue
		}

		prereqPtr = pMap[prereq]
		if prereqPtr != nil {
			t.immediatePrereqs[prereqPtr] = struct{}{}

			continue
		}

		return fmt.Errorf("unknown Prereqs %s in %s: %w", prereq, t.DataName, ErrBadValue)
	}

	return nil
}

func (t *TechTemplate) buildAllPrereqs() map[*TechTemplate]struct{} {
	if t.allPrereqs != nil {
		return t.allPrereqs
	}

	t.allPrereqs = make(map[*TechTemplate]struct{})
	for immediatePrereq := range t.immediatePrereqs {
		t.allPrereqs[immediatePrereq] = struct{}{}
		maps.Copy(t.allPrereqs, immediatePrereq.buildAllPrereqs())
	}

	return t.allPrereqs
}

func SumTechEffectiveRPCost(techs ...*TechTemplate) int {
	allTechs := make(map[*TechTemplate]struct{})

	for _, tech := range techs {
		if tech != nil {
			allTechs[tech] = struct{}{}
			maps.Copy(allTechs, tech.allPrereqs)
		}
	}

	var cost int
	for tech := range allTechs {
		cost += tech.effectiveCost
	}

	return cost
}
