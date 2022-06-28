package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

func startBot() {
	botToken := "5514527614:AAFbQ4kuwG3QrbJeO12haFH3inoFIiAQvJk"
	botApi := "https://api.telegram.org/bot"
	botUrl := botApi + botToken
	offset := 0
	for {
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("Updates error: ", err.Error())
		}
		for _, update := range updates {
			err = respond(botUrl, update)
			if err != nil {
				log.Println("Updates error: ", err.Error())
			}
			offset = update.UpdateId + 1
			fmt.Println(update)
		}
	}
}

func getUpdates(botUrl string, offset int) ([]Update, error) {
	resp, err := http.Get(botUrl + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var restResponse RestResponse
	err = json.Unmarshal(body, &restResponse)
	if err != nil {
		return nil, err
	}
	return restResponse.Result, nil
}

func respond(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	switch update.Message.Text {
	case "/start":
		botMessage.Text = "Hi, my name is my_weather_bot. Please, send me your city."
	default:
		weather, err := weatherAPI(trueCityName(update.Message.Text))
		if err != nil {
			return err
		}
		if weather.Location.Name == "" {
			botMessage.Text = "Неизвестный город."
		} else {
			botMessage.Text = fmt.Sprintln(update.Message.Text, weather.Current.Temp,
				"°С\nВетер", weather.Current.WindMPH, "метров в секуду")
		}
	}
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	_, err = http.Post(botUrl+"/sendMessage", "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	return nil
}

func main() {
	startBot()
}
