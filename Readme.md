# Weather Station Sensor
Reads data from a [Tinkerforge Weather Station](https://www.tinkerforge.com/en/shop/outdoor-weather-station-ws-6147.html) and sends it to a message hub (AWS Simple Notification Service).
The service sends the weather data only once on execution and is meant to be run as a cron job.