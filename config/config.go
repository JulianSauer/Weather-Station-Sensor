package config

import (
    "encoding/json"
    "log"
    "os"
)

const CONFIG_NAME = "config.json"

type Config struct {
    SensorType                       string `json:"sensorType"`
    TinkerforgeWeatherStationAddress string `json:"tinkerforgeWeatherStationAddress"`
    TinkerforgeWeatherStationPort    string `json:"tinkerforgeWeatherStationPort"`
    TinkerforgeWeatherStationUID     string `json:"tinkerforgeWeatherStationUID"`
    MobileAlertsUrl                  string `json:"mobileAlertsUrl"`
    MobileAlertsDeviceId             string `json:"mobileAlertsDeviceId"`
    AWSSNSWeatherTopic               string `json:"AWS-SNS-Weather-Topic"`
    AWSSNSBatteryTopic               string `json:"AWS-SNS-Battery-Topic"`
}

func Load() *Config {
    file, e := os.Open(CONFIG_NAME)
    if e != nil {
        log.Fatal(e.Error())
    }

    defer file.Close()
    decoder := json.NewDecoder(file)
    config := Config{}
    if e = decoder.Decode(&config); e != nil {
        log.Fatal("could not parse config.json")
    }
    return &config
}
