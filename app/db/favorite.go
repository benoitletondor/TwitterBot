package db

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type Favorite struct {
	id        int
	UserId    int64
	UserName  string
	TweetId   int64
	Status    string
	FavDate   time.Time
	UnfavDate time.Time
}

const (
	_TABLE_FAVORITE = "favorite"
)

func (fav Favorite) Persist() error {
	stmtIns, err := database.Prepare("INSERT INTO " + _TABLE_FAVORITE + "(userId, userName, tweetId, status, favDate, unfavDate) VALUES( ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	unfavDate := mysql.NullTime{Time: fav.UnfavDate, Valid: !fav.UnfavDate.IsZero()}

	_, err = stmtIns.Exec(fav.UserId, fav.UserName, fav.TweetId, fav.Status, fav.FavDate, unfavDate)
	return err
}

func HasAlreadyFav(tweetId int64) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_FAVORITE + " WHERE tweetId = ? LIMIT 1")
	if err != nil {
		return true, err
	}

	defer stmtOut.Close()

	var size int

	err = stmtOut.QueryRow(tweetId).Scan(&size)
	if err != nil {
		return true, err
	}

	return size > 0, nil
}
