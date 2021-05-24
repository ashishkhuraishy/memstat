package pkg

import (
	"fmt"
	"runtime"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

const (
	// Number of bars in bar chart
	heapAllocBarCount = 6
)

// The controller interface is the core interface which
// will contain the basic functions required to use on the
// UI with our custom data
type Controller interface {
	// Render will take in a `memStats` data which then will
	// be converted to UI renderable format
	Render(*runtime.MemStats)
	// Resize method will reset all the size parameters provided
	// to the ui and redraw them with new values with respect to
	// the current terminal size
	Resize()
}

// `controller` is a dType created to implemet the `Controller`
// interface which also stores all the required data to render
// our UI
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

// NewController acts as a constructor for the `controller`
// which will initialise all the values and initialise the
// UI to render
func NewController() *controller {
	// getting the current terminal width to initialise
	// all the widgets
	terminalwidth, _ := ui.TerminalDimensions()

	// initialising all the elements inside the `controller`
	ctl := &controller{
		Grid: ui.NewGrid(),

		HeapObjectsSparkLine:     widgets.NewSparkline(),
		HeapObjectSparkLineGroup: widgets.NewSparklineGroup(),
		HeapObjectsData:          newChartRing(terminalwidth),

		SysText:        widgets.NewParagraph(),
		GCCCPUFraction: widgets.NewGauge(),

		HeapAllocBarChart:     widgets.NewBarChart(),
		HeapAllocBarChartData: newChartRing(heapAllocBarCount),

		HeapPie: widgets.NewPieChart(),
	}

	// After initialising all the elements this will give initial
	// data to the controller and render the UI
	ctl.initUI()

	return ctl
}

// Render will take in the current data of memory stats and then
// extarct each value, then normalise it into a form in which it
// can be fed as the data for the UI.
func (c *controller) Render(data *runtime.MemStats) {
	// Adding the new data on to the previously initialised ring
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

	// Adding the new data on to the previously initialised ring
	c.HeapAllocBarChartData.Push(data.HeapAlloc)
	c.HeapAllocBarChart.Data = c.HeapAllocBarChartData.Data()
	c.HeapAllocBarChart.Labels = nil
	for _, v := range c.HeapAllocBarChart.Data {
		c.HeapAllocBarChart.Labels = append(c.HeapAllocBarChart.Labels, byteToBinary(uint64(v)))
	}

	c.HeapPie.Data = []float64{float64(data.HeapIdle), float64(data.HeapInuse)}

	// After setting up all the required data on to the widgets inside
	// the controller we will then proceed to rendner that UI on to the
	// terminal. The `Grid` encapsulates all the other widgets so render
	// that will return all the other as its children
	ui.Render(c.Grid)
}

// The `Resize` method will calculate the current size of the terminal
// and redraw the terminal with new bounds which will make the UI responsive
func (c *controller) Resize() {
	c.resize()
	ui.Render(c.Grid)
}

// This function will initialise elements iniside the `controller`
// with required static data and positioning each element in the
// screen wrt to one another
func (c *controller) initUI() {
	// Setting the boundary
	c.resize()

	// Initialising all the elements with data / titles and
	// colours and adding styling to each element
	c.HeapObjectsSparkLine.LineColor = ui.Color(89)
	c.HeapObjectSparkLineGroup = widgets.NewSparklineGroup(c.HeapObjectsSparkLine)

	c.SysText.Title = "Sys, the total bytes of memory obtained from OS"
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

	// After initialising every element, this will position the elements
	// on the screen. To make the grid responsive we are using relative
	// size where the height and width can be expressed as factors same
	// as in `css-flexbox`. The total size of height and weight is from
	// 0 - 1, we can assign each element a relative weight for their value
	// which will then be rendered accordingly
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

// This method will find the current height and width of the terminal
// and set that as the bound for the `Grid` with encapsulates all the
// widgets rendered in the terminal
func (c *controller) resize() {
	w, h := ui.TerminalDimensions()
	c.Grid.SetRect(0, 0, w, h)
}
