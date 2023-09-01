package main

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"slices"

	"github.com/elordis/ti_drive_plot/ti"
	"github.com/elordis/ti_drive_plot/tiplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/font"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/text"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func ShipIsWorseThan(s *ti.Ship) func(*ti.Ship) bool {
	return func(other *ti.Ship) bool {
		return (s.DriveAssembly.Drive == other.DriveAssembly.Drive &&
			s.PayloadMassT > other.PayloadMassT &&
			s.PayloadTPerMiningDay() > other.PayloadTPerMiningDay())
	}
}

func ShipIsBetterThan(s *ti.Ship) func(*ti.Ship) bool {
	return func(other *ti.Ship) bool {
		return (s.DriveAssembly.Drive == other.DriveAssembly.Drive &&
			other.PayloadMassT > s.PayloadMassT &&
			other.PayloadTPerMiningDay() > s.PayloadTPerMiningDay())
	}
}

func CmpMaxPayload(a, b *ti.Ship) int {
	return cmp.Compare(a.PayloadMassT, b.PayloadMassT)
}

func CmpMaxPtmpd(a, b *ti.Ship) int {
	return cmp.Compare(a.PayloadTPerMiningDay(), b.PayloadTPerMiningDay())
}

func ResultPlotters(s []*ti.Ship) []plot.Plotter {
	xys := make(plotter.XYs, len(s))
	names := make([]string, len(s))
	styles := make([]text.Style, len(s))

	for i := range s {
		xys[i].X, xys[i].Y = s[i].PayloadMassT, s[i].PayloadTPerMiningDay()
		names[i] = s[i].DriveAssembly.Drive.String()
		styles[i] = text.Style{ //nolint:exhaustruct
			Font:     font.From(plotter.DefaultFont, plotter.DefaultFontSize),
			Rotation: tiplot.FourtyFiveDegree,
			Handler:  plot.DefaultTextHandler,
		}
	}

	return []plot.Plotter{
		&plotter.Scatter{ //nolint:exhaustruct
			XYs: xys,
			GlyphStyleFunc: func(i int) draw.GlyphStyle {
				return tiplot.DriveGlyph(s[i].DriveAssembly.Drive)
			},
		},
		&plotter.Labels{
			XYs:       xys,
			Labels:    names,
			TextStyle: styles,
			Offset:    vg.Point{X: 4, Y: 4}, //nolint:gomnd
		},
	}
}

func main() {
	templatePath := flag.String("t", "", "folder containing template JSON files")
	constraintsFile := flag.String("c", "", "JSON-encoded constraints file")
	minDV := flag.Float64("dv", 10, "minimum delta-V to target")                           //nolint:gomnd
	minCruiseAccel := flag.Float64("cra", 0.000001, "minimum allowed cruise acceleration") //nolint:gomnd
	minCombatAccel := flag.Float64("coa", 0.000001, "minimum allowed combat acceleration") //nolint:gomnd
	maxRPCost := flag.Int("rp", 10000000, "maximum allowed RP cost for drive assembly")    //nolint:gomnd
	minPayload := flag.Float64("pl", 200, "maximum allowed RP cost for drive assembly")    //nolint:gomnd
	minPtpmd := flag.Float64("ptpmd", 200, "maximum allowed RP cost for drive assembly")   //nolint:gomnd
	outputFile := flag.String("o", "output.png", "file to write output to")
	outputSize := flag.Int("size", 1000, "size of output image in pixels") //nolint:gomnd
	logScale := flag.Bool("log", false, "use logarighmic scale")
	cpuprofile := flag.String("cpuprofile", "", "(DEBUG) write cpu profile to file")

	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatalf("creating cpuprofile: %s", err)
		}

		err = pprof.StartCPUProfile(f)
		if err != nil {
			log.Fatalf("starting cpuprofile: %s", err)
		}
	}

	log.Println("init started")

	engine, err := ti.NewGameEngine(*templatePath)
	if err != nil {
		log.Fatalf("intializing engine: %s", err)
	}

	constraints, err := ti.NewConstraints(*constraintsFile)
	if err != nil {
		log.Fatalf("intializing constraints: %s", err)
	}

	engine.ApplyConstraints(constraints)
	log.Printf("mining amount: %s\n", engine.MiningAmount())

	log.Println("simulation started")

	results := make([]*ti.Ship, 0, len(engine.DriveTemplates))
	engine.ForAllDriveAssemblies(false, func(d *ti.DriveAssembly) {
		rpCost := d.EffectiveRPCost()
		if rpCost > *maxRPCost {
			return
		}
		maxPayload := d.MaxPayloadConstrained(*minDV, *minCruiseAccel, *minCombatAccel)
		if maxPayload < *minPayload {
			return
		}
		currentShip := &ti.Ship{
			DriveAssembly: d,
			Tanks:         d.TanksForPayloadDV(maxPayload, *minDV),
			PayloadMassT:  maxPayload,
		}
		if currentShip.PayloadTPerMiningDay() < *minPtpmd {
			return
		}

		if slices.ContainsFunc(results, ShipIsBetterThan(currentShip)) {
			return
		}
		results = slices.DeleteFunc(results, ShipIsWorseThan(currentShip))
		results = append(results, currentShip)
	})

	for _, r := range results {
		log.Printf("%s: %.0f T, %.0f T/day\n", r.DriveAssembly, r.PayloadMassT, r.PayloadTPerMiningDay())
	}

	p := plot.New()
	p.Title.Text = "Payload per Mining Day / Payload Diagram\n"
	p.Title.Text += fmt.Sprintf(
		"Mining: %.2f*%v+(%s)\n",
		1+engine.SimulationConstraints.MiningBonus,
		engine.SimulationConstraints.MiningBodies,
		engine.SimulationConstraints.MiningFlat,
	)
	p.Title.Text += fmt.Sprintf(
		"Research: %.2f Tech Bonus, %.2f Project Bonus, \n",
		engine.SimulationConstraints.TechBonus,
		engine.SimulationConstraints.ProjectBonus,
	)
	p.Title.Text += fmt.Sprintf(
		"DV: %.0f Kps, Cruise: %.3f G, Combat: %.3f G, RP: %d, Payload: %.0f T, PTperMD: %.0f T/day",
		*minDV,
		*minCruiseAccel,
		*minCombatAccel,
		*maxRPCost,
		*minPayload,
		*minPtpmd,
	)
	p.X.Label.Text = "Payload (T)"
	p.Y.Label.Text = "Payload per Mining Day (T/day)"
	p.X.Min, p.X.Max = *minPayload, slices.MaxFunc(results, CmpMaxPayload).PayloadMassT
	p.Y.Min, p.Y.Max = *minPtpmd, slices.MaxFunc(results, CmpMaxPtmpd).PayloadTPerMiningDay()

	if *logScale {
		p.X.Scale, p.Y.Scale = plot.LogScale{}, plot.LogScale{}
		p.X.Tick.Marker, p.Y.Tick.Marker = tiplot.MajorLogTicks{}, tiplot.MajorLogTicks{}
		p.X.Tick.Label.Rotation = -tiplot.FourtyFiveDegree
		p.X.Tick.Label.XAlign, p.X.Tick.Label.YAlign = text.XLeft, text.YTop
	}

	p.Add(plotter.NewGrid())
	p.Add(ResultPlotters(results)...)

	err = p.Save(vg.Length(*outputSize), vg.Length(*outputSize), *outputFile)
	if err != nil {
		log.Fatalf("saving plot: %s", err)
	}

	pprof.StopCPUProfile()
}
