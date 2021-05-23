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

func (s *StatRing) Data() []float64 {
	var res []float64

	s.ring.Do(func(i interface{}) {
		if val, ok := i.(uint64); ok {
			res = append(res, float64(val))
		} else {
			res = append(res, 0.0)
		}
	})

	return res
}

func (s *StatRing) Normalised() []float64 {
	max := s.max()

	if max == 0 {
		/// Here
		return make([]float64, s.ring.Len())
	}

	var res []float64

	s.ring.Do(func(i interface{}) {
		var pct float64
		if val, ok := i.(uint64); ok {
			pct = float64(val) / float64(max)
		}

		res = append(res, pct)
	})

	return res
}

func (s *StatRing) max() uint64 {
	var max uint64

	s.ring.Do(func(i interface{}) {
		if val, ok := i.(uint64); ok && val > max {
			max = val
		}
	})

	return max
}
