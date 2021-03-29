package tests

import (
    "bufio"
    "github.com/JulianSauer/Weather-Station-Pi/cache"
    "io/ioutil"
    "os"
    "strings"
    "testing"
)

var messagesToWrite = []string{
    `"timestamp":"20210101-000000","temperature":22.1,"humidity":46,"windSpeed":0,"gustSpeed":0,"rain":28.8,"windDirection":0}`,
    `"timestamp":"20210101-000001","temperature":22.1,"humidity":46,"windSpeed":0,"gustSpeed":0,"rain":28.8,"windDirection":0}`,
    `"timestamp":"20210101-000002","temperature":22.1,"humidity":46,"windSpeed":0,"gustSpeed":0,"rain":28.8,"windDirection":0}`,
}

func TestWriteALl(t *testing.T) {
    cleanup(t)

    cache.WriteAll(&messagesToWrite)

    // Check file content
    c, e := os.Open(cache.CACHE_NAME)
    if e != nil {
        t.Error(e)
    }
    defer c.Close()

    var foundMessages []string
    scanner := bufio.NewScanner(c)
    for scanner.Scan() {
        foundMessages = append(foundMessages, scanner.Text())
    }

    if len(foundMessages) != 3 {
        t.Errorf("Expected %d but found %d messages", 3, len(foundMessages))
    }
    for i, foundMessage := range foundMessages {
        if foundMessage != messagesToWrite[i] {
            t.Errorf("Expected: %s", messagesToWrite[i])
            t.Errorf("Actual:   %s", foundMessage)
        }
    }

    cleanup(t)
}

func TestReadAll(t *testing.T) {
    cleanup(t)

    foundMessages := cache.ReadAll()
    if foundMessages != nil {
        t.Errorf("Expected %d messages, found %d", 0, len(*foundMessages))
    }

    messageBytes := []byte(strings.Join(messagesToWrite[:], "\n"))
    if e := ioutil.WriteFile(cache.CACHE_NAME, messageBytes, 0644); e != nil {
        t.Error(e)
    }

    foundMessages = cache.ReadAll()
    if len(*foundMessages) != 3 {
        t.Errorf("Expected %d messages, found %d", 3, len(*foundMessages))
    }
    for i, foundMessage := range *foundMessages {
        if foundMessage != messagesToWrite[i] {
            t.Errorf("Expected: %s", messagesToWrite[i])
            t.Errorf("Actual:   %s", foundMessage)
        }
    }
    if _, e := os.Stat(cache.CACHE_NAME); !os.IsNotExist(e) {
        t.Errorf("File should have been deleted but exists: %s", cache.CACHE_NAME)
    }

    cleanup(t)
}

func TestReadWrite(t *testing.T) {
    cleanup(t)

    messagesToWritePart1 := messagesToWrite[0:2]
    cache.WriteAll(&messagesToWritePart1)
    foundMessages := cache.ReadAll()

    if len(*foundMessages) != 2 {
        t.Errorf("Expected %d but found %d messages", 2, len(*foundMessages))
    }
    if _, e := os.Stat(cache.CACHE_NAME); !os.IsNotExist(e) {
        t.Errorf("File should have been deleted but exists: %s", cache.CACHE_NAME)
    }

    cache.WriteAll(foundMessages)
    messagesToWritePart2 := []string{messagesToWrite[2]}
    cache.WriteAll(&messagesToWritePart2)
    if _, e := os.Stat(cache.CACHE_NAME); os.IsNotExist(e) {
        t.Errorf("File should exist: %s", cache.CACHE_NAME)
    }

    foundMessages = cache.ReadAll()

    if len(*foundMessages) != 3 {
        t.Errorf("Expected %d but found %d messages", 3, len(*foundMessages))
    }
    if _, e := os.Stat(cache.CACHE_NAME); !os.IsNotExist(e) {
        t.Errorf("File should have been deleted but exists: %s", cache.CACHE_NAME)
    }

    cleanup(t)
}

func cleanup(t *testing.T) {
    if _, e := os.Stat(cache.CACHE_NAME); !os.IsNotExist(e) {
        if e = os.Remove(cache.CACHE_NAME); e != nil {
            t.Error(e)
        }
    }
}
