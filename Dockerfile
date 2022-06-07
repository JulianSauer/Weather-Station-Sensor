FROM golang:1.18.3 as build

WORKDIR /go/src/github.com/JulianSauer/Weather-Station-Sensor
ADD . /go/src/github.com/JulianSauer/Weather-Station-Sensor

RUN go get -d -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o Weather-Station-Sensor .

FROM alpine:3.8

RUN apk add --no-cache tzdata

COPY --from=build /go/src/github.com/JulianSauer/Weather-Station-Sensor /bin/
RUN chmod 755 /bin/Weather-Station-Sensor

RUN echo "*/15 * * * *  /bin/Weather-Station-Sensor >> /var/log/Weather-Station-Sensor.log" >> /home/cron-schedule.txt
RUN /usr/bin/crontab /home/cron-schedule.txt

CMD ["/usr/sbin/crond", "-f", "-l", "8"]
