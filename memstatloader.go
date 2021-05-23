package main

import (
	"net/http"
	"runtime"
)

type memStatsLoader struct {
	url    string
	client *http.Client
}

func (p *memStatsLoader) Load() *runtime.MemStats {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	return &rtm
}
