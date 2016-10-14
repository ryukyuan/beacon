package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/kazuph/go-binenv"
	"fmt"
)

func main() {
	port := "80"

	router := gin.New()

	router.GET("/", func(c *gin.Context) {
		c.String(200, "OK")
	})

	router.POST("/callback", func(c *gin.Context) {
		bot := getLineBot()
		events, err := bot.ParseRequest(c.Request)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				c.Writer.WriteHeader(400)
			} else {
				c.Writer.WriteHeader(500)
			}
			return
		}
		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(message.Text + "だってばよ")).Do(); err != nil {
						log.Print(err)
					}
				}
			}
			if event.Type == linebot.EventTypeBeacon {
				if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("ビーコンが検知したってばよ")).Do(); err != nil {
					log.Print(err)
				}
			}
		}
	})

	router.Run(":" + port)
}

func getLineBot() (bot *linebot.Client) {

	env, err := binenv.Load(Asset)
	if err != nil {
		fmt.Println(err)
	}

	bot, botErr := linebot.New(
		env["CHANNEL_SECRET"],
		env["CHANNEL_TOKEN"],
	)
	if botErr != nil {
		log.Fatal(botErr)
	}
	return
}