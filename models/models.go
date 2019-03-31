//go:generate enumer -json -trimprefix Color -transform snake -type Color -output models_string.go
package models

type Color int

const (
	ColorBlue Color = iota
	ColorGreen
	ColorDeepRed
	ColorWarmWhite
	ColorCoolWhite
)

type Time struct {
	Seconds  int    `json:"seconds"`
	Minutes  int    `json:"minutes"`
	Hours    int    `json:"hours"`
	Day      int    `json:"day"`
	Month    int    `json:"month"`
	Year     int    `json:"year"`
	Location string `json:"location"`
}

type Schedule struct {
	Name  string `json:"name"`
	Ramps []Ramp `json:"ramps"`
}

type Ramp struct {
	Color  Color   `json:"color"`
	Points []Point `json:"points"`
}

type Point struct {
	// Time is time of day in minutes
	Time int `json:"time"`
	// Intensity is a percentage from 0 to 2000
	Intensity int `json:"intensity"`
}

// PointSlice attaches the methods of Interface to []Point, sorting in increasing order.
type PointSlice []Point

func (p PointSlice) Len() int           { return len(p) }
func (p PointSlice) Less(i, j int) bool { return p[i].Time < p[j].Time }
func (p PointSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type LunarSchedule struct {
	Enable bool `json:"enable"`
	Start  int  `json:"start"`
	End    int  `json:"end"`
}

