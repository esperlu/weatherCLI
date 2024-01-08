package lang

// Languge switch
func Language(code string) map[string]string {

	lang := make(map[string]map[string]string)

	lang["en"] = map[string]string{
		"Sunday":               "Sunday",
		"Monday":               "Monday",
		"Tuesday":              "Tuesday",
		"Wednesday":            "Wednesday",
		"Thursday":             "Thursday",
		"Friday":               "Friday",
		"Saturday":             "Saturday",
		"precipitation":        "precipitation",
		"Forecast for the day": "Forecast for the day",
		"humidity":             "humidity",
	}

	lang["fr"] = map[string]string{
		"Sunday":               "Dimanche",
		"Monday":               "Lundi",
		"Tuesday":              "Mardi",
		"Wednesday":            "Mercredi",
		"Thursday":             "Jeudi",
		"Friday":               "Vendredi",
		"Saturday":             "Samedi",
		"precipitation":        "precipitation",
		"Forecast for the day": "Prévisions du jour",
		"humidity":             "humidité",
	}

	lang["de"] = map[string]string{
		"Sunday":               "Sonntag",
		"Monday":               "Montag",
		"Tuesday":              "Dienstag",
		"Wednesday":            "Mittwoch",
		"Thursday":             "Donnerstag",
		"Friday":               "Freitag",
		"Saturday":             "Samstag",
		"Forecast for the day": "Die heutige Wettervorhersage",
		"precipitation":        "Niederschlag",
		"humidity":             "Luftfeuchtigkeit",
	}

	lang["es"] = map[string]string{
		"Sunday":               "Domingo",
		"Monday":               "Lunes",
		"Tuesday":              "Martes",
		"Wednesday":            "Miercoles",
		"Thursday":             "Jueves",
		"Friday":               "Viernes",
		"Saturday":             "Sabado",
		"Forecast for the day": "Pronóstico del día",
		"precipitation":        "precipitación",
		"humidity":             "humedad",
	}

	lang["pt"] = map[string]string{
		"Sunday":               "Domingo",
		"Monday":               "Segunda",
		"Tuesday":              "Terca",
		"Wednesday":            "Quarta",
		"Thursday":             "Quinta",
		"Friday":               "Sexta",
		"Saturday":             "Sabado",
		"Forecast for the day": "Previsão para hoje",
		"precipitation":        "precipitação",
		"humidity":             "umidade",
	}

	lang["nl"] = map[string]string{
		"Sunday":               "Zondag",
		"Monday":               "Maandag",
		"Tuesday":              "Dinsdag",
		"Wednesday":            "Woensdag",
		"Thursday":             "Donderdag",
		"Friday":               "Vrijdag",
		"Saturday":             "Zaterdag",
		"Forecast for the day": "Verwachtingen voor volgende uren",
		"precipitation":        "neerslag",
		"humidity":             "vochtigheid",
	}

	return lang[code]
}
