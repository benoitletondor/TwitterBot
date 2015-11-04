/*
 *   Copyright 2015 Benoit LETONDOR
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package db

import (
	"time"
)

type Reply struct {
	id        int
	UserId    int64
	UserName  string
	TweetId   int64
	Status    string
	Answer    string
	ReplyDate time.Time
}

const (
	_TABLE_REPLY = "reply"
)

func (reply Reply) Persist() error {
	stmtIns, err := database.Prepare("INSERT INTO " + _TABLE_REPLY + "(userId, userName, tweetId, status, answer, replyDate) VALUES( ?, ?, ?, ?, ?, ? )")
	if err != nil {
		return err
	}

	defer stmtIns.Close()

	_, err = stmtIns.Exec(reply.UserId, reply.UserName, reply.TweetId, reply.Status, reply.Answer, reply.ReplyDate)
	return err
}

func HasAlreadyReplied(tweetId int64) (bool, error) {
	stmtOut, err := database.Prepare("SELECT count(*) FROM " + _TABLE_REPLY + " WHERE tweetId = ? LIMIT 1")
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
