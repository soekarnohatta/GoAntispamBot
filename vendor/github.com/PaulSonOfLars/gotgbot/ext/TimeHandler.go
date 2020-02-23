package ext

import (
	"time"
)

func GetDelay(timeToProcess int) float64 {
	timeSpan := float64(time.Now().UTC().Nanosecond() * 1000000) - float64(timeToProcess)
	return timeSpan/100000000000000
}

func GetJeda(jeda time.Time) float64 {
	timeSpan := time.Since(jeda)
	return timeSpan.Seconds()
}
