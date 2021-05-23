package main

import (
	"container/ring"
	"runtime"
)

type Controller interface {
	Render(*runtime.MemStats)
	Resize()
}

type StatRing struct {
	ring *ring.Ring
}

func NewChartRing(n int) *StatRing {
	return &StatRing{
		ring: ring.New(n),
	}
}

func (s *StatRing) Push(n uint64) {
	s.ring.Value = n
	s.ring = s.ring.Next()
}
