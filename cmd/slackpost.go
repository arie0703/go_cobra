/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"github.com/joho/godotenv"
	"github.com/slack-go/slack"
	"github.com/spf13/cobra"
	"time"
	"strconv"
)

// slackpostCmd represents the slackpost command
var slackpostCmd = &cobra.Command{
	Use:   "slackpost",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		//環境変数を呼び出す
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Printf("環境変数が呼び出せませんでした！: %v", err)
		} 
		
		
		switch globalFlags.action {
		case "help":
			help()
		case "message":
			message(args)
		case "new":
			new()
		case "record":
			getHistory()
		default:
			fmt.Printf("メッセージを入力してね！")
		}
		
	},
}

func help() {
	fmt.Println("slackpost -a message ~~~~~:  新規メッセージを投稿")
	fmt.Println("slackpost -a new:  今日の日付のスレッドを投稿")
	fmt.Println("slackpost -a record:  本日の投稿記録を取得")
}

func message(args []string) {
	tkn := os.Getenv("SLACK_API_TOKEN")
	channelName := os.Getenv("SLACK_CHANNEL_NAME")
	c := slack.New(tkn)

	if len(args) == 0 {
		fmt.Println("メッセージを入力してね")
	} else {
		_, _, err := c.PostMessage(channelName, slack.MsgOptionText(args[0], true))
		if err != nil {
			panic(err)
		}

		fmt.Println("投稿完了！: " + args[0])
	}
	
}

func new() {
	
	tkn := os.Getenv("SLACK_API_TOKEN")
	channelName := os.Getenv("SLACK_CHANNEL_NAME")
	c := slack.New(tkn)

	t := time.Now()
	const layout = "1月2日" 
	msg := t.Format(layout) + "の活動記録"
	//今日の日付のスレッドを作成
	_, _, err := c.PostMessage(channelName, slack.MsgOptionText(msg, true))
	if err != nil {
		panic(err)
	}
	fmt.Println("投稿完了！: " + msg)
}

func getHistory() {
	tkn := os.Getenv("SLACK_API_TOKEN")
	channelId:= os.Getenv("SLACK_CHANNEL_ID")
	c := slack.New(tkn)

	now := time.Now()
	oldest := strconv.FormatInt(time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).UnixNano(), 10)
	latest := strconv.FormatInt(time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, time.Local).UnixNano(), 10)
	param := slack.GetConversationHistoryParameters{
		ChannelID:	channelId,
		Oldest: oldest[:10] + "." + oldest[11:16],
		Latest: latest[:10] + "." + latest[11:16],
	}

	history, err := c.GetConversationHistory(&param)
	if err != nil {
		fmt.Printf("%s\n", err)
		return
	}
	fmt.Println("〜本日の投稿〜")
	for i, m := range history.Messages {
		fmt.Println(i+1, m.Text)
	}
}

func init() {
	rootCmd.AddCommand(slackpostCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// slackpostCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// slackpostCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// var action string
	// rootCmd.PersistentFlags().StringVarP(&action, "action", "a", "help", "action command")
}
