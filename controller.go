package main

import (
	"fmt"
	"runtime"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	terminalWidth     = 120
	heapAllocBarCount = 6
)

type controller struct {
	Grid *ui.Grid

	HeapObjectsSparkLine     *widgets.Sparkline
	HeapObjectSparkLineGroup *widgets.SparklineGroup
	HeapObjectsData          *StatRing

	SysText        *widgets.Paragraph
	GCCCPUFraction *widgets.Gauge

	HeapAllocBarChart     *widgets.BarChart
	HeapAllocBarChartData *StatRing

	HeapPie *widgets.PieChart
}

func newController() *controller {
	ctl := &controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkLine:     widgets.NewSparkline(),
		HeapObjectSparkLineGroup: widgets.NewSparklineGroup(),
		HeapObjectsData:          NewChartRing(terminalWidth),

		SysText:        widgets.NewParagraph(),
		GCCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChart:     widgets.NewBarChart(),
		HeapAllocBarChartData: NewChartRing(heapAllocBarCount),

		HeapPie: widgets.NewPieChart(),
	}

	ctl.initUI()

	return ctl
}

func (c *controller) Render(data *runtime.MemStats) {
	c.HeapObjectsData.Push(data.HeapObjects)
	c.HeapObjectsSparkLine.Data = c.HeapObjectsData.Normalised()
	c.HeapObjectSparkLineGroup.Title = fmt.Sprintf("HeapObjects, live heap object count: %d", data.HeapObjects)

	c.SysText.Text = fmt.Sprint(byteToBinary(data.Sys))

	fNormalize := func() int {
		f := data.GCCPUFraction
		if f < 0.01 && f > 0 {
			for f < 1 {
				f = f * 10.0
			}
		}

		return int(f)
	}

	c.GCCCPUFraction.Percent = fNormalize()
	c.GCCCPUFraction.Label = fmt.Sprintf("%.2f%%", data.GCCPUFraction*100)

	c.HeapAllocBarChartData.Push(data.HeapAlloc)
	c.HeapAllocBarChart.Data = c.HeapAllocBarChartData.Data()
	c.HeapAllocBarChart.Labels = nil
	for _, v := range c.HeapAllocBarChart.Data {
		c.HeapAllocBarChart.Labels = append(c.HeapAllocBarChart.Labels, byteToBinary(uint64(v)))
	}

	c.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	ui.Render(c.Grid)
}

func (c *controller) Resize() {
	c.resize()
	ui.Render(c.Grid)
}

func (c *controller) initUI() {
	c.resize()

	c.HeapObjectsSparkLine.LineColor = ui.Color(89)
	c.HeapObjectSparkLineGroup = widgets.NewSparklineGroup(c.HeapObjectsSparkLine)

	c.SysText.Text = "Sys, the total bytes of memory obtained from OS"
	c.SysText.PaddingLeft = 25
	c.SysText.PaddingTop = 1

	c.HeapAllocBarChart.BarGap = 2
	c.HeapAllocBarChart.BarWidth = 8
	c.HeapAllocBarChart.Title = "HeapAlloc, bytes of allocated heap objects"
	c.HeapAllocBarChart.NumFormatter = func(f float64) string { return "" }

	c.GCCCPUFraction.Title = "GCCCPUFraction 0%~100%"
	c.GCCCPUFraction.BarColor = ui.Color(50)

	c.HeapPie.Title = "HeapInUse vs HeapIdle"
	c.HeapPie.LabelFormatter = func(dataIndex int, _ float64) string { return []string{"Idle", "Inuse"}[dataIndex] }

	c.Grid.Set(
		ui.NewRow(0.2, c.HeapObjectSparkLineGroup),
		ui.NewRow(0.8,
			ui.NewCol(0.5,
				ui.NewRow(0.2, c.SysText),
				ui.NewRow(0.2, c.GCCCPUFraction),
				ui.NewRow(0.6, c.HeapAllocBarChart),
			),
			ui.NewCol(0.5, c.HeapPie),
		),
	)
}

func (c *controller) resize() {
	_, h := ui.TerminalDimensions()
	c.Grid.SetRect(0, 0, terminalWidth, h)
}
