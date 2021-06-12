package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-Pi/cache"
    "github.com/JulianSauer/Weather-Station-Pi/config"
    "github.com/JulianSauer/Weather-Station-Pi/sns"
    "github.com/JulianSauer/Weather-Station-Pi/weather"
    "github.com/robfig/cron/v3"
    "os"
    "os/signal"
)

func main() {
    cronConfig := config.Load()
    c := cron.New()

    _, e := c.AddFunc(cronConfig.Cron, func() {
        checkCache()
        messages := collectData()
        if messages != nil {
            sns.PublishSensorData(messages)
        }
    })
    if e != nil {
        fmt.Println(e.Error())
        return
    }

    _, e = c.AddFunc(cronConfig.CronBattery, func() {
        if weather.BatteryIsLow() {
            sns.PublishLowBattery()
        }
    })
    if e != nil {
        fmt.Println(e.Error())
        return
    }

    c.Start()
    fmt.Printf("Starting cron job: %s\n", cronConfig.Cron)
    fmt.Printf("Starting cron job: %s\n", cronConfig.CronBattery)
    fmt.Println("Hit ctrl+c to stop")
    sig := make(chan os.Signal)
    signal.Notify(sig, os.Interrupt, os.Kill)
    <-sig
}

func collectData() *[]string {
    weatherData := weather.GetWeatherData()
    if weatherData == nil {
        return nil
    }
    weatherDataAsJson := make([]string, len(weatherData) * 2)
    for i := 0; i + 1 < len(weatherDataAsJson); i += 2 {
        jsonBytes, e := json.Marshal(weatherData[i])
        if e != nil {
            fmt.Println(e.Error())
        }
        weatherDataAsJson[i] = string(jsonBytes)

        weatherData[i].Timestamp = "latest"
        jsonBytes, e = json.Marshal(weatherData[i])
        if e != nil {
            fmt.Println(e.Error())
        }
        weatherDataAsJson[i + 1] = string(jsonBytes)
    }
    return &weatherDataAsJson
}

func checkCache() {
    unpublishedMessages := cache.ReadAll()
    if unpublishedMessages == nil {
        return
    }
    sns.PublishSensorData(unpublishedMessages)
}
