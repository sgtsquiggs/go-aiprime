package main

import (
	"context"
	"flag"
	"log"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/nathan-osman/go-sunrise"

	"github.com/sgtsquiggs/go-aiprime/client"
	"github.com/sgtsquiggs/go-aiprime/config"
	"github.com/sgtsquiggs/go-aiprime/models"
	"github.com/sgtsquiggs/go-aiprime/photoperiod"
)

var (
	ConfigFlag = flag.String("config", "", "config path")
	Config     = &config.Config{}
	Client     *client.Client
)

func updateTime(ctx context.Context) {
	now := time.Now().In(Config.Location)

	timeRequest := &client.TimeRequest{
		Seconds:  now.Second(),
		Minutes:  now.Minute(),
		Hours:    now.Hour(),
		Day:      now.Day(),
		Month:    int(now.Month()),
		Year:     now.Year(),
		Location: Config.Timezone,
	}

	_, err := Client.Time.Update(ctx, timeRequest)
	if err != nil {
		log.Printf("error updating time: %v", err)
	}
}

func updateSchedule(ctx context.Context) {
	var (
		riseInt = make(map[string]float64)
		midInt  = make(map[string]float64)
		setInt  = make(map[string]float64)
	)
	for _, color := range models.ColorValues() {
		riseInt[color] = 0.05
		midInt[color] = 0.80
		setInt[color] = 0
	}

	setInt[models.ColorDeepRed] = 0.02
	setInt[models.ColorWarmWhite] = 0.02
	setInt[models.ColorBlue] = 0.02

	period := &photoperiod.CyclodialPhotoperiod{
		SunriseIntensity: riseInt,
		MiddayIntensity:  midInt,
		SunsetIntensity:  setInt,
		Resolution:       30,
		Config:           *Config,
	}

	schedule, err := period.Schedule()
	if err != nil {
		log.Fatal(err)
	}

	// hack for moonlight until I add a better way
	now := time.Now().In(Config.Location)
	_, set := sunrise.SunriseSunset(Config.Latitude, Config.Longitude, now.Year(), now.Month(), now.Day())
	setTime := int(set.In(Config.Location).Hour()*60 + set.In(Config.Location).Minute())

	twilightTime := setTime + 60

	var points []models.Point
	for i, ramp := range schedule.Ramps {
		switch ramp.Color {
		case models.ColorBlue:
			points = []models.Point{
				{Time: twilightTime, Intensity: 20},
				{Time: 0, Intensity: 20},
				{Time: 180, Intensity: 0},
			}
			points = append(ramp.Points, points...)
		default:
			points = []models.Point{
				{Time: twilightTime, Intensity: 0},
				{Time: 0, Intensity: 0},
				{Time: 180, Intensity: 0},
			}
			points = append(ramp.Points, points...)
		}
		sort.Sort(models.PointSlice(points))
		schedule.Ramps[i].Points = points
	}

	scheduleRequest := &client.ScheduleRequest{Name: schedule.Name, Ramps: schedule.Ramps}

	lunarRequest := &client.LunarScheduleRequest{
		Enable: true,
		Start:  twilightTime,
		End:    180,
	}

	_, err = Client.Schedule.Update(ctx, scheduleRequest)
	if err != nil {
		log.Printf("error updating schedule: %v", err)
	}

	_, err = Client.Schedule.UpdateLunar(ctx, lunarRequest)
	if err != nil {
		log.Printf("error updating lunar schedule: %v", err)
	}
}

func main() {
	flag.Parse()

	if ConfigFlag == nil || *ConfigFlag == "" {
		log.Fatal("missing config path")
	}

	_, err := toml.DecodeFile(*ConfigFlag, Config)
	if err != nil {
		log.Fatalf("invalid configuration: %v", err)
	}

	Config.Location, err = time.LoadLocation(Config.Timezone)
	if err != nil {
		log.Fatalf("invalid timezone %v: %v", Config.Timezone, err)
	}

	Client, err = client.New(client.SetBaseURL(Config.DeviceURL))
	if err != nil {
		log.Fatalf("invalid device url: %v", err)
	}

	deadline := time.Now().Add(15 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), deadline)

	updateTime(ctx)
	updateSchedule(ctx)

	cancel()
}
