package sensor

import (
    "github.com/JulianSauer/Weather-Station-Sensor/config"
    "github.com/JulianSauer/Weather-Station-Sensor/dto"
    "github.com/JulianSauer/Weather-Station-Sensor/sensor/mobilealerts"
    "github.com/JulianSauer/Weather-Station-Sensor/sensor/tinkerforge"
    "log"
)

type Sensor interface {
    Setup()
    GetWeatherData() ([]dto.WeatherData, error)
    BatteryIsLow() (bool, error)
}

func GetSensor() Sensor {
    configuration := config.Load()
    var sensor Sensor
    if configuration.SensorType == "Tinkerforge" {
        sensor = tinkerforge.TfSensor{}
    } else if configuration.SensorType == "MobileAlerts" {
        sensor = mobilealerts.MaSensor{}
    } else {
        log.Fatal("Invalid sensor type in config: " + configuration.SensorType)
    }
    sensor.Setup()
    return sensor
}
