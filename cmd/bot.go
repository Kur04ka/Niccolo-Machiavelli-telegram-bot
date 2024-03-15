package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"

	"github.com/Kur04ka/telegram_bot/internal/config"
	"github.com/Kur04ka/telegram_bot/internal/quote/db"
	"github.com/Kur04ka/telegram_bot/pkg/client/open_ai"
	"github.com/Kur04ka/telegram_bot/pkg/client/postgresql"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sashabaranov/go-openai"
)

var (
	pathToScripts string = "C:/Users/cobak/OneDrive/Рабочий стол/programming/tg_bot/internal/src/script_answers/"
	pathToImages  string = "C:/Users/cobak/OneDrive/Рабочий стол/programming/tg_bot/internal/src/images/"
)

func main() {
	cfg := config.GetConfig()
	cfgDB := cfg.Postgresql

	log.Println("Creating new openAI client...")
	openai_client, err := open_ai.NewOpenAIClient(cfg.Proxy.User, cfg.Proxy.Password, cfg.Proxy.ProxyAddress, cfg.Proxy.Port, cfg.Tokens.OpenAIToken)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("OpenAI client was successfully created")

	log.Println("Creating new bot instance...")
	bot, err := tgbotapi.NewBotAPI(cfg.Tokens.TelegramToken)
	if err != nil {
		log.Panic(err)
	}
	log.Println("New bot instance was successfully created")

	log.Println("Connecting to Postgresql...")
	connection := postgresql.NewPostgresqlClient(cfgDB.User, cfgDB.Password, cfgDB.Host, cfgDB.Port, cfgDB.DbName)
	db := db.NewStorage(connection)
	defer connection.Close(context.Background())
	log.Println("Successfully connected to Postgresql")

	bot.Debug = true
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := bot.GetUpdatesChan(updateConfig)

	var isChat bool = false
	for update := range updates {
		if update.Message == nil {
			continue
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Неизвестная команда, чтобы выести список команд введи /help")
		if !update.Message.IsCommand() {
			if isChat {
				resp, err := openai_client.CreateChatCompletion(
					context.Background(),
					openai.ChatCompletionRequest{
						Model:     openai.GPT3Dot5Turbo,
						MaxTokens: 1000,
						Messages: []openai.ChatCompletionMessage{
							{
								Role: openai.ChatMessageRoleUser,
								Content: fmt.Sprintf(`Веди весь диалог со мной от имени Никколо Макиавелли, 
								учитывая все известные работы данного философа, его образ мышления, на сколько это возможно,
								сделай так, чтобы я почувствовал, что я общаюсь с Никколо Макиавелли. Вот мой вопрос:
								"%s"`, update.Message.Text),
							},
						},
					},
				)
				if err != nil {
					log.Fatalf("ChatCompletion error: %v\n", err)
				}

				msg.Text = resp.Choices[0].Message.Content
				msg.ReplyToMessageID = update.Message.MessageID

			}

			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("error sending message, error: %v\n", err)
			}

		} else {
			switch update.Message.Command() {
			case "start":
				msg.Text = readFromFile(fmt.Sprint(pathToScripts, "start.txt"))
			case "help":
				msg.Text = readFromFile(fmt.Sprint(pathToScripts, "help.txt"))
			case "biography":
				photo := tgbotapi.NewPhoto(update.Message.Chat.ID, tgbotapi.FilePath(fmt.Sprintf("%sportrait.png", pathToImages)))
				if _, err := bot.Send(photo); err != nil {
					log.Fatalf("error sending message, error: %v\n", err)
				}
				msg.Text = readFromFile(fmt.Sprint(pathToScripts, "biography.txt"))
			case "quote":
				quote, err := db.FindOne(context.TODO(), strconv.Itoa(rand.Intn(20)+1))
				if err != nil {
					log.Fatalln(err)
				}
				msg.Text = fmt.Sprintf("%s\n@Никколо Макиавелли", quote.Quote)
			case "chat":
				isChat = true
				msg.Text = readFromFile(fmt.Sprint(pathToScripts, "chat_start.txt"))
			case "stop":
				isChat = false
				msg.Text = readFromFile(fmt.Sprint(pathToScripts, "chat_stop.txt"))
			}

			if _, err := bot.Send(msg); err != nil {
				log.Fatalf("error sending message, error: %v\n", err)
			}
		}
	}
}

func readFromFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading file, error: %v\n", err)
	}
	return string(data)
}
