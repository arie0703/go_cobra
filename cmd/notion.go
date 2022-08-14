/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os"
	"io"
	"github.com/joho/godotenv"
	"net/http"
	"github.com/spf13/cobra"
	"log"
	"encoding/json"
)

// type Result struct {
// 	Results []string
// 	Paragraph string
// 	Rich_text string
// 	Type string
// 	Block string
// }

// notionCmd represents the notion command
var notionCmd = &cobra.Command{
	Use:   "notion",
	Short: "Notion APIを使って色々呼び出す",
	Long: `GoでNotion APIを呼び出す実験・JSON処理の学習`,
	Run: func(cmd *cobra.Command, args []string) {
		//環境変数を呼び出す
		err := godotenv.Load(".env")
		if err != nil {
			fmt.Printf("環境変数が呼び出せませんでした！: %v", err)
		}

		tkn := os.Getenv("NOTION_SECRET_TOKEN")
		target := os.Getenv("JUL_1")
		base_url := "https://api.notion.com/v1/blocks/" + target
		client := &http.Client{}
		req, err := http.NewRequest("GET", base_url + "/children", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Add("Authorization", tkn)
		req.Header.Add("Notion-Version", "2022-06-28")
		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		// APIレスポンスをbyte型として読み込む
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		var r interface{}
		json.Unmarshal(b, &r);

		//　JSONレスポンスをマッピング（めんどくさい）
		results := r.(map[string]interface{})["results"].([]interface{})
		length := len(results)

		// 対象のページの内容を出力する
		for i := 0; i < length; i++ {
			paragraph := results[i].(map[string]interface{})["paragraph"]
			if paragraph != nil && len(paragraph.(map[string]interface{})["rich_text"].([]interface{})) > 0 {
				content := paragraph.(map[string]interface{})["rich_text"].([]interface{})[0].(map[string]interface{})["text"].(map[string]interface{})["content"]
				fmt.Println(content)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(notionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// notionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// notionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
