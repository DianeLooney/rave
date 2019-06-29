package waves

import (
	"math"
)

func Triangle(offset float64) float64 {
	div := int(offset) % 2
	rem := offset - float64(int(offset))
	if div == 0 {
		return -1 + 2*rem
	}

	return 1 - 2*rem
}

func Sin(offset float64) float64 {
	return math.Sin(2 * math.Pi * offset)
}

func Square(offset float64) float64 {
	if int(offset)%2 == 0 {
		return 1
	}

	return -1
}
