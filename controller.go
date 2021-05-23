package main

import (
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
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
