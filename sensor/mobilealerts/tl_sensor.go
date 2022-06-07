package mobilealerts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/JulianSauer/Weather-Station-Sensor/config"
	"github.com/JulianSauer/Weather-Station-Sensor/dto"
	"gopkg.in/resty.v1"
	"time"
)

type MaSensor struct{}

type MADevice struct {
	Deviceid    string `json:"deviceid"`
	Lastseen    int    `json:"lastseen"`
	Lowbattery  bool   `json:"lowbattery"`
	Measurement struct {
		Idx int     `json:"idx"`
		Ts  int     `json:"ts"`
		C   int     `json:"c"`
		Lb  bool    `json:"lb"`
		T1  float64 `json:"t1"`
	} `json:"measurement"`
}

type MADeviceRequest struct {
	Devices []MADevice `json:"devices"`
	Success bool       `json:"success"`
}

func (maSensor MaSensor) Setup() {}

func (maSensor MaSensor) GetWeatherData() ([]dto.WeatherData, error) {
	fmt.Println("Checking sensor data")

	deviceResult, e := querySensor()
	if e != nil {
		return nil, e
	}

	temperature := deviceResult.Measurement.T1
	return convertWeatherData(temperature), nil
}

func (maSensor MaSensor) BatteryIsLow() (bool, error) {
	deviceResult, e := querySensor()
	if e != nil {
		return false, e
	}
	return deviceResult.Lowbattery, nil
}

func querySensor() (*MADevice, error) {
	configFile := config.Load()
	url := configFile.MobileAlertsUrl
	parameters := "deviceids=" + configFile.MobileAlertsDeviceId

	client := resty.New()
	response, e := client.R().
		SetHeader("Content-Type", "application/x-www-form-urlencoded").
		SetBody(parameters).
		Post(url)

	if e != nil {
		return nil, e
	}

	var deviceRequest MADeviceRequest
	if e = json.Unmarshal(response.Body(), &deviceRequest); e != nil {
		return nil, e
	}
	if len(deviceRequest.Devices) != 1 {
		return nil, errors.New(fmt.Sprintf("Found %d devices, expected 1", len(deviceRequest.Devices)))
	}

	return &deviceRequest.Devices[0], nil
}

func convertWeatherData(temperature float64) []dto.WeatherData {
	location, _ := time.LoadLocation("Europe/Berlin")
	timestamp := time.Now().In(location)
	t := timestamp.Format("20060102-150405")
	return []dto.WeatherData{
		{
			Source:        "TemperatureSensor",
			Timestamp:     t,
			DataFor:       []string{t},
			Temperature:   []string{fmt.Sprintf("%.1f", temperature)},
			Humidity:      []string{""},
			WindSpeed:     []string{""},
			GustSpeed:     []string{""},
			Rain:          []string{""},
			WindDirection: []string{""},
		},
	}
}
