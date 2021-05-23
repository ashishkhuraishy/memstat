package main

import (
	"runtime"
)

func LoadStats() *runtime.MemStats {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	return &rtm
}
