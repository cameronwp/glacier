package ui

import (
	"math"

	"github.com/gizak/termui"
)

// Render is a blocking call that generates the UI.
func Render() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	vaultPar := termui.NewPar("Valve Documentary\nus-east-1\npath/to/target\n10 files")
	vaultPar.BorderLabel = "Batch Info"
	vaultPar.Height = 6
	vaultPar.TextFgColor = termui.ColorWhite
	vaultPar.BorderFg = termui.ColorCyan

	totalGauge := termui.NewGauge()
	totalGauge.Percent = 50
	totalGauge.Height = 3
	totalGauge.BorderLabel = "Overall Progress"
	totalGauge.Label = "{{percent}}% (100MBs free)"
	totalGauge.BarColor = termui.ColorRed
	totalGauge.BorderFg = termui.ColorWhite
	totalGauge.BorderLabelFg = termui.ColorCyan

	shift := float64(0)
	sinps := func(x float64) []float64 {
		n := 110
		ps := make([]float64, n)
		for i := range ps {
			ps[i] = 1 + math.Sin(float64(i)/5+x)
		}
		return ps
	}

	rateGraph := termui.NewLineChart()
	rateGraph.BorderLabel = "Upload Rate" // update if using KB, MB, etc
	rateGraph.Data = sinps(shift)
	rateGraph.DataLabels = make([]string, 220)
	rateGraph.Height = 9
	rateGraph.AxesColor = termui.ColorWhite
	rateGraph.LineColor = termui.ColorGreen | termui.AttrBold

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, vaultPar, totalGauge),
			termui.NewCol(6, 0, rateGraph),
		),
	)

	// display controls at the bottom

	termui.Body.Align()
	termui.Render(termui.Body)

	// kill it
	termui.Handle("/sys/kbd/q", func(termui.Event) {
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/C-c", func(termui.Event) {
		termui.StopLoop()
	})

	termui.Handle("/sys/kbd/l", func(termui.Event) {
		shift = shift + 0.25
		rateGraph.Data = sinps(shift)
		termui.Render(termui.Body)
	})
	termui.Handle("/sys/kbd/h", func(termui.Event) {
		shift = shift - 0.25
		rateGraph.Data = sinps(shift)
		termui.Render(termui.Body)
	})

	termui.Loop()
}
