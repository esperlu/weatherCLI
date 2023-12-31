package utils

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/esperlu/weatherCLI/lang"
)

var Help = `
Usage: $ weather BRU paris new-york

Default station : Brussels

option      Description                    Defaut
--------------------------------------------------     
-l  string   Language code e.g. de or es         en
-d  int      Number of day forecast               5
-t  int      Warning threshold for rain (%)      50
-f  bool     Print tomorrow's forecast        false
-v  bool     Print version                    false
`

// Extract file name from path
func Filename(path string) string {
	pos := strings.LastIndex(path, "/")
	return path[pos+1:]
}

// Format time from 03:04 PM to 15:04
func FormatPMtime(PMtime string) string {
	t, err := time.Parse("03:04 PM", PMtime)
	if err != nil {
		panic(err)
	}
	return t.Format("15:04")
}

// Print program name and runtime version
func PrintProgName(ver string) string {
	return fmt.Sprintf(
		"%s %s compiled with %s",
		Filename(os.Args[0]),
		ver,
		runtime.Version(),
	)
}

// Couleurs
const (
	CReset    = "\033[0m"
	CRed      = "\033[31m"
	CGreen    = "\033[32m"
	UpArrow   = '\u2197'
	DownArrow = '\u2198'
)

// JSON struct
type Weather struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
	Location struct {
		Name    string  `json:"name"`
		Country string  `json:"country"`
		Lat     float32 `json:"lat"`
		Long    float32 `json:"lon"`
	} `json:"location"`
	Current struct {
		LastUpdate    string  `json:"last_updated"`
		Temp          float32 `json:"temp_c"`
		Wind          float32 `json:"wind_kph"`
		WindDir       string  `json:"wind_dir"`
		WindDegree    int     `json:"wind_degree"`
		Qnh           float32 `json:"pressure_mb"`
		FeelsLike     float32 `json:"feelslike_c"`
		Precipitation float32 `json:"precip_mm"`
		Humidity      int     `json:"humidity"`
		Visibility    float32 `json:"vis_km"`
		Uv            float32 `json:"uv"`
		Condition     struct {
			Text string `json:"text"`
		} `json:"condition"`
	} `json:"current"`

	Forecast struct {
		ForecastDay []struct {
			Date      string `json:"date"`
			DateEpoch int    `json:"date_epoch"`
			Day       struct {
				MinTemp       float32 `json:"mintemp_c"`
				MaxTemp       float32 `json:"maxtemp_c"`
				MaxWind       float32 `json:"maxwind_kph"`
				AvgHumidity   float32 `json:"avghumidity"`
				ChanceOfRain  int     `json:"daily_chance_of_rain"`
				Precipitation float32 `json:"totalprecip_mm"`
				Condition     struct {
					Text string `json:"text"`
				} `json:"condition"`
			} `json:"day"`
			Astro struct {
				SunRise string `json:"sunrise"`
				SunSet  string `json:"sunset"`
			} `json:"astro"`
			Hour []Forecast
			// } `json:"hour"`
		} `json:"forecastday"`
	} `json:"forecast"`
}

// JSON struct for forecast func
type Forecast struct {
	TimeEpoch int    `json:"time_epoch"`
	Time      string `json:"time"`
	Condition struct {
		Text string `json:"text"`
	} `json:"condition"`
	Temp          float32 `json:"temp_c"`
	DewPoint      float32 `json:"dewpoint_c"`
	Humidity      int     `json:"humidity"`
	Wind          float32 `json:"wind_kph"`
	WindDir       int     `json:"wind_degree"`
	Qnh           float32 `json:"pressure_mb"`
	FeelsLike     float32 `json:"feelslike_c"`
	Precipitation float32 `json:"precip_mm"`
	ChanceOfRain  int     `json:"chance_of_rain"`
}

// Print forecast
func PrintForecast(forecast Forecast, arrow rune, thresholdRain int, language string) string {

	rainForecast := ""
	if forecast.ChanceOfRain >= thresholdRain {
		rainForecast = fmt.Sprintf(
			"%s %s: %d%% %.2fmm %c%s",
			CRed,
			lang.Language(language)["rain"],
			forecast.ChanceOfRain,
			forecast.Precipitation,
			arrow,
			CReset,
		)
	}
	return fmt.Sprintf(
		"  %s %03d/%.0f %s %.0f/%0.f %d%% %0.f [%0.f°]%s\n",
		forecast.Time[11:],
		forecast.WindDir,
		forecast.Wind,
		forecast.Condition.Text,
		forecast.Temp,
		forecast.DewPoint,
		forecast.Humidity,
		forecast.Qnh,
		forecast.FeelsLike,
		rainForecast,
	)

}

// WhichArrow determine the rain trend
// returns the arrow rune and previous precipitation value
func WhichArrow(previousRain float32, forecast Forecast) (rune, float32) {

	// Firts line of forecast or no change from previous
	if previousRain == 9999.0 || forecast.Precipitation == previousRain {
		previousRain = forecast.Precipitation
		return 0, 0
	}

	// UpArrow
	if forecast.Precipitation > previousRain {
		return UpArrow, forecast.Precipitation
	}
	// DownArrow
	return DownArrow, forecast.Precipitation
}
