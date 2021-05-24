package pkg

import (
	"fmt"
	"runtime"
)

// LoadStats will read the current memory stats and
// write it on to a variable and return it
func LoadStats() *runtime.MemStats {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	return &rtm
}

// byteToBinary will take in a `uint64` in a byte
// size format and convert it into common human
// readable memory size formats
func byteToBinary(b uint64) string {
	// Setting the basic memory size as unitSize
	const unit = 1024

	// checking if the inputted size is less than
	// the unit size. If it is then we can return
	// it as is and is in byte size (500 -> 500B)
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	// If the size is greater thab the unit size we
	// can divide it up until it reaches a standard
	// size. (2048 -> 2KB)
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
