package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
	"github.com/cookienyancloud/quotestemplate/configs"
	"github.com/cookienyancloud/quotestemplate/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
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

	ctx, cancel := chromedp.NewContext(
		context.Background(),
		chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	for update := range updates {

		if update.Message == nil {
			continue
		}

		if update.Message.Command() == "start" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcome)
			_, _ = bot.Send(msg)
			continue
		}
		//mes := make([]string, 2)
		if update.Message.Caption != "" {
			//mes = strings.Split(update.Message.Caption, "\n")
			var buf []byte
			if err := chromedp.Run(ctx, elementScreenshot(`https://pkg.go.dev/`, `img.Homepage-logo`, &buf)); err != nil {
				log.Printf("error during screenshot: %v", err)
			}
			if err := ioutil.WriteFile("elementScreenshot.png", buf, 0o644); err != nil {
				log.Printf("error during writing: %v", err)
			}

		} else if update.Message.Text != "" {
			//mes = strings.Split(update.Message.Caption, "\n")

		} else {

		}

	}

}

func elementScreenshot(urlstr, sel string, res *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Screenshot(sel, res, chromedp.NodeVisible),
	}
}
