package cache

import (
    "bufio"
    "fmt"
    "os"
)

const CACHE_NAME = "WeatherStationCache.json"

func WriteAll(messages *[]string) {
    fmt.Printf("Writing %d messages to cache\n", len(*messages))
    c, e := os.OpenFile(CACHE_NAME, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if e != nil {
        fmt.Println(e.Error())
        return
    }
    defer c.Close()
    for _, message := range *messages {
        if _, e = c.WriteString(message + "\n"); e != nil {
            fmt.Println(e.Error())
        }
    }
}

func ReadAll() *[]string {
    if _, e := os.Stat(CACHE_NAME); os.IsNotExist(e) {
        return nil
    }
    fmt.Println("Reading from cache")
    c, e := os.Open(CACHE_NAME)
    if e != nil {
        fmt.Println(e.Error())
        return nil
    }
    defer c.Close()

    var messages []string
    scanner := bufio.NewScanner(c)
    for scanner.Scan() {
        messages = append(messages, scanner.Text())
    }
    fmt.Printf("Found %d messages in cache\n", len(messages))
    if e = os.Remove(CACHE_NAME); e != nil {
        fmt.Printf("Cannot delete %s: %s\n", CACHE_NAME, e.Error())
    }
    return &messages
}
