package main

import (
	"cmp"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime/pprof"
	"slices"
	"strings"

	"github.com/elordis/ti_drive_plot/ti"
	"github.com/elordis/ti_drive_plot/tiplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
)

func Abbreviate(s string, abbrevLimit int, keepLast bool) string {
	abbrev := ""
	parts := strings.Split(s, " ")

	max := len(parts)
	if keepLast {
		max--
	}

	for i := 0; i < max; i++ {
		abbrev += parts[i][:abbrevLimit]
	}

	if keepLast {
		abbrev += " " + parts[max]
	}

	return abbrev
}

func ApproximatelyEqual(a, b, delta float64) bool {
	return math.Abs(a-b) < delta*(a+b)/2
}

type PlotNailer interface {
	plot.Plotter
	plot.Thumbnailer
}

type Result struct {
	MinShip *ti.Ship
	MaxShip *ti.Ship
}

func (r *Result) Name() string {
	evName := "None"
	if r.MinShip.DriveAssembly.EVUtility != nil {
		evName = r.MinShip.DriveAssembly.EVUtility.FriendlyName
	}

	thrustName := "None"
	if r.MinShip.DriveAssembly.ThrustUtility != nil {
		thrustName = r.MinShip.DriveAssembly.ThrustUtility.FriendlyName
	}

	return fmt.Sprintf(
		"%s %s %s %s %s",
		Abbreviate(r.MinShip.DriveAssembly.Drive.FriendlyName, 2, true),      //nolint:gomnd
		Abbreviate(r.MinShip.DriveAssembly.PowerPlant.FriendlyName, 2, true), //nolint:gomnd
		Abbreviate(r.MinShip.DriveAssembly.Radiator.FriendlyName, 2, false),  //nolint:gomnd
		Abbreviate(evName, 2, false),                                         //nolint:gomnd
		Abbreviate(thrustName, 2, false),                                     //nolint:gomnd
	)
}

func (r *Result) String() string {
	return fmt.Sprintf(
		"%s - %drp Min: %.0ft %.0ft/d Max: %.0ft %.0ft/d",
		r.Name(),
		r.MinShip.DriveAssembly.EffectiveRPCost(),
		r.MinShip.PayloadMassT,
		r.MinShip.PayloadTPerMiningDay(),
		r.MaxShip.PayloadMassT,
		r.MaxShip.PayloadTPerMiningDay(),
	)
}

const Nearness = 0.05

func (r *Result) IsBetterThan(s *Result) bool { //nolint:gocognit
	if r.MinShip.DriveAssembly.Drive != s.MinShip.DriveAssembly.Drive {
		return false
	}

	minorBetter, minorWorse, majorBetter, majorWorse := 0, 0, 0, 0

	if ApproximatelyEqual(r.MinShip.PayloadMassT, s.MinShip.PayloadMassT, Nearness) { //nolint:nestif
		if r.MinShip.PayloadMassT < s.MinShip.PayloadMassT {
			minorBetter++
		} else {
			minorWorse++
		}
	} else {
		if r.MinShip.PayloadMassT < s.MinShip.PayloadMassT {
			majorBetter++
		} else {
			majorWorse++
		}
	}

	if ApproximatelyEqual(r.MinShip.PayloadTPerMiningDay(), s.MinShip.PayloadTPerMiningDay(), Nearness) { //nolint:nestif
		if r.MinShip.PayloadTPerMiningDay() > s.MinShip.PayloadTPerMiningDay() {
			minorBetter++
		} else {
			minorWorse++
		}
	} else {
		if r.MinShip.PayloadTPerMiningDay() > s.MinShip.PayloadTPerMiningDay() {
			majorBetter++
		} else {
			majorWorse++
		}
	}

	if ApproximatelyEqual(r.MaxShip.PayloadMassT, s.MaxShip.PayloadMassT, Nearness) { //nolint:nestif
		if r.MaxShip.PayloadMassT > s.MaxShip.PayloadMassT {
			minorBetter++
		} else {
			minorWorse++
		}
	} else {
		if r.MaxShip.PayloadMassT > s.MaxShip.PayloadMassT {
			majorBetter++
		} else {
			majorWorse++
		}
	}

	if ApproximatelyEqual(r.MaxShip.PayloadTPerMiningDay(), s.MaxShip.PayloadTPerMiningDay(), Nearness) { //nolint:nestif
		if r.MaxShip.PayloadTPerMiningDay() > s.MaxShip.PayloadTPerMiningDay() {
			minorBetter++
		} else {
			minorWorse++
		}
	} else {
		if r.MaxShip.PayloadTPerMiningDay() > s.MaxShip.PayloadTPerMiningDay() {
			majorBetter++
		} else {
			majorWorse++
		}
	}

	if majorBetter == 0 && majorWorse == 0 {
		return minorBetter > minorWorse
	}

	if majorBetter > 0 && majorWorse > 0 {
		return false
	}

	return majorBetter > majorWorse
}

func (r *Result) IsWorseThan(s *Result) bool {
	return s.IsBetterThan(r)
}

