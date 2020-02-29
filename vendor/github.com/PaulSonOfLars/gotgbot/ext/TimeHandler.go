package ext

import (
	"time"
)

func GetJeda(jeda time.Time) float64 {
	timeSpan := time.Since(jeda)
	return timeSpan.Seconds()
}
