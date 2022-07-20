/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
    "sync"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/service/sqs"
)

// awsCmd represents the aws command
var awsCmd = &cobra.Command{
	Use:   "aws",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		switch globalFlags.action {
		case "sqs":
			// SQSにメッセージ送信 -> メッセージ内容取得 -> 削除までの一連の流れ
			sendQueue(args)
		default:
			fmt.Println("please specify action")
			fmt.Println("go run main.go aws -a [service]")
			os.Exit(1)
		}


	},
}

var svcSQS *sqs.SQS

func receiveMessages() {

	params := &sqs.ReceiveMessageInput{
        QueueUrl: aws.String(os.Getenv("AWS_SQS_URL")),
        MaxNumberOfMessages: aws.Int64(10),
        WaitTimeSeconds: aws.Int64(20),
    }
    resp, err := svcSQS.ReceiveMessage(params)

    if err != nil {
        fmt.Println(err)
	}

	if len(resp.Messages) == 0 {
        fmt.Println("empty queue.")
    }

    var wg sync.WaitGroup
    for _, m := range resp.Messages {
        wg.Add(1)
        go func(msg *sqs.Message) {
            defer wg.Done()
			fmt.Println(msg)
			// メッセージ削除する
			if err := DeleteMessage(msg); err != nil {
                fmt.Println(err)
            }
        }(m)
    }

    wg.Wait()
}

func DeleteMessage(msg *sqs.Message) error {
    params := &sqs.DeleteMessageInput{
        QueueUrl:      aws.String(os.Getenv("AWS_SQS_URL")),
        ReceiptHandle: aws.String(*msg.ReceiptHandle),
    }
    _, err := svcSQS.DeleteMessage(params)

    if err != nil {
        return err
    }
    return nil
}


func sendQueue(args []string) {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Printf("環境変数が呼び出せませんでした！: %v", err)
	}

	fmt.Println("aws called")

	const (
		AWS_REGION = "ap-northeast-1"
	)


	// AWS SDKを利用するための設定
	sess := session.Must(session.NewSession())
	svcSQS = sqs.New(
		sess,
		aws.NewConfig().WithRegion("ap-northeast-1"),
	)

	for _, msg := range args {
		params := &sqs.SendMessageInput{
			MessageBody:  aws.String(msg),
			QueueUrl:     aws.String(os.Getenv("AWS_SQS_URL")),
			DelaySeconds: aws.Int64(0),
			MessageGroupId: aws.String("TEST"),
		}
		sqsRes, err := svcSQS.SendMessage(params)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(sqsRes)
	}

	// Receive Messages
	receiveMessages()

}

func init() {
	rootCmd.AddCommand(awsCmd)
}
