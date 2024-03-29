package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/esperlu/weatherCLI/lang"
	"github.com/esperlu/weatherCLI/utils"
)

const (
	configFile = "/home/jeanluc/.config/weatherCLI/config.json"
	ver        = "1.1"
)

// Config struct to assign values from config.json
type Config struct {
	APIKey                  string `json:"APIkey"`
	DefaultCity             string `json:"defaultCity"`
	DefaultDaysForecast     int    `json:"defaultDaysForecast"`
	DefaultThresholdRain    int    `json:"defaultThresholdRain"`
	DefaultTomorrowForecast bool   `json:"defaultTomorrowForecast"`
	DefaultLanguage         string `json:"defaultLanguage"`
}

func main() {

	startTime := time.Now()

	// Read config file
	var config Config
	configFile, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal(err)
	}
	// Parse config file
	err = json.Unmarshal(configFile, &config)

	// Flags
	daysForecast := flag.Int("d", config.DefaultDaysForecast, "Number of day forecast")
	thresholdRain := flag.Int("t", config.DefaultThresholdRain, "Threshold % of rain -t 55")
	tomerrowForecast := flag.Bool("f", config.DefaultTomorrowForecast, "Tomorrow's forecast")
	language := flag.String("l", config.DefaultLanguage, "Language (fr default)")
	version := flag.Bool("v", false, "Version")

	// Custom help message
	flag.Usage = func() {
		fmt.Println(utils.Help)
	}
	flag.Parse()

	if *version {
		fmt.Printf(
			"%s %s Compiled with %s \n",
			utils.Filename(os.Args[0]),
			ver,
			runtime.Version(),
		)
		return
	}

	// Populate cities []string
	cities := flag.Args()

	// Default city
	nbArgs := flag.NArg()
	if nbArgs == 0 {
		cities = []string{config.DefaultCity}
		nbArgs = 1
	}

	// Loop through cities in go routines
	var wg sync.WaitGroup
	var i int
	ch := make(chan string, nbArgs)

	for _, city := range cities {
		i++
		wg.Add(1)
		go func(city string, i int) {

			defer wg.Done()
			days := fmt.Sprintf("%d", *daysForecast)
			url := "https://api.weatherapi.com/v1/forecast.json?key=" + config.APIKey + "&q=" + city + "&days=" + days + "&aqi=no&lang=" + *language

			resp, err := http.Get(url)
			if err != nil {
				log.Fatal(err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatal(err)
			}

			// Unmarshal JSON response.Body
			var wx utils.Weather
			err = json.Unmarshal(body, &wx)
			if err != nil {
				log.Fatal(err)
			}

			// Print API error message
			if wx.Error.Message != "" {
				fmt.Printf(
					"\n%s%s%s\n",
					utils.CRed,
					wx.Error.Message,
					utils.CReset,
				)
				return
			}

			// Print current wx
			sunRise := utils.FormatPMtime(wx.Forecast.ForecastDay[0].Astro.SunRise)
			sunSet := utils.FormatPMtime(wx.Forecast.ForecastDay[0].Astro.SunSet)
			outputString := fmt.Sprintf(
				"%d-%s%s - %s (%.2f,%.2f) %s-%s %s",
				i,
				utils.CGreen,
				wx.Location.Name,
				wx.Location.Country,
				wx.Location.Lat, wx.Location.Long,
				sunRise, sunSet,
				utils.CReset,
			)

			// If rain in current wx
			var currentPrecipitation string
			if wx.Current.Precipitation > 0 {
				currentPrecipitation = fmt.Sprintf(
					"%s%s:%.2fmm%s",
					utils.CRed,
					lang.Language(*language)["precipitation"],
					wx.Current.Precipitation,
					utils.CReset,
				)
			}

			outputString += fmt.Sprintf(
				"\n%s %03d/%.0f %.0f %s"+
					" %.0f/%.0f %d%% Q%.0f [%.0f°] %s\n"+
					"%s:\n",
				wx.Current.LastUpdate[11:],
				wx.Current.WindDegree,
				wx.Current.Wind,
				wx.Current.Visibility*1000,
				wx.Current.Condition.Text,
				wx.Current.Temp,
				wx.Current.Temp-((100.0-float32(wx.Current.Humidity))/5.0),
				wx.Current.Humidity,
				wx.Current.Qnh,
				wx.Current.FeelsLike,
				currentPrecipitation,
				lang.Language(*language)["Forecast for the day"],
			)

			//  Print forecast of day 0 (today)
			var previousRain float32 = 9999.0
			var arrow rune
			for _, forecast := range wx.Forecast.ForecastDay[0].Hour {

				if int(time.Now().Unix()) > forecast.TimeEpoch {
					continue
				}

				// Test rain trend and assign arrow value and previousRain
				arrow, previousRain = utils.RainTrendArrow(previousRain, forecast)
				// Print to output
				outputString += utils.PrintForecast(forecast, arrow, *thresholdRain, *language)

			}

			// Print forecast days after today
			for _, forecast := range wx.Forecast.ForecastDay {
				// skip current day
				if int(time.Now().Unix()) >= forecast.DateEpoch {
					continue
				}

				date, err := time.Parse("2006-01-02", forecast.Date)
				if err != nil {
					log.Fatal(err)
				}
				weekDay := fmt.Sprintf("%s", date.Weekday())
				outputString += fmt.Sprintf("%s - %-9s",
					forecast.Date,
					lang.Language(*language)[weekDay],
				)

				// If rain in forecast
				rainForecast := ""
				red := ""
				if forecast.Day.ChanceOfRain >= *thresholdRain {
					red = utils.CRed
				}
				if forecast.Day.ChanceOfRain > 0 {
					rainForecast = fmt.Sprintf(
						"%s %s: %d%% %.2fmm%s",
						red,
						lang.Language(*language)["precipitation"],
						forecast.Day.ChanceOfRain,
						forecast.Day.Precipitation,
						utils.CReset,
					)
				}

				// Print forecast
				outputString += fmt.Sprintf(
					" %s %.0f->%.0f° %.0f kmh %s %.0f%% %s %s\n",
					forecast.Day.Condition.Text,
					forecast.Day.MinTemp,
					forecast.Day.MaxTemp,
					forecast.Day.MaxWind,
					lang.Language(*language)["humidity"],
					forecast.Day.AvgHumidity,
					rainForecast,
					utils.CReset,
				)

				// print forecast for tomorrow only
				previousRain = 9999.0
				arrow = 0
				if *tomerrowForecast {
					for _, tomorrow := range forecast.Hour {
						tomorrowTime := time.Now().AddDate(0, 0, 1)
						tomorrowString := tomorrowTime.Format("2006-01-02")
						// skip if not tomorrow between 09:00 and 22:00
						if tomorrowString != tomorrow.Time[:10] ||
							tomorrow.Time[11:] < "09:00" ||
							tomorrow.Time[11:] > "22:00" {
							continue
						}

						// Test rain trend and assign arrow value and previousRain
						arrow, previousRain = utils.RainTrendArrow(previousRain, tomorrow)
						// Print to output
						outputString += utils.PrintForecast(tomorrow, arrow, *thresholdRain, *language)
					}
				}

			}

			ch <- outputString
		}(city, i)
	}

	wg.Wait()
	close(ch)

	// Read channel
	var values []string
	for val := range ch {
		values = append(values, val)
	}

	// Sort the output to preserve the arguments order
	sort.Strings(values)

	// Print final output (omitting the int used to sort the args)
	for _, value := range values {
		dashPos := strings.Index(value, "-") + 1
		fmt.Println(value[dashPos:])
	}

	// print timing
	fmt.Printf(
		"%s | %0.3f sec.\n%s\n",
		utils.PrintProgName(ver),
		time.Since(startTime).Seconds(),
		"Source: WeatherAPI https://www.weatherapi.com",
	)

}
