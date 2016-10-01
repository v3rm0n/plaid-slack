package main

import (
	"net/http"
	"encoding/json"
	"bytes"
	"log"
	"os"
)

var hookUrl string
var channel string

func init() {
	hookUrl = os.Getenv("SLACK_HOOK_URL")
	channel = os.Getenv("SLACK_CHANNEL")
	switch {
	case hookUrl == "":
		log.Fatal("$SLACK_HOOK_URL must be set")
	case channel == "":
		log.Fatal("$SLACK_CHANNEL must be set")
	}
}

type SlackMessage struct {
	Text        string `json:"text"`
	Username    string `json:"username"`
	Channel     string `json:"channel"`
	Attachments []SlackMessageAttachment `json:"attachments"`
}

type SlackMessageAttachment struct {
	Fallback string `json:"fallback"`
	Color    string `json:"color"`
	Title    string `json:"title"`
	Text     string `json:"text"`
	Ts       int64 `json:"ts"`
}

func getColor(level string) string {
	if "error" == level {
		return "#d9514e";
	} else if "warning" == level {
		return "#fd9357";
	} else {
		return "#35ce8d";
	}
}

func SendSlackMessage(event TimelineEvent) {
	log.Println("Sending message to slack: ", event.Title)
	attachments := []SlackMessageAttachment{
		{
			Text: event.Description,
			Fallback:event.Description,
			Title: event.Title,
			Color: getColor(event.Level),
			Ts: event.DateOpened.Unix(),
		},
	}
	slackMessage := SlackMessage{
		Channel: channel,
		Username: "Plaid Slackbot",
		Attachments: attachments,
	}
	msg, err := json.Marshal(slackMessage)
	if err != nil {
		log.Println("error:", err)
		return
	}
	resp, err := http.Post(hookUrl, "application/json", bytes.NewBuffer(msg))
	if err != nil {
		log.Println("error:", err)
		return
	}
	defer resp.Body.Close()
	log.Println("Response from Slack:", resp.Status)
}