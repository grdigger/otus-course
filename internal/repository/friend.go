package repository

import (
	"database/sql"
)

type FriendRepository interface {
	Add(userID, friendID int) (err error)
	Remove(userID, friendID int) (err error)
	GetFriends(userID int) (map[int]string, error)
	Find(name, surname string) (map[int]string, error)
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

func (f *Friend) Find(name, surname string) (map[int]string, error) {
	out := make(map[int]string)
	var params []interface{}
	q := ""

	if name == "" || surname == "" {
		q = "select u.user_id, concat(u.name, ' ', u.surname) as name from user u order by user_id limit 50"
	} else {
		q = "select u.user_id, concat(u.name, ' ', u.surname) as name from user u where name like ? and surname like ? order by user_id"
		params = []interface{}{"%" + name + "%", "%" + surname + "%"}
	}
	rows, err := f.db.Query(q, params...)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var i int
		var n string
		err = rows.Scan(&i, &n)
		if err != nil {
			return nil, err
		}
		out[i] = n
	}
	return out, nil
}
