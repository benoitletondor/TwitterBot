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

	err = stmtOut.QueryRow(content + "%").Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}

func GetNumberOfTweetsBetweenDates(from time.Time, to time.Time) (int, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_TWEET + " WHERE date >= ? AND date <= ? LIMIT 1")
	if err != nil {
		return 0, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(from, to).Scan(&size)
	if err != nil {
		return 0, err
	}

	return size, nil
}
