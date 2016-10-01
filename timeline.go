package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"log"
)

type DateTime struct {
	time.Time
}

func (dateTime *DateTime) UnmarshalJSON(b []byte) error {
	timeStr := string(b[:])
	if "null" == timeStr {
		return nil
	}
	t, err := time.Parse("\"2006-01-02 15:04\"", timeStr)
	if err != nil {
		return err
	}
	dateTime.Time = t
	return nil
}

type TimelineEvent struct {
	Description, Title, Level string
	DateOpened                *DateTime `json:"date_opened"`
	DateClosed                *DateTime `json:"date_closed"`
	Id                        int
}

func (t TimelineEvent) String() string {
	return fmt.Sprintf("%s, %s", t.Title, t.DateClosed)
}

func GetTimeline(result chan []TimelineEvent) {
	log.Print("Getting Plaid issues timeline")
	var timeline []TimelineEvent
	resp, err := http.Get("https://status.plaid.com/issues/timeline/")
	if err != nil {
		log.Fatal("error:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("error:", err)
	}
	err = json.Unmarshal(body, &timeline)
	if err != nil {
		log.Fatal("error:", err)
	}
	result <- timeline
}
