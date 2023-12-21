# CLI weather app in Go

This CLI application is written in the Go language and is designed to connect to a weather API provided by [weatherAPI.com](https://www.weatherapi.com/)

The CLI app enables developers to access weather data while working in a console.


## Installation

1. Run the following command to install the weatherCLI repo in the directory defined in your `$GOPATH` environment variable:  
`go get github.com/esperlu/weatherCLI` 
2. Navigate to the now local sources: `$GOPATH/src/github.com/esperlu/weatherCLI`
3. Give it a try:
`go run weather.go Brussels`
4. If successfull, compile the weather sources and data:
    * To compile the binary and save it in the current directory, run the following command:  
    `go build weather.go`
    * To compile the binary and install it in the binary folder defined in the `$GOBIN` environment variable, run the following command:  
    `go install weather.go`  
    This will make the binary accessible and executable system wide.

## Setup and Config file ##

You have first to create an account on [WeatherAPI](https://www.weatherapi.com/) and insert your API key in a `config.json` file. You can also set there the default values for other variables.

    {
        "APIkey": "your_api_key_here",
        "defaultCity": "Brussels",
        "defaultDaysForecast": 5,
        "defaultThresholdRain": 60,
        "defaultTomorrowForecast": false,
        "defaultLanguage": "fr"
    }

The path to your `config.json` file is to be set in the `Config` section of `weather.go`

    // Config
    const (
        configFile = "/path/to/config.json"
        ver        = "1.1"
    )


## Usage

    Usage: $ weather BRU paris new-york

    Default station : Brussels

    option      Description                      Defaut
    ---------------------------------------------------
    -l  string   Language code e.g. de, fr or es     fr     
    -d  int      Number of day forecast               5
    -t  int      Warning threshold for rain (%)      50
    -f  bool     Print tomorrow's forecast        false
    -v  bool     Print version                    false

The default values for the flags can be changed in your [config.json file](#Setup-and-config-file)


## Examples

To get the Berlin's current weather and forecast for the next 3 days in German language with the probability threshold for rain set to 60%:

    $ weather -l de -t 60 -d 3 Berlin

Will return:

    Berlin - Germany (52.52,13.40) 08:15-15:54 
    18:00 230/31 10000 leichter Regenfall 7/6 93% Q994 [3°]
    Die heutige Wettervorhersage:
        19:00 280/36 Klar 7/1 66% 983 [2°]
        20:00 277/38 leichte Regenschauer 6/0 66% 984 [1°] Regen: 100% 0mm
        21:00 271/38 bewölkt 6/-1 64% 984 [0°]
        22:00 269/44 bedeckt 5/-1 65% 984 [-1°]
        23:00 272/45 stellenweise Regenfall 4/0 75% 985 [-2°] Regen: 100% 0mm
    2023-12-22 - Freitag   stellenweise Regenfall 2->4° 45 kmh Luftfeuchtigkeit 65%  Regen: 89% 1mm 
    2023-12-23 - Samstag   starker Regenfall 1->6° 31 kmh Luftfeuchtigkeit 89%  Regen: 79% 25mm 


- First line format:

    `City name - country (lat, long) sunrise-sunset`
- Second line (close to aviation METAR format):

    `Time | Wind direction/Wind speed (kmh) | Visibility | WX Description | Temp/DewPoint | Relative Humidity | QNH | [Wind Chill factor]`

- Forecast lines:

    `Time | Wind direction/Wind speed (kmh) | WX Description | Temp/DewPoint | Relative Humidity | QNH | [Wind Chill factor] | Rain forecast`

## Add your language

For the time beeing, the languages available are:
- English (en)
- French (fr)
- German (de)
- Spanish (es)
- Portuguese (pt)
- Deutch (nl)

I you want to add a new language or change an existing one, just edit the [`lang/lang.go`](https://github.com/esperlu/weatherCLI/blob/main/lang/lang.go) file to your liking.

## Wheather data source

- API provided by [Weather API](https://www.weatherapi.com)
- WeatherAPI Docs [Weather API](https://www.weatherapi.com/docs/)
