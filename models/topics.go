package models

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/hexcraft-biz/model"
	"github.com/jmoiron/sqlx"
)

//================================================================
// Data Struct
//================================================================
type EntityTopic struct {
	*model.Prototype `dive:""`
	Name             string `db:"name"`
}

func (e *EntityTopic) GetAbsTopic() (*AbsTopic, error) {
	return &AbsTopic{
		ID:        *e.ID,
		Name:      e.Name,
		CreatedAt: e.Ctime.Format("2006-01-02 15:04:05"),
		UpdatedAt: e.Mtime.Format("2006-01-02 15:04:05"),
	}, nil
}

type AbsTopic struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	CreatedAt string    `db:"ctime" json:"createdAt"`
	UpdatedAt string    `db:"mtime" json:"updatedAt"`
}

//================================================================
// Engine
//================================================================
type TopicsTableEngine struct {
	*model.Engine
}

func NewTopicsTableEngine(db *sqlx.DB) *TopicsTableEngine {
	return &TopicsTableEngine{
		Engine: model.NewEngine(db, "topics"),
	}
}

func (e *TopicsTableEngine) Insert(name string) (*EntityTopic, error) {
	ee := &EntityTopic{
		Prototype: model.NewPrototype(),
		Name:      name,
	}

	_, err := e.Engine.Insert(ee)
	return ee, err
}

func (e *TopicsTableEngine) GetByID(id string) (*EntityTopic, error) {
	row := EntityTopic{}
	q := `SELECT * FROM ` + e.TblName + ` WHERE id = UUID_TO_BIN(?);`
	if err := e.Engine.Get(&row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func (e *TopicsTableEngine) List(pg *model.Pagination) ([]*AbsTopic, error) {
	rows := []*AbsTopic{}
	if pg == nil {
		pg = model.NewDefaultPagination()
	}

	q := `SELECT * FROM ` + e.TblName + ` ` + pg.ToString() + `;`
	errDB := e.Select(&rows, q)
	return rows, errDB
}
