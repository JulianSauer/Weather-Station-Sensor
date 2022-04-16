package sns

import (
    "fmt"
    "github.com/JulianSauer/Weather-Station-Sensor/cache"
    "github.com/JulianSauer/Weather-Station-Sensor/config"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
)

var topic string

func init() {
    topic = config.Load().AWSSNSWeatherTopic
}

func PublishLowBattery() {
    session, e := session.NewSession(&aws.Config{
        Region: aws.String("eu-central-1"),
    })

    if e != nil {
        fmt.Println(e.Error())
        return
    }

    client := sns.New(session)
    input := &sns.PublishInput{
        Message:  aws.String("Battery is low"),
        TopicArn: aws.String(topic),
    }

    result, e := client.Publish(input)
    if e != nil {
        fmt.Printf("Cannot publish: %s\n", e.Error())

    } else {
        fmt.Printf("%s: Battery is low\n", *result.MessageId)
    }
}

func PublishSensorData(messages *[]string) {
    session, e := session.NewSession(&aws.Config{
        Region: aws.String("eu-central-1"),
    })

    if e != nil {
        fmt.Println(e.Error())
        return
    }

    client := sns.New(session)
    var unpublishableMessages []string
    for _, message := range *messages {
        input := &sns.PublishInput{
            Message:  aws.String(message),
            TopicArn: aws.String(topic),
        }

        result, e := client.Publish(input)
        if e != nil {
            fmt.Printf("Cannot publish: %s\n", e.Error())
            unpublishableMessages = append(unpublishableMessages, message)
        } else {
            fmt.Printf("%s: %s\n", *result.MessageId, message)
        }
    }
    if len(unpublishableMessages) > 0 {
        cache.WriteAll(&unpublishableMessages)
    }
}
