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
    AWSSNSTopic           string `json:"AWS-SNS-Topic"`
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