package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-Pi/weather"
)

func main() {

}

func collectData() []string {
    weatherData := weather.GetWeatherData()
    weatherDataAsJson := make([]string, len(weatherData))
    for i, data := range weatherData {
        jsonBytes, e := json.Marshal(data)
        if e != nil {
            fmt.Println(e.Error())
        }
        weatherDataAsJson[i] = string(jsonBytes)
    }
    return weatherDataAsJson
}
