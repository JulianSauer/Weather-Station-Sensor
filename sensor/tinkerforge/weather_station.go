package tinkerforge

import (
    "errors"
    "fmt"
    "github.com/JulianSauer/Weather-Station-Sensor/config"
    "github.com/JulianSauer/Weather-Station-Sensor/dto"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
    "github.com/Tinkerforge/go-api-bindings/outdoor_weather_bricklet"
    "os"
    "time"
)

var address string
var uid string
var location *time.Location
var identifiers []uint8

type TfSensor struct {}

func (tfSensor TfSensor) Setup() {
    configuration := config.Load()
    address = configuration.TinkerforgeWeatherStationAddress + ":" + configuration.TinkerforgeWeatherStationPort
    uid = configuration.TinkerforgeWeatherStationUID
    location, _ = time.LoadLocation("Europe/Berlin")
    fmt.Println("Connecting to weather sensor")
    outdoorWeatherBricklet, connection := reconnect()
    defer disconnect(connection)
    if outdoorWeatherBricklet == nil || connection == nil {
        return
    }
    var e error
    identifiers, e = outdoorWeatherBricklet.GetStationIdentifiers()
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
}

func (tfSensor TfSensor) BatteryIsLow() (bool, error) {
    fmt.Println("Checking battery")
    outdoorWeatherBricklet, connection := reconnect()
    defer disconnect(connection)
    if outdoorWeatherBricklet == nil || connection == nil {
        return false, errors.New("could not connect to weather station")
    }
    for _, identifier := range identifiers {
        _, _, _, _, _, _, batteryLow, _, e := outdoorWeatherBricklet.GetStationData(identifier)
        if e != nil {
            return false, e
        }

        if batteryLow {
            return true, nil
        }
    }
    return false, nil
}

func (tfSensor TfSensor) GetWeatherData() ([]dto.WeatherData, error) {
    fmt.Println("Checking sensor data")
    result := make([]dto.WeatherData, len(identifiers))
    outdoorWeatherBricklet, connection := reconnect()
    defer disconnect(connection)
    if outdoorWeatherBricklet == nil || connection == nil {
        return nil, errors.New("could not connect to weather station")
    }
    for i, identifier := range identifiers {
        temperature, humidity, windSpeed, gustSpeed, rain, windDirection, _, _, e := outdoorWeatherBricklet.GetStationData(identifier)
        if e != nil {
            return nil, e
        }

        result[i] = convertWeatherData(time.Now().In(location), temperature, humidity, windSpeed, gustSpeed, rain, windDirection)
    }
    return result, nil
}

func reconnect() (*outdoor_weather_bricklet.OutdoorWeatherBricklet, *ipconnection.IPConnection) {
    connection := ipconnection.New()
    outdoorWeatherBricklet, e := outdoor_weather_bricklet.New(uid, &connection)
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }

    if e = connection.Connect(address); e != nil {
        fmt.Printf("Could not connect to weather sensor: %s", e.Error())
        return nil, nil
    }
    return &outdoorWeatherBricklet, &connection
}

func disconnect(connection *ipconnection.IPConnection) {
    connection.Disconnect()
    connection.Close()
}

func convertWeatherData(timestamp time.Time, temperature int16, humidity uint8, windSpeed uint32, gustSpeed uint32, rain uint32, windDirection outdoor_weather_bricklet.WindDirection) dto.WeatherData {
    t := timestamp.Format("20060102-150405")
    return dto.WeatherData{
        Source:        "WeatherStation",
        Timestamp:     t,
        DataFor:       []string{t},
        Temperature:   []string{fmt.Sprintf("%.1f", float64(temperature)/10.0)},
        Humidity:      []string{fmt.Sprintf("%d", humidity)},
        WindSpeed:     []string{fmt.Sprintf("%.1f", float64(windSpeed)/10.0)},
        GustSpeed:     []string{fmt.Sprintf("%.1f", float64(gustSpeed)/10.0)},
        Rain:          []string{fmt.Sprintf("%.1f", float64(rain)/10.0)},
        WindDirection: []string{fmt.Sprintf("%.1f", float32(windDirection)*22.5)}, // Convert to degrees
    }
}
