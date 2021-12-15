package main

import (
	"context"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
	"github.com/cookienyancloud/quotestemplate/configs"
	"github.com/cookienyancloud/quotestemplate/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	welcome = `
бот цитат. Пришлите цитату, автора, возможно фото.
Одним сообщением, разделяя переносом.`
)

type Temp struct {
	QuoteText   string
	Name        string
	PhotoHeight string
	PhotoWidth  string
	PhotoURL    string
}

var tmplPage = template.Must(template.ParseFiles("./static/index.html"))

func main() {
	conf, err := configs.InitConf()
	if err != nil {
		log.Fatalf("error init config:%v\n", err)
	}

	bot, updates, err := tgBot.StartBot(conf.TgToken)
	if err != nil {
		log.Fatalf("error init bot:%v\n", err)
	}
	go func() {
		http.HandleFunc("/image", makeScreen)
		http.ListenAndServe(conf.Addr, nil)
	}()

	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
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

		mes := make([]string, 2)
		if update.Message.Caption != "" {
			mes = strings.Split(update.Message.Caption, "\n")
			leng := len(update.Message.Photo)
			phUrl, err := bot.GetFileDirectURL((update.Message.Photo)[leng-1].FileID)
			if err != nil {
				log.Printf("err getting photo URL: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "err getting photo URL")
				_, _ = bot.Send(msg)
			}
			h := strconv.Itoa((update.Message.Photo)[leng-1].Height)
			w := strconv.Itoa((update.Message.Photo)[leng-1].Width)
			q := fmt.Sprintf("%s%s%s?quote-text=%s&name=%s&photo-height=%s&photo-width=%s&photo-url=%s",
				"localhost", conf.Addr, "/image",
				mes[0], mes[1], h, w, phUrl)
			fmt.Println(q)
			filename := "screen.png"
			var imageBuf []byte
			if err := chromedp.Run(ctx, ScreenshotTasks(q, &imageBuf)); err != nil {
				log.Printf("err running chrome: %v", err)
			}
			if err := ioutil.WriteFile(filename, imageBuf, 0644); err != nil {
				log.Fatal(err)
			}

		} else if update.Message.Text != "" {
			//mes = strings.Split(update.Message.Caption, "\n")

		} else {

		}

	}

}

func makeScreen(w http.ResponseWriter, r *http.Request) {
	data := &Temp{
		QuoteText:   r.URL.Query()["quote-text"][0],
		Name:        r.URL.Query()["name"][0],
		PhotoHeight: r.URL.Query()["photo-height"][0],
		PhotoWidth:  r.URL.Query()["photo-width"][0],
		PhotoURL:    r.URL.Query()["photo-url"][0],
	}
	err := tmplPage.Execute(w, data)
	if err != nil {
		log.Printf("error execuring: %v", err)
	}
}

func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(url),
		chromedp.ActionFunc(func(ctx context.Context) (err error) {
			*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
			return err
		}),
	}
}