package photoperiod

import (
	"math"
	"time"

	"github.com/sgtsquiggs/go-aiprime/models"
)

type Photoperiod interface {
	Schedule() (*models.Schedule, error)
	LunarSchedule() (*models.LunarSchedule, error)
}

func timeToMinutes(t time.Time) int {
	return t.Hour()*60 + t.Minute()
}

func percentToIntensity(p float64) int {
	if p < 0 {
		return 0
	}
	if p > 1.25 {
		return 2000
	}
	if p > 1 {
		return 1000 + int(math.Round(1000*math.Mod(p, 1)/0.25))
	}
	return int(math.Round(1000 * p))
}

func fixIntensityMap(mp map[models.Color]float64) map[models.Color]int {
	mi := make(map[models.Color]int, len(models.ColorValues()))
	for _, color := range models.ColorValues() {
		p, ok := mp[color]
		if ok {
			mi[color] = percentToIntensity(p)
		} else {
			mi[color] = 0
		}
	}
	return mi
}
