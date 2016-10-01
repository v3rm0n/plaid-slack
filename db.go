package main

import (
	"log"
	"database/sql"
	"time"
	"github.com/lib/pq"
)

func GetDB(url string) *sql.DB {
	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func GetLastEventTime(db *sql.DB, result chan time.Time) {
	log.Print("Getting last event time from DB")
	rows, err := db.Query("SELECT max(date) FROM event")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		var date string
		if err := rows.Scan(&date); err != nil {
			log.Fatal(err)
		}
		log.Printf("Last event time: %s\n", date)
		timestamp, err := time.Parse("2006-01-02T15:04:05Z", date)
		if err != nil {
			log.Fatal(err)
		}
		result <- timestamp
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
}

func SaveEvent(event TimelineEvent, db *sql.DB) {
	log.Print("Saving timeline event", event.Title, "to DB with date", event.DateOpened)
	result, err := db.Exec("INSERT INTO event(date,description) VALUES($1,$2)", pq.FormatTimestamp(event.DateOpened.Time), event.Title)
	if err != nil {
		log.Fatal(err)
	}
	i, err := result.RowsAffected()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Rows inserted: %s", i)
}