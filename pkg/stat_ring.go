package pkg

import (
	"container/ring"
)

// StatRing struct stores the data in a ring format,
// which is same as a circular list. New data will be
// added to the tail and the process will continue in
// a cyclic fashion which will help us to render data
// in a moving format
type StatRing struct {
	ring *ring.Ring
}

// newChartRing works as a constructor for the statRing
// struct which initialises a ring inside it with the size
// of `n`.
func newChartRing(n int) *StatRing {
	return &StatRing{
		ring: ring.New(n),
	}
}

// Push will add a data to the current head of the list
// and moves the head to the next position in the list
func (s *StatRing) Push(n uint64) {
	s.ring.Value = n
	s.ring = s.ring.Next()
}

// Data will convert the `uint64` type data inside the
// ring on to a `float64` which then can be fed into
// graphs to render on the terminal
func (s *StatRing) Data() []float64 {
	var res []float64

	// Iterating through all the elements inside the
	// ring and converting them to `float64`
	s.ring.Do(func(i interface{}) {
		if val, ok := i.(uint64); ok {
			res = append(res, float64(val))
		} else {
			// If the conversion fails then set that
			// value as 0
			res = append(res, 0.0)
		}
	})

	return res
}

// Normalised will normalise all the data inside the
// ring, ie, convert all the data inside the ring which
// is in `uint64` to float and make it between 0 & 1.
func (s *StatRing) Normalised() []float64 {
	// finiding the max term inside the ring
	max := s.max()

	// If the max is 0 then there are no items inside
	// the ring, so returning an empty list with the
	// ring size
	if max == 0 {
		return make([]float64, s.ring.Len())
	}

	var res []float64

	// If there are values inisde the ring then itrate
	// over them and divide the value with the `max`
	// element, so that the result will always be 1
	// or less that 1.
	s.ring.Do(func(i interface{}) {
		var pct float64
		if val, ok := i.(uint64); ok {
			pct = float64(val) / float64(max)
		}

		res = append(res, pct)
	})

	return res
}

// `max` will iterate through all the elements inside the
// ring and find the biggest `uint64` in there.
func (s *StatRing) max() uint64 {
	var max uint64

	s.ring.Do(func(i interface{}) {
		if val, ok := i.(uint64); ok && val > max {
			max = val
		}
	})

	return max
}