func (r *Result) PlotLine() ([]PlotNailer, error) {
	elems := make([]PlotNailer, 2) //nolint:gomnd
	line, points, err := plotter.NewLinePoints(
		plotter.XYs{
			{
				X: r.MinShip.PayloadMassT,
				Y: r.MinShip.PayloadTPerMiningDay(),
			},
			{
				X: r.MaxShip.PayloadMassT,
				Y: r.MaxShip.PayloadTPerMiningDay(),
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("generating line for result: %w", err)
	}

	line.LineStyle, points.GlyphStyle = tiplot.DriveStyle(r.MinShip.DriveAssembly.Drive)
	elems[0], elems[1] = line, points

	return elems, nil
}

func CmpMaxPayload(a, b *Result) int {
	return cmp.Compare(a.MaxShip.PayloadMassT, b.MaxShip.PayloadMassT)
}

func CmpMaxPtmpd(a, b *Result) int {
	return cmp.Compare(a.MaxShip.PayloadTPerMiningDay(), b.MaxShip.PayloadTPerMiningDay())
}

func main() {
	templatePath := flag.String("t", "", "folder containing template JSON files")
	constraintsFile := flag.String("c", "", "JSON-encoded constraints file")
	minDV := flag.Float64("dv", 10, "minimum delta-V to target")                            //nolint:gomnd
	minCruiseAccel := flag.Float64("cra", 0.001, "minimum allowed cruise acceleration")     //nolint:gomnd
	minCombatAccel := flag.Float64("coa", 0.001, "minimum allowed combat acceleration")     //nolint:gomnd
	maxRPCost := flag.Int("rp", 10000000, "maximum allowed RP cost for drive assembly")     //nolint:gomnd
	minPtpmd := flag.Float64("ptpmd", 200, "minimum allowed payload tonns per mining days") //nolint:gomnd
	minPayload := flag.Float64("pl", 500, "minimum allowed payload tonns per mining days")  //nolint:gomnd
	outputFile := flag.String("o", "output.png", "file to write output to")
	cpuprofile := flag.String("cpuprofile", "", "write cpu profile to file")

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

	results := make([]*Result, 0, len(engine.DriveTemplates))
	engine.ForAllDriveAssemblies(false, func(d *ti.DriveAssembly) {
		rpCost := d.EffectiveRPCost()
		if rpCost > *maxRPCost {
			return
		}
		maxPayload := d.MaxPayloadConstrained(*minDV, *minCruiseAccel, *minCombatAccel)
		if maxPayload <= *minPayload {
			return
		}
		maxShip := &ti.Ship{
			DriveAssembly: d,
			Tanks:         d.TanksForPayloadDV(maxPayload, *minDV),
			PayloadMassT:  maxPayload,
		}
		if maxShip.PayloadTPerMiningDay() < *minPtpmd {
			return
		}
		minShip := &ti.Ship{
			DriveAssembly: d,
			Tanks:         d.TanksForPayloadDV(maxPayload, *minDV),
			PayloadMassT:  maxPayload,
		}
		low := *minPayload
		high := maxShip.PayloadMassT
		for high-low > 1 {
			minShip.PayloadMassT = (high + low) / 2 //nolint:gomnd
			minShip.Tanks = d.TanksForPayloadDV(minShip.PayloadMassT, *minDV)

			if minShip.PayloadTPerMiningDay() < *minPtpmd {
				low = minShip.PayloadMassT
			} else {
				high = minShip.PayloadMassT
			}
		}
		r := Result{MinShip: minShip, MaxShip: maxShip}

		if slices.ContainsFunc(results, r.IsWorseThan) {
			return
		}
		results = slices.DeleteFunc(results, r.IsBetterThan)
		results = append(results, &r)
	})

	p := plot.New()
	p.Title.Text = "Payload per Mining Day / Payload\n"
	p.Title.Text += fmt.Sprintf(
		"Mining: %.2f*%v+(%s)\n",
		1+engine.SimulationConstraints.MiningBonus,
		engine.SimulationConstraints.MiningBodies,
		engine.SimulationConstraints.MiningFlat,
	)
	p.Title.Text += fmt.Sprintf(
		"DV: %.0f Kps, Cruise: %.3f G, Combat: %.3f G, RP: %d, PTpMD: %.0f T/d, Payload: %.0f T",
		*minDV,
		*minCruiseAccel,
		*minCombatAccel,
		*maxRPCost,
		*minPtpmd,
		*minPayload,
	)
	p.X.Label.Text = "Payload (T)"
	p.Y.Label.Text = "Payload per Mining Day (T/day)"
	p.X.Min, p.X.Max = *minPayload, slices.MaxFunc(results, CmpMaxPayload).MaxShip.PayloadMassT
	p.Y.Min, p.Y.Max = *minPtpmd, slices.MaxFunc(results, CmpMaxPtmpd).MaxShip.PayloadTPerMiningDay()
	// p.X.Scale, p.Y.Scale = plot.LogScale{}, plot.LogScale{}
	// p.X.Tick.Marker, p.X.Tick.Marker = plot.LogTicks{Prec: -1}, plot.LogTicks{Prec: -1}
	p.Legend = plot.NewLegend()
	p.Legend.ThumbnailWidth = 50
	p.Add(plotter.NewGrid())

	for _, r := range results {
		elems, err := r.PlotLine()
		if err != nil {
			log.Fatalf("generating plot data for '%s': %s", r, err)
		}

		elemsThumb := make([]plot.Thumbnailer, len(elems))

		for i, e := range elems {
			p.Add(e)
			elemsThumb[i] = e
		}

		p.Legend.Add(r.Name(), elemsThumb...)
		log.Println(r)
	}

	err = p.Save(1000, 1000, *outputFile) //nolint:gomnd
	if err != nil {
		log.Fatalf("saving plot: %s", err)
	}

	pprof.StopCPUProfile()
}
