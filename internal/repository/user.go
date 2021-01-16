package repository

import (
	"context"
	"database/sql"
	"github.com/grdigger/otus-course/internal/model"
)

type UserRepository interface {
	Auth(email, password string) (user *model.User, err error)
	GetByID(id int) (user *model.User, err error)
	UpdateById(id int, name, surname string, age int, gender, city int) (bool, error)
	Add(email, password, name, surname string, age, gender_id, city_id int) (int, error)
}
type User struct {
	db *sql.DB
}

func NewUser(db *sql.DB) UserRepository {
	return &User{db: db}
}

func (u *User) Auth(email, password string) (user *model.User, err error) {
	out := new(model.User)
	rows, err := u.db.Query("select a.id, u.name, u.surname, u.age, u.gender_id, u.city_id from auth a left join user u on (a.id = u.user_id) where email = ? and password = md5(?) limit 1", email, password)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		model := new(model.User)
		// for each row, scan the result into our tag composite object
		err = rows.Scan(&model.ID, &model.Name, &model.Surname, &model.Age, &model.GenderId, &model.CityId)
		if err != nil {
			return nil, err
		}
		out = model
	}
	return out, err
}

func (u *User) GetByID(id int) (user *model.User, err error) {
	out := new(model.User)
	rows, err := u.db.Query("select u.name, u.surname, u.age, u.gender_id, u.city_id, u.user_id from user u where user_id = ? limit 1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		model := new(model.User)
		// for each row, scan the result into our tag composite object
		err = rows.Scan(&model.Name, &model.Surname, &model.Age, &model.GenderId, &model.CityId, &model.ID)
		if err != nil {
			return nil, err
		}
		out = model
	}
	return out, err
}

func (u *User) UpdateById(id int, name, surname string, age int, gender, city int) (bool, error) {
	_, err := u.db.Exec("update user set name = ?, surname = ?, age = ?, gender_id = ?, city_id = ? where user_id = ?", name, surname, age, gender, city, id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (u *User) Add(email, password, name, surname string, age, genderId, cityId int) (int, error) {
	ctx := context.Background()
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}
	stmt, err := tx.Prepare("INSERT INTO auth (email, password) VALUES (?, md5(?))")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	res, err := stmt.Exec(email, password)
	lid, err := res.LastInsertId()

	stmt, err = tx.Prepare("INSERT INTO user (name, surname, age, gender_id, city_id, user_id) VALUES (?, ?, ?, ?, ?, ?);")
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	_, err = stmt.Exec(name, surname, age, genderId, cityId, lid)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return int(lid), nil
}
