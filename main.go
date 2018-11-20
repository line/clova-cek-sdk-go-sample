package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/line/clova-cek-sdk-go/cek"
	"github.com/line/line-bot-sdk-go/linebot"
)

func sessionEndSpeech(speech string) *cek.ResponseMessage {
	return cek.NewResponseBuilder().
		OutputSpeech(
			cek.NewOutputSpeechBuilder().
				AddSpeechText(speech, cek.SpeechInfoLangJA).Build()).
		ShouldEndSession(true).
		Build()
}

func sessionContinueSpeech(speech string, session *cek.Session) *cek.ResponseMessage {
	return cek.NewResponseBuilder().
		OutputSpeech(
			cek.NewOutputSpeechBuilder().
				AddSpeechText(speech, cek.SpeechInfoLangJA).Build()).
		SessionAttributes(session.SessionAttributes).
		ShouldEndSession(false).
		Build()
}

func handleIntentRequest(req *cek.IntentRequest, session *cek.Session) *cek.ResponseMessage {
	switch req.Intent.Name {
	case "Clova.GuidIntent":
		return sessionContinueSpeech("飲み物を注文して下さい。", session)
	case "OrderBeverage":
		return orderBeverage(req, session)
	case "Clova.YesIntent":
		return orderConfirm(req, session)
	}
	return nil
}

func orderBeverage(req *cek.IntentRequest, session *cek.Session) *cek.ResponseMessage {
	beverage, ok := req.Intent.Slots["beverage"]
	if !ok {
		return sessionEndSpeech("最初からやり直して下さい。")
	}
	amount, ok := req.Intent.Slots["amount"]
	if !ok {
		return sessionEndSpeech("最初からやり直して下さい。")
	}
	log.Printf("OrderCoffee: beverage=%s, amount=%s", beverage.Value, amount.Value)

	session.SessionAttributes["beverage"] = beverage.Value
	session.SessionAttributes["amount"] = amount.Value

	speech := fmt.Sprintf("%sを%s杯でよろしいですか?", beverage.Value, amount.Value)
	return sessionContinueSpeech(speech, session)
}

func orderConfirm(req *cek.IntentRequest, session *cek.Session) *cek.ResponseMessage {
	beverage, ok := session.SessionAttributes["beverage"]
	if !ok {
		return sessionEndSpeech("最初からやり直して下さい。")
	}
	amount, ok := session.SessionAttributes["amount"]
	if !ok {
		return sessionEndSpeech("最初からやり直して下さい。")
	}
	log.Printf("OrderConfirm: beverage=%s, amount=%s", beverage, amount)
	go func() {
		time.Sleep(10 * time.Second)
		sendMessage(session.User.UserID, fmt.Sprintf("%sを%s杯ご用意いたしました。", beverage, amount))
	}()
	return sessionEndSpeech("ご注文ありがとうございました。しばらくお待ち下さい。")
}

func sendMessage(id, text string) {
	client := &http.Client{}
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET"), os.Getenv("CHANNEL_ACCESS_TOKEN"), linebot.WithHTTPClient(client))
	if err != nil {
		log.Printf("LINE bot client initialization error.")
		return
	}
	msg := linebot.NewTextMessage(text)
	_, err = bot.PushMessage(id, msg).Do()
	if err != nil {
		log.Printf("message pushing failed. err=%q", err)
		return
	}
	log.Printf("message pushing succeeded.")
}

func main() {
	ext := cek.NewExtension(os.Getenv("EXTENSION_ID"))
	log.Printf("ExtensionID=%s", ext.ID)
	if os.Getenv("DEBUG_MODE") == "true" {
		cek.WithDebugMode(ext)
	}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		reqMsg, err := ext.ParseRequest(r)
		if err != nil {
			log.Printf("invalid request. err=%+v", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		var response *cek.ResponseMessage
		switch request := reqMsg.Request.(type) {
		case *cek.IntentRequest:
			response = handleIntentRequest(request, reqMsg.Session)
		case *cek.LaunchRequest:
			response = sessionContinueSpeech("いらっしゃいませ。ご注文をどうぞ。", reqMsg.Session)
		case *cek.SessionEndedRequest:
			response = sessionEndSpeech("ありがとうございました。")
		}
		if response != nil {
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		}
	})
	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}
