package main

import (
	"os"
	"log"
	"time"
)

var dbUrl string

func init() {
	dbUrl = os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		log.Fatal("$DATABASE_URL must be set")
	}
}

func main() {
	log.Println("Plaid Slackbot started")
	timelineResult := make(chan []TimelineEvent)
	go GetTimeline(timelineResult)

	db := GetDB(dbUrl)
	defer db.Close()
	lastEventTimeResult := make(chan time.Time)
	go GetLastEventTime(db, lastEventTimeResult)

	timeline := <-timelineResult
	lastEventTime := <-lastEventTimeResult

	for _, event := range timeline {
		if (event.DateOpened.Time.After(lastEventTime)) {
			SendSlackMessage(event)
			SaveEvent(event, db)
		}
	}
	log.Println("Plaid Slackbot finished")
}


