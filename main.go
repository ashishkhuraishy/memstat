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

func Run() {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialise terminal : %v", err)
	}
	defer ui.Close()

	controller := pkg.NewController()

	events := ui.PollEvents()
	tick := time.Tick(time.Second)

	for {
		select {
		case event := <-events:
			switch event.Type {
			case ui.KeyboardEvent:
				if event.ID == "<C-c>" {
					return
				}
			case ui.ResizeEvent:
				controller.Resize()
			}
		case <-tick:
			stat := pkg.LoadStats()
			controller.Render(stat)
		}
	}

}
