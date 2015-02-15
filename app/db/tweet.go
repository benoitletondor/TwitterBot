package db

import (
	"time"
)

type Tweet struct {
	id      int
	Content string
	Date    time.Time
}

const (
	_TABLE_TWEET = "tweet"
)

func (tweet Tweet) Persist() error {
	stmtIns, err := database.Prepare("INSERT INTO " + _TABLE_TWEET + "(content, date) VALUES( ?, ? )")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(tweet.Content, tweet.Date)
	return err
}
