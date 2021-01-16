package repository

import (
	"database/sql"
	"fmt"
	"github.com/grdigger/otus-course/internal/model"
)

type HelperInterface interface {
	GetAll(helperName string) model.Helper
}

type HelperRepo struct {
	helperName string
	db         *sql.DB
}

func NewHelperRepo(db *sql.DB) *HelperRepo {
	return &HelperRepo{db: db}
}

func (h *HelperRepo) GelAll(helperHame string) (*model.Helper, error) {
	out := model.NewHelper()
	outData := make(map[int]string)
	rows, err := h.db.Query(fmt.Sprintf("select id, name from %s", helperHame))
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var id int
		var name string
		// for each row, scan the result into our tag composite object
		err = rows.Scan(&id, &name)
		if err != nil {
			return nil, err
		}
		outData[id] = name
	}
	out.SetData(outData)
	return out, err
}
