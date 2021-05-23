package pkg

import (
	"fmt"
	"runtime"
)

func LoadStats() *runtime.MemStats {
	var rtm runtime.MemStats
	runtime.ReadMemStats(&rtm)
	return &rtm
}

// Converting size in bytes to a human readble format
func byteToBinary(b uint64) string {
	const unit = 1024

	if b < unit {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
