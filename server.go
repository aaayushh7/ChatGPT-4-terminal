package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func GetResponse2(client gpt3.Client, ctx context.Context, question string) string {
	var response string
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt: []string{
			question,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		response = resp.Choices[0].Text
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	return response
}

type NullWriter2 int

func (NullWriter) Write2([]byte) (int, error) { return 0, nil }

func main2() {
	log.SetOutput(new(NullWriter))
	viper.SetConfigFile(".env")
	viper.ReadInConfig()
	apiKey := viper.GetString("API_KEY")
	if apiKey == "" {
		panic("Missing API key")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	r := gin.Default()
	r.POST("/send-message", func(c *gin.Context) {
		var json struct {
			Message string `json:"message"`
		}
		if c.BindJSON(&json) == nil {
			response := GetResponse2(client, ctx, json.Message)
			c.JSON(200, gin.H{
				"response": response,
			})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
