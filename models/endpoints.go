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
type EntityEndpoint struct {
	*model.Prototype `dive:""`
	Path             string `db:"path"`
}

func (e *EntityEndpoint) GetAbsEndpoint() (*AbsEndpoint, error) {
	return &AbsEndpoint{
		ID:        *e.ID,
		Path:      e.Path,
		CreatedAt: e.Ctime.Format("2006-01-02 15:04:05"),
		UpdatedAt: e.Mtime.Format("2006-01-02 15:04:05"),
	}, nil
}

type AbsEndpoint struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Path      string    `db:"path" json:"path"`
	CreatedAt string    `db:"ctime" json:"createdAt"`
	UpdatedAt string    `db:"mtime" json:"updatedAt"`
}

//================================================================
// Engine
//================================================================
type EndpointsTableEngine struct {
	*model.Engine
}

func NewEndpointsTableEngine(db *sqlx.DB) *EndpointsTableEngine {
	return &EndpointsTableEngine{
		Engine: model.NewEngine(db, "endpoints"),
	}
}

func (e *EndpointsTableEngine) Insert(path string) (*EntityEndpoint, error) {
	ee := &EntityEndpoint{
		Prototype: model.NewPrototype(),
		Path:      path,
	}

	_, err := e.Engine.Insert(ee)
	return ee, err
}

func (e *EndpointsTableEngine) GetByID(id string) (*EntityEndpoint, error) {
	row := EntityEndpoint{}
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

func (e *EndpointsTableEngine) List(pg *model.Pagination) ([]*AbsEndpoint, error) {
	rows := []*AbsEndpoint{}
	if pg == nil {
		pg = model.NewDefaultPagination()
	}

	q := `SELECT * FROM ` + e.TblName + ` ` + pg.ToString() + `;`
	errDB := e.Select(&rows, q)
	return rows, errDB
}
