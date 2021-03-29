package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-Pi/config"
    "github.com/JulianSauer/Weather-Station-Pi/sns"
    "github.com/JulianSauer/Weather-Station-Pi/weather"
    "github.com/robfig/cron/v3"
    "os"
    "os/signal"
)

func main() {
    cronConfig := config.Load().Cron
    c := cron.New()
    _, e := c.AddFunc(cronConfig, func() {
        data := collectData()
        for _, d := range data {
            sns.Publish(d)
        }
    })
    if e != nil {
        fmt.Println(e.Error())
    } else {
        c.Start()
        fmt.Printf("Starting cron job: %s\n", cronConfig)
        fmt.Println("Hit ctrl+c to stop")
        sig := make(chan os.Signal)
        signal.Notify(sig, os.Interrupt, os.Kill)
        <-sig
    }
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
