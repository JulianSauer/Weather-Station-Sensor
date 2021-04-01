package config

import (
    "encoding/json"
    "log"
    "os"
)

const CONFIG_NAME = "config.json"

type Config struct {
    WeatherStationAddress string `json:"weatherStationAddress"`
    WeatherStationPort    string `json:"weatherStationPort"`
    WeatherStationUID     string `json:"weatherStationUID"`
    AWSSNSWeatherTopic    string `json:"AWS-SNS-Weather-Topic"`
    AWSSNSBatteryTopic    string `json:"AWS-SNS-Battery-Topic"`
    Cron                  string `json:"cron"`
    CronBattery           string `json:"cronBattery"`
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
