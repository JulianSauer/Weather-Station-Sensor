# Weather Station Sensor
Reads data from a [Tinkerforge Weather Station](https://www.tinkerforge.com/en/shop/outdoor-weather-station-ws-6147.html) or a [Mobile Alerts temperature sensor](https://mobile-alerts.eu/temperatursensor-ma10100/) and sends it to a message hub (AWS Simple Notification Service).
The service sends the weather data only once on execution and is meant to be run as a cron job.

Example config.json:
```json
{
  "sensorType":  "",
  "tinkerforgeWeatherStationAddress": "",
  "tinkerforgeWeatherStationPort":  "",
  "tinkerforgeWeatherStationUID": "",
  "mobileAlertsUrl":  "",
  "mobileAlertsDeviceId": "",
  "AWS-SNS-Weather-Topic":  "",
  "AWS-SNS-Battery-Topic": ""
}
```

The sensor type has to be `Tinkerforge` or `MobileAlerts`.
Depending on that either fill in the values for the Tinkerforge weather station or for the Mobile Alerts sensor.
The topics are ARNs.

This also requires a .aws/credentials file with an access key and a secret access key that can push messages to the SNS topics.