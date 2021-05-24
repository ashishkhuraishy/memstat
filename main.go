package main

import (
	"log"
	"time"

	"github.com/ashishkhuraishy/memstat/pkg"
	ui "github.com/gizak/termui/v3"
)

func main() {
	Run()
}

// The run function is responsible for initialising and rendering
// the ui on the screen and updating it on each second.
func Run() {
	// Initialising the ui package
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialise terminal : %v", err)
	}
	defer ui.Close()

	// Initialising a new controller which stores the data
	// for rendering and resizing our ui
	controller := pkg.NewController()

	// `events` will be a channel which contains data for
	// each time a UI event occurs (like mouse click, keyboard
	// typed, resized). Here we will listean to the values from
	// the channel to quit and redraw the ui each time  user
	// resizes the terminal
	events := ui.PollEvents()
	// `tick` is also a channel which will trigger each second
	// so that we can use that to update our on screen data and
	// display new data every second
	tick := time.Tick(time.Second)

	// Starting an infinte loop which will listean for
	// events and time tiker
	for {
		select {
		// This will listean for every events that occur
		// related to the UI
		case event := <-events:
			switch event.Type {
			case ui.KeyboardEvent:
				// If the user clicks <CTRL+c> then the
				// loop will break and the program will
				// be executed
				if event.ID == "<C-c>" {
					return
				}
			// When the user resizes the terminal, we are
			// triggering the controller to redraw the UI,
			// so that it'll be responsive
			case ui.ResizeEvent:
				controller.Resize()
			}
		// Here the ticker will trigger every second and
		// we are getting the current status from the
		// system , then using that data to render a new
		// UI state.
		case <-tick:
			// Loading the status
			stat := pkg.LoadStats()
			// rendering the new state with new data
			controller.Render(stat)
		}
	}

}
