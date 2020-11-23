package photoperiod

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/nathan-osman/go-sunrise"

	"github.com/sgtsquiggs/go-aiprime/config"
	"github.com/sgtsquiggs/go-aiprime/models"
)

type CyclodialPhotoperiod struct {
	SunriseIntensity map[string]float64
	MiddayIntensity  map[string]float64
	SunsetIntensity  map[string]float64
	MinimumPeriod    time.Duration
	MaximumPeriod    time.Duration
	Resolution       int
	Config           config.Config

	riseTime   int
	noonTime   int
	setTime    int
	calcConfig config.Config
}

var _ Photoperiod = &CyclodialPhotoperiod{}

func (p *CyclodialPhotoperiod) Schedule() (*models.Schedule, error) {
	if p.Resolution <= 0 {
		return nil, errors.New("resolution must be greater than zero")
	}

	p.calculateTimes()
	riseMap := fixIntensityMap(p.SunriseIntensity)
	midMap := fixIntensityMap(p.MiddayIntensity)
	setMap := fixIntensityMap(p.SunsetIntensity)

	ramps := make([]models.Ramp, 0, len(models.ColorValues()))
	for _, color := range models.ColorValues() {
		points, err := cycloidCurvePoints(
			p.riseTime,
			p.noonTime,
			p.setTime,
			riseMap[color],
			midMap[color],
			setMap[color],
			p.Resolution)
		if err != nil {
			return nil, fmt.Errorf("error generating ramp %s: %w", color, err)
		}
		ramp := models.Ramp{Color: color, Points: points}
		ramps = append(ramps, ramp)
	}

	return &models.Schedule{Ramps: ramps}, nil
}

func (p *CyclodialPhotoperiod) LunarSchedule() (*models.LunarSchedule, error) {
	p.calculateTimes()
	return &models.LunarSchedule{Enable: true, Start: p.setTime, End: p.riseTime}, nil
}

func (p *CyclodialPhotoperiod) calculateTimes() {
	if p.Config == p.calcConfig {
		return
	}
	now := time.Now().In(p.Config.Location)
	rise, set := sunrise.SunriseSunset(p.Config.Latitude, p.Config.Longitude, now.Year(), now.Month(), now.Day())
	noon := sunrise.JulianDayToTime(sunrise.MeanSolarNoon(p.Config.Longitude, now.Year(), now.Month(), now.Day()))
	p.riseTime = timeToMinutes(rise.In(p.Config.Location))
	p.noonTime = timeToMinutes(noon.In(p.Config.Location))
	p.setTime = timeToMinutes(set.In(p.Config.Location))
	p.calcConfig = p.Config
}

// x(t) = a(t - sin(t)), y(t) = a(1 - cos(t))
func cycloidCurvePoints(
	startTime int,
	midTime int,
	endTime int,
	startIntensity int,
	midIntensity int,
	endIntensity int,
	resolution int) ([]models.Point, error) {

	// sanity checks

	if resolution <= 0 {
		return nil, errors.New("resolution must be greater than zero")
	}
	if startIntensity > midIntensity {
		return nil, errors.New("midIntensity must be larger than startIntensity")
	}
	if endIntensity > midIntensity {
		return nil, errors.New("midIntensity must be larger than endIntensity")
	}
	if startTime >= midTime {
		return nil, errors.New("midTime must be after startTime")
	}
	if midTime >= endTime {
		return nil, errors.New("midTime must be before endTime")
	}

	// increment resolution by two to exclude start/end points
	resolution += 2
	if int(math.Mod(float64(resolution), 2)) > 0 {
		// decrement resolution by one to remove midpoint
		resolution -= 1
	}
	median := float64(resolution) / 2

	aX := float64(endTime - startTime)

	var t, x, y, aY float64
	var time, intensity int
	var points []models.Point
	var prevPoint *models.Point
	for i := 0; i <= resolution; i++ {

		t = 2 * math.Pi * float64(i) / float64(resolution)
		x = t - math.Sin(t)
		y = 1 - math.Cos(t)

		time = startTime + int(math.Round(aX*x/(2*math.Pi)))

		if float64(i) <= median {
			aY = float64(midIntensity - startIntensity)
			intensity = startIntensity + int(math.Round(aY*y/2))
		} else {
			aY = float64(midIntensity - endIntensity)
			intensity = endIntensity + int(math.Round(aY*y/2))
		}

		// check for duplicate time
		if prevPoint != nil && prevPoint.Time == time {
			// approximate...
			prevPoint.Intensity = (prevPoint.Intensity + intensity) / 2
			continue
		}

		point := models.Point{
			Time:      time,
			Intensity: intensity,
		}
		points = append(points, point)
		prevPoint = &point
	}

	return points, nil
}
