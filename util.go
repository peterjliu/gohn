package gohn

import (
	"math"
	"math/rand"
	"time"
)

func ExpWithMax(mean, maxi float64) time.Duration {
	return time.Duration(math.Min(maxi, mean*rand.ExpFloat64()))
}
