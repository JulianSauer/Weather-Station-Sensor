package sns

import (
    "fmt"
    "github.com/JulianSauer/Weather-Station-Pi/config"
    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/sns"
)

var topic string

func init() {
    topic = config.Load().AWSSNSTopic
}

func Publish(message string) {
    session, e := session.NewSession(&aws.Config{
        Region: aws.String("eu-central-1"),
    })

    if e != nil {
        fmt.Println(e.Error())
        return
    }

    client := sns.New(session)
    input := &sns.PublishInput{
        Message:  aws.String(message),
        TopicArn: aws.String(topic),
    }

    result, e := client.Publish(input)
    if e != nil {
        fmt.Println(e.Error())
    } else {
        fmt.Printf("%s: %s\n", *result.MessageId, message)
    }
}
