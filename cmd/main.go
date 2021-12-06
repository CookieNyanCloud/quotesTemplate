package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/cookienyancloud/quotestemplate/chrome"
	"github.com/cookienyancloud/quotestemplate/configs"
	"github.com/cookienyancloud/quotestemplate/photo"
	"github.com/cookienyancloud/quotestemplate/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
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

	abc, cancel:= chromedp.NewContext(context.Background(),chromedp.WithDebugf(log.Printf))
	defer cancel()

	url:="http://127.0.0.1:5500/"
	filename:="test.png"
	var imageBuf []byte
	err = chromedp.Run(abc, chrome.Screen(url, &imageBuf))
	if err != nil {
		log.Fatalf("error ran chrome:%v\n", err)
	}
	err = ioutil.WriteFile(filename, imageBuf, 9544)
	if err != nil {
		log.Fatalf("error write:%v\n", err)
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

			fileName, err := photo.DownloadFile(s, query[1])
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
				continue
			}
			all, err := ioutil.ReadFile(fileName)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
				continue
			}
			buffer := bytes.NewBuffer(all)
			img, err := jpeg.Decode(buffer)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
				continue
			}
			x := img.Bounds().Size().X
			y := img.Bounds().Size().Y
			fmt.Println(x, y)
			send := tgbotapi.FileBytes{
				Name:  query[1] + "jpg",
				Bytes: all,
			}

			upload := tgbotapi.NewPhotoUpload(update.Message.Chat.ID, send)
			_, err = bot.Send(upload)
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)

			}

			//decode, _, err := image.Decode(file)
			//if err != nil {
			//	msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
			//	_, _ = bot.Send(msg)
			//	continue
			//}
			err = os.Remove(query[1] + ".jpeg")
			if err != nil {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, err.Error())
				_, _ = bot.Send(msg)
			}
		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "что-то не так")
			_, _ = bot.Send(msg)
			continue
		}

	}
}
