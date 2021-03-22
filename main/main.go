package main

import (
    "fmt"
    "github.com/Tinkerforge/go-api-bindings/ipconnection"
    "github.com/Tinkerforge/go-api-bindings/outdoor_weather_bricklet"
    "os"
)

const ADDRESS string = ""
const UID string = ""

func main() {

    connection := ipconnection.New()
    defer connection.Close()
    outdoorWeatherBricklet, e := outdoor_weather_bricklet.New(UID, &connection)
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }

    if e = connection.Connect(ADDRESS); e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
    defer connection.Disconnect()

    identifiers, e := outdoorWeatherBricklet.GetStationIdentifiers()
    if e != nil {
        fmt.Println(e.Error())
        os.Exit(1)
    }
    for _, identifier := range identifiers {
        temperature, humidity, windSpeed, gustSpeed, rain, windDirection, batteryLow, _, e := outdoorWeatherBricklet.GetStationData(identifier)
        if e != nil {
            fmt.Println(e.Error())
        }
        fmt.Printf("Temperature (Station): %f °C\n", float64(temperature)/10.0)
        fmt.Printf("Humidity (Station): %d %%RH\n", humidity)
        fmt.Printf("Wind Speed (Station): %f m/s\n", float64(windSpeed)/10.0)
        fmt.Printf("Gust Speed (Station): %f m/s\n", float64(gustSpeed)/10.0)
        fmt.Printf("Rain (Station): %f mm\n", float64(rain)/10.0)
        if windDirection == outdoor_weather_bricklet.WindDirectionN {
            fmt.Println("Wind Direction (Station): N")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionNNE {
            fmt.Println("Wind Direction (Station): NNE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionNE {
            fmt.Println("Wind Direction (Station): NE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionENE {
            fmt.Println("Wind Direction (Station): ENE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionE {
            fmt.Println("Wind Direction (Station): E")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionESE {
            fmt.Println("Wind Direction (Station): ESE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionSE {
            fmt.Println("Wind Direction (Station): SE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionSSE {
            fmt.Println("Wind Direction (Station): SSE")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionS {
            fmt.Println("Wind Direction (Station): S")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionSSW {
            fmt.Println("Wind Direction (Station): SSW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionSW {
            fmt.Println("Wind Direction (Station): SW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionWSW {
            fmt.Println("Wind Direction (Station): WSW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionW {
            fmt.Println("Wind Direction (Station): W")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionWNW {
            fmt.Println("Wind Direction (Station): WNW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionNW {
            fmt.Println("Wind Direction (Station): NW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionNNW {
            fmt.Println("Wind Direction (Station): NNW")
        } else if windDirection == outdoor_weather_bricklet.WindDirectionError {
            fmt.Println("Wind Direction (Station): Error")
        }

        if batteryLow {
            fmt.Println("Battery low!")
        }
    }
}
