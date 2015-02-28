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

func HasTweetWithContent(content string) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_TWEET + " WHERE content LIKE ? LIMIT 1")
	if err != nil {
		return true, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(content+"%").Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}
