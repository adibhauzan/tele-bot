package main

import (
	"context"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/sashabaranov/go-openai"
)
var (
	openaiClient *openai.Client
)

func extractOrderID(message string) string {
	re := regexp.MustCompile(`ID → ([\w\[\]]+)`)
	match := re.FindStringSubmatch(message)
	if len(match) > 1 {
		return match[1]
	}
	return ""
}

func main() {
	_ = godotenv.Load()
	botToken := os.Getenv("BOT_TOKEN")
	groupID := os.Getenv("GROUP_ID")
	OpenAIKey := os.Getenv("OPENAI_KEY")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		log.Fatalf("Failed to create bot: %v", err)
	}

	log.Printf("Authorized on account %s", bot.Self.UserName)
	openaiClient = openai.NewClient(OpenAIKey)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates , err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalf("Failed to get updates channel: %v", err)
	}

	GroupID, err := strconv.ParseInt(groupID, 10, 64)
	if err != nil {
		log.Fatalf("Failed to parse group ID: %v", err)
	}
	for update := range updates {
		if update.Message == nil || update.Message.Chat == nil {
			continue
		}

		if update.Message.Chat.ID != GroupID {
			continue
		}

		text := update.Message.Text
		log.Printf("Received message: %s", text)

		if strings.Contains(text, "NEW ORDER") {
			orderID := extractOrderID(text)
			if orderID != "" {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "/p "+orderID)
				_, err := bot.Send(msg)
				if err != nil {
					log.Printf("Failed to send message: %v", err)
				} else {
					log.Println("Order ID sent successfully. Exiting bot.")
				}

				time.Sleep(2 * time.Second) 
				return                      
			}
		}
		if strings.HasPrefix(text, "!akumaunanya") {
			prompt := strings.TrimPrefix(text, "!akumaunanya")
			prompt = strings.TrimSpace(prompt)
		
			if prompt == "" {
				prompt = "reset, and say hello"
			}
		
			resp, err := openaiClient.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
				Model: openai.GPT3Dot5Turbo,
				Messages: []openai.ChatCompletionMessage{
					{
						Role:    "user",
						Content: prompt,
					},
				},
			})
			if err != nil {
				errMsg := err.Error()
				if strings.Contains(errMsg, "You exceeded your current quota") {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, "KREDIT ANDA SUDAH HABIS, KREDITNYA DI TOP UP BEGOOOOOOO")
					bot.Send(msg)
				} else {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, errMsg)
					bot.Send(msg)
				}
				continue
			}
		
			answer := resp.Choices[0].Message.Content
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, answer)
			bot.Send(msg)
		}
	}
}