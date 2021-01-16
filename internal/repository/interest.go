package repository

import (
	"context"
	"database/sql"
)

type InterestRepository interface {
	GetByUserId(id int) (out []int, err error)
	UpdateInterestsById(id int, interests []string) (bool, error)
}

type Interest struct {
	db *sql.DB
}

func NewInterest(db *sql.DB) InterestRepository {
	return &Interest{db: db}
}

func (i *Interest) GetByUserId(id int) (out []int, err error) {
	out = make([]int, 0)
	rows, err := i.db.Query("select interest_id from user_interest where user_id = ? order by interest_id", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var iID int
		err = rows.Scan(&iID)
		if err != nil {
			return nil, err
		}
		out = append(out, iID)
	}
	return out, err
}

func (i *Interest) UpdateInterestsById(id int, interests []string) (bool, error) {
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return false, err
	}
	_, err = tx.ExecContext(ctx, "delete from user_interest where user_id = ?", id)
	if err != nil {
		tx.Rollback()
		return false, err
	}

	for _, v := range interests {
		_, err = tx.ExecContext(ctx, "INSERT INTO user_interest (user_id, interest_id) VALUES (?, ?)", id, v)
		if err != nil {
			tx.Rollback()
			return false, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}
