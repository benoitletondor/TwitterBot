package db

import (
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

type Favorite struct {
	id         int
	UserId     int64
	UserName   string
	TweetId    int64
	Status     string
	FavDate    time.Time
	UnfavDate  time.Time
	LastAction time.Time
}

const (
	_TABLE_FAVORITE = "favorite"
)

func (fav Favorite) Persist() error {
	var stmtIns *sql.Stmt
	var err error

	if fav.id == 0 {
		stmtIns, err = database.Prepare("INSERT INTO " + _TABLE_FAVORITE + "(userId, userName, tweetId, status, favDate, unfavDate, lastAction) VALUES( ?, ?, ?, ?, ?, ?, ?)")
	} else {
		stmtIns, err = database.Prepare("UPDATE " + _TABLE_FAVORITE + " SET userId = ?, userName = ?, tweetId = ?, status = ?, favDate = ?, unfavDate = ?, lastAction = ? WHERE id = ?")
	}

	if err != nil {
		return err
	}

	defer stmtIns.Close()

	unfavDate := mysql.NullTime{Time: fav.UnfavDate, Valid: !fav.UnfavDate.IsZero()}

	if fav.id == 0 {
		_, err = stmtIns.Exec(fav.UserId, fav.UserName, fav.TweetId, fav.Status, fav.FavDate, unfavDate, time.Now())
	} else {
		_, err = stmtIns.Exec(fav.UserId, fav.UserName, fav.TweetId, fav.Status, fav.FavDate, unfavDate, fav.LastAction, fav.id)
	}

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

func GetNotUnfavorite(maxFavDate time.Time, limit int) ([]Favorite, error) {
	favs := make([]Favorite, 0)

	stmtOut, err := database.Prepare("SELECT * FROM " + _TABLE_FAVORITE + " WHERE unfavDate IS NULL AND favDate <= ? ORDER BY lastAction LIMIT ?")
	if err != nil {
		return favs, err
	}

	defer stmtOut.Close()

	rows, err := stmtOut.Query(maxFavDate, limit)
	if err != nil {
		return favs, err
	}

	defer rows.Close()

	for rows.Next() {
		fav, err := mapFav(rows)
		if err != nil {
			return favs, err
		}

		favs = append(favs, fav)
	}

	return favs, nil
}

func mapFav(rows *sql.Rows) (Favorite, error) {
	var id int
	var userId int64
	var userName string
	var tweetId int64
	var status string
	var favDate time.Time
	var unfavDate mysql.NullTime
	var lastAction time.Time

	err := rows.Scan(&id, &userId, &userName, &tweetId, &status, &favDate, &unfavDate, &lastAction)
	if err != nil {
		return Favorite{}, err
	}

	var unfavTime time.Time
	if unfavDate.Valid {
		unfavTime = unfavDate.Time
	}

	return Favorite{id: id, UserId: userId, UserName: userName, TweetId: tweetId, Status: status, FavDate: favDate, UnfavDate: unfavTime, LastAction: lastAction}, nil
}
