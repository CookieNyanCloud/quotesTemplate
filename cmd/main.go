package main

import (
	"github.com/cookienyancloud/quotestemplate/configs"
	"github.com/cookienyancloud/quotestemplate/photo"
	"github.com/cookienyancloud/quotestemplate/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"log"
	"strings"
)

const (
	welcome = `
бот цитат. Пришлите цитату, автора, возможно фото.
Одним сообщением, разделяя переносом.`
)

func main() {
	conf, err := configs.InitConf()
	if err != nil {
		log.Fatalf("error init config:%v\n", err)
	}

	bot, updates, err := tgBot.StartBot(conf.TgToken)
	if err != nil {
		log.Fatalf("error init bot:%v\n", err)
	}
	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcome)
			_, _ = bot.Send(msg)
			continue
		}

		var query []string

		if update.Message.Text != "" {
			query = strings.Split(update.Message.Text, "\n")
			if len(query) != 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "чего-то не хватает")
				_, _ = bot.Send(msg)
				continue
			}

		} else if update.Message.Caption != "" {
			query = strings.Split(update.Message.Caption, "\n")
			if len(query) != 2 {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "чего-то не хватает")
				_, _ = bot.Send(msg)
				continue
			}
			var s string
			if update.Message.Document != nil {
				s, err = bot.GetFileDirectURL(update.Message.Document.FileID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					_, _ = bot.Send(msg)
					continue
				}
			} else if (*update.Message.Photo)[len((*update.Message.Photo))-1].FileID != "" {
				s, err = bot.GetFileDirectURL((*update.Message.Photo)[len((*update.Message.Photo))-1].FileID)
				if err != nil {
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
					_, _ = bot.Send(msg)
					continue
				}
			}

			file, err := photo.DownloadFile(s, query[1])
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
				continue
			}
			defer file.Close()

			send := tgbotapi.FileReader{
				Name:   file.Name(),
				Reader: file,
			}

			upload := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, send)
			res, err := bot.Send(upload)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, res.Text)
			_, _ = bot.Send(msg)

			//decode, _, err := image.Decode(file)
			//if err != nil {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			//	_, _ = bot.Send(msg)
			//	continue
			//}

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "что-то не так")
			_, _ = bot.Send(msg)
			continue
		}

	}
}
