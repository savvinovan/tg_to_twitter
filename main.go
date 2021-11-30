package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/spf13/viper"
)

func main() {
	config := initConfig()
	bot, err := tgbotapi.NewBotAPI(config.TgBotToken)
	if err != nil {
		log.Panic(err)
	}
	//bot.Debug = true
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		var response string

		//log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
		log.Println(update.Message.Text)
		switch update.Message.Text {
		case "huy":
			response = "pizda"
		case "/hello":
			response = "doroobo doroobo"
		case "/catfact":
			response = getCatFact(config)
		case "/dogfact":
			response = getDogFact(config)
		default:
			response = update.Message.Text
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, response)
		msg.ReplyToMessageID = update.Message.MessageID

		_, err := bot.Send(msg)
		if err != nil {
			log.Println("err ", err)
		}
		//fmt.Println(m)
	}
}

func getCatFact(c *Config) string {
	resp, err := http.Get(c.CatFactURL)
	if err != nil {
		log.Fatalln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Print(err)
	}
	textRaw, ok := result["text"]
	if !ok {
		textRaw = ""
	}
	text, ok := textRaw.(string)
	if !ok {
		text = ""
	}
	return text
}

func getDogFact(c *Config) string {
	resp, err := http.Get(c.DogFactURL)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(string(body))
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Print(err)
	}
	textRaw, ok := result["facts"]
	if !ok {
		textRaw = ""
	}
	text, ok := textRaw.([]interface{})
	if !ok {
		text = []interface{}{"asdf"}
	}
	res := text[0].(string)
	return res
}

func initConfig() *Config {
	var c Config
	viper.SetConfigName("dev")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(fmt.Errorf("Fatal error config file: %w \n", err))
	}
	err = viper.Unmarshal(&c)
	if err != nil {
		log.Panic(fmt.Errorf("Fatal error unmarshalling config file: %w \n", err))
		return nil
	}
	return &c

}

type Config struct {
	TgBotToken string
	CatFactURL string
	DogFactURL string
}
