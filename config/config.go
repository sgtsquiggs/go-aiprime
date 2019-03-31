package config

import "time"

type Config struct {
	Latitude  float64 `toml:"latitude"`
	Longitude float64 `toml:"longitude"`
	DeviceURL string  `toml:"device_url"`
	Timezone  string  `toml:"timezone"`

	Location *time.Location
}

