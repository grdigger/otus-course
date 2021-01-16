package repository

import (
	"database/sql"
)

type FriendRepository interface {
	Add(userID, friendID int) (err error)
	Remove(userID, friendID int) (err error)
	GetFriends(userID int) (map[int]string, error)
}
type Friend struct {
	db *sql.DB
}

func NewFriend(db *sql.DB) FriendRepository {
	return &Friend{db: db}
}

func (f *Friend) Add(userID, friendID int) (err error) {
	_, err = f.db.Exec("INSERT INTO friend (user_id, friend_id) VALUES (?, ?)", userID, friendID)
	if err != nil {
		return err
	}

	return nil
}

func (f *Friend) Remove(userID, friendID int) (err error) {
	// по хорошему надо сделать soft_delete
	_, err = f.db.Exec("delete from friend where user_id = ? and friend_id = ?", userID, friendID)
	if err != nil {
		return err
	}
	return nil
}

func (f *Friend) GetFriends(userID int) (map[int]string, error) {
	out := make(map[int]string)
	rows, err := f.db.Query("select f.friend_id, concat(u.surname, ' ',u.name) as name from friend f join user u on f.friend_id = u.user_id where f.user_id = ? ", userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i int
		var n string
		// for each row, scan the result into our tag composite object
		err = rows.Scan(&i, &n)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}