package main

import "fmt"

// Converting size in bytes to a human readble format
func byteToBinary(b uint64) string {
	const uint = 1024

	if b < uint {
		return fmt.Sprintf("%d B", b)
	}

	div, exp := uint64(uint), 0
	for n := b / uint; n >= uint; n /= uint {
		div *= uint
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "KMGTPE"[exp])
}
