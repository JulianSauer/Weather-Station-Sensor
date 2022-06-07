package main

import (
    "encoding/json"
    "fmt"
    "github.com/JulianSauer/Weather-Station-Sensor/cache"
    "github.com/JulianSauer/Weather-Station-Sensor/sensor"
    "github.com/JulianSauer/Weather-Station-Sensor/sns"
    "time"
)

func main() {
    location, _ := time.LoadLocation("Europe/Berlin")
    timestamp := time.Now().In(location).Format("02 January 2006 15:04:05")
    fmt.Printf("Starting Weather Station at %s\n", timestamp)

    s := sensor.GetSensor()

    publishSensorData(s)
    checkBattery(s)
}

func publishSensorData(s sensor.Sensor) {
    checkCache()
    messages, e := collectData(s)
    if e != nil {
        fmt.Println("Could not read sensor data")
        fmt.Println(e.Error())
    }
    if messages != nil {
        sns.PublishSensorData(messages)
    }
}

func checkBattery(s sensor.Sensor) {
    batteryLow, e := s.BatteryIsLow()
    if e != nil {
        fmt.Println("Could not ready battery state")
        fmt.Println(e.Error())
    }
    if batteryLow {
        sns.PublishLowBattery()
    }
}

func collectData(s sensor.Sensor) (*[]string, error) {
    weatherData, e := s.GetWeatherData()
    if e != nil {
        return nil, e
    }
    weatherDataAsJson := make([]string, len(weatherData) * 2)
    for i := 0; i + 1 < len(weatherDataAsJson); i += 2 {
        jsonBytes, e := json.Marshal(weatherData[i])
        if e != nil {
            return nil, e
        }
        weatherDataAsJson[i] = string(jsonBytes)

        weatherData[i].Timestamp = "latest"
        jsonBytes, e = json.Marshal(weatherData[i])
        if e != nil {
            return nil, e
        }
        weatherDataAsJson[i + 1] = string(jsonBytes)
    }
    return &weatherDataAsJson, nil
}

func checkCache() {
    unpublishedMessages := cache.ReadAll()
    if unpublishedMessages == nil {
        return
    }
    sns.PublishSensorData(unpublishedMessages)
}
