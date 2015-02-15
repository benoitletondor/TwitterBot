package db

import (
	"time"
)

type Follow struct {
	id           int
	UserId       int64
	UserName     string
	Status       string
	FollowDate   time.Time
	UnfollowDate time.Time
}

const (
	_TABLE_FOLLOW = "follow"
)

func AlreadyFollow(userId int64) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_FOLLOW + " WHERE userId = ? LIMIT 1")
	if err != nil {
		return true, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(userId).Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}

func (follow Follow) Persist() error {
	stmtIns, err := database.Prepare("INSERT INTO " + _TABLE_FOLLOW + "(userId, userName, status, followDate) VALUES( ?, ?, ?,? )")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(follow.UserId, follow.UserName, follow.Status, follow.FollowDate)
	return err
}
