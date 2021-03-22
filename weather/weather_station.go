package weather

import (
    "fmt"
    "github.com/JulianSauer/Weather-Station-Pi/dto"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
    "github.com/Tinkerforge/go-api-bindings/outdoor_weather_bricklet"
    "os"
    "time"
)

const ADDRESS string = ""
const UID string = ""

var location *time.Location
var identifiers []uint8

func init() {
    location, _ = time.LoadLocation("Europe/Berlin")
    outdoorWeatherBricklet, connection := reconnect()
    defer disconnect(connection)
    var e error
    identifiers, e = outdoorWeatherBricklet.GetStationIdentifiers()
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
}

func GetWeatherData() []dto.WeatherData {
    result := make([]dto.WeatherData, len(identifiers))
    outdoorWeatherBricklet, connection := reconnect()
    defer disconnect(connection)
    for i, identifier := range identifiers {
        temperature, humidity, windSpeed, gustSpeed, rain, windDirection, batteryLow, _, e := outdoorWeatherBricklet.GetStationData(identifier)
        if e != nil {
            fmt.Println(e.Error())
        }

        result[i] = convertWeatherData(time.Now().In(location), temperature, humidity, windSpeed, gustSpeed, rain, windDirection)

        if batteryLow {
            fmt.Println("Battery low!")
        }
    }
    return result
}

func reconnect() (*outdoor_weather_bricklet.OutdoorWeatherBricklet, *ipconnection.IPConnection) {
    connection := ipconnection.New()
    outdoorWeatherBricklet, e := outdoor_weather_bricklet.New(UID, &connection)
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }

    if e = connection.Connect(ADDRESS); e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
    return &outdoorWeatherBricklet, &connection
}

func disconnect(connection *ipconnection.IPConnection) {
    connection.Disconnect()
    connection.Close()
}

func convertWeatherData(timestamp time.Time, temperature int16, humidity uint8, windSpeed uint32, gustSpeed uint32, rain uint32, windDirection outdoor_weather_bricklet.WindDirection) dto.WeatherData {
    return dto.WeatherData{
        Timestamp:     timestamp,
        Temperature:   float64(temperature) / 10.0,
        Humidity:      humidity,
        WindSpeed:     float64(windSpeed) / 10.0,
        GustSpeed:     float64(gustSpeed) / 10.0,
        Rain:          float64(rain) / 10.0,
        WindDirection: float32(windDirection) * 22.5, // Convert to degrees
    }
}
