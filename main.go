package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/cookienyancloud/quotestemplate/configs"
	"github.com/cookienyancloud/quotestemplate/tgBot"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	welcome = `
бот цитат. Пришлите фото, цитату, автора цитаты, возможно автора фото.
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

		if update.Message.Command() == "change" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, welcome)
			_, _ = bot.Send(msg)
			continue
		}

		if update.Message.Caption != "" {
			mes := strings.Split(update.Message.Caption, "\n")
			if len(mes) < 2 || len(mes) >3 {
				fmt.Println(len(mes))
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "цитата\nавтор\nавтор фото")
				_, _ = bot.Send(msg)
				continue
			}
			leng := len(update.Message.Photo)
			phUrl, err := bot.GetFileDirectURL((update.Message.Photo)[leng-1].FileID)
			if err != nil {
				log.Printf("err getting photo URL: %v", err)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "err getting photo URL")
				_, _ = bot.Send(msg)
				continue
			}
			h := strconv.Itoa((update.Message.Photo)[leng-1].Height)
			w := strconv.Itoa((update.Message.Photo)[leng-1].Width)
			author := ""
			if len(mes) == 3 {
				author = mes[2]
			}
			q := fmt.Sprintf("%s?quote-text=%s&name=%s&author=%s&photo-height=%s&photo-width=%s&photo-url=%s",
				conf.URL,
				mes[0], mes[1], author, h, w, phUrl)
			fmt.Println(q)
			apiFlashEndpoint := "https://api.apiflash.com/v1/urltoimage"
			request, _ := http.NewRequest("GET", apiFlashEndpoint, nil)
			query := request.URL.Query()
			query.Add("access_key", conf.ApiKey)
			query.Add("element", ".container")
			query.Add("width", h)
			query.Add("height", h)
			query.Add("url", q)
			request.URL.RawQuery = query.Encode()

			client := &http.Client{}
			response, _ := client.Do(request)
			file := tgbotapi.FileReader{
				Name:   "screenshot.jpeg",
				Reader: response.Body,
			}
			photo := tgbotapi.NewPhoto(update.Message.From.ID, file)
			bot.Send(photo)
			response.Body.Close()
		} else if update.Message.Text != "" {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "пока так не буду")
			_, _ = bot.Send(msg)
			continue

		} else {
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "хз")
			_, _ = bot.Send(msg)
			continue
		}

	}

}

//
//func makeScreen(w http.ResponseWriter, r *http.Request) {
//	data := &Temp{
//		QuoteText:   r.URL.Query()["quote-text"][0],
//		Name:        r.URL.Query()["name"][0],
//		PhotoHeight: r.URL.Query()["photo-height"][0],
//		PhotoWidth:  r.URL.Query()["photo-width"][0],
//		PhotoURL:    r.URL.Query()["photo-url"][0],
//	}
//	err := tmplPage.Execute(w, data)
//	if err != nil {
//		log.Printf("error execuring: %v", err)
//	}
//}

//
//func ScreenshotTasks(url string, imageBuf *[]byte) chromedp.Tasks {
//	return chromedp.Tasks{
//		chromedp.Navigate(url),
//		chromedp.ActionFunc(func(ctx context.Context) (err error) {
//			*imageBuf, err = page.CaptureScreenshot().WithQuality(90).Do(ctx)
//			return err
//		}),
//	}
//}
