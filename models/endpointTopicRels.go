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
type EntityEndpointTopicRel struct {
	*model.Prototype `dive:""`
	EndpointId       uuid.UUID `db:"endpoint_id"`
	TopicId          uuid.UUID `db:"topic_id"`
}

//================================================================
// View Table: endpoint_topics
//================================================================
type EndpointTopicRel struct {
	ID         uuid.UUID `db:"id" json:"id"`
	EndpointID uuid.UUID `db:"endpoint_id" json:"endpointId"`
	Path       string    `db:"path" json:"path"`
	TopicID    uuid.UUID `db:"topic_id" json:"topicId"`
	Name       string    `db:"name" json:"name"`
	CreatedAt  string    `db:"ctime" json:"createdAt"`
	UpdatedAt  string    `db:"mtime" json:"updatedAt"`
}

//================================================================
// Engine
//================================================================
type EndpointTopicRelsTableEngine struct {
	*model.Engine
	viewName string
}

func NewEndpointTopicRelsTableEngine(db *sqlx.DB) *EndpointTopicRelsTableEngine {
	return &EndpointTopicRelsTableEngine{
		Engine:   model.NewEngine(db, "endpoint_topic_rels"),
		viewName: "view_endpoint_topic_rels",
	}
}

func (e *EndpointTopicRelsTableEngine) Insert(EndpointId, TopicId string) (*EndpointTopicRel, error) {
	eid, _ := uuid.Parse(EndpointId)
	tid, _ := uuid.Parse(TopicId)

	etr := &EntityEndpointTopicRel{
		Prototype:  model.NewPrototype(),
		EndpointId: eid,
		TopicId:    tid,
	}

	if _, err := e.Engine.Insert(etr); err != nil {
		return nil, err
	}

	return e.GetByID(etr.ID.String())
}

func (e *EndpointTopicRelsTableEngine) GetByID(id string) (*EndpointTopicRel, error) {
	row := EndpointTopicRel{}
	q := `SELECT * FROM ` + e.viewName + ` WHERE id = UUID_TO_BIN(?);`
	if err := e.Engine.Get(&row, q, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &row, nil
}

func (e *EndpointTopicRelsTableEngine) GetByEndpointID(endpointId string) ([]*EndpointTopicRel, error) {
	rows := []*EndpointTopicRel{}

	if i, err := uuid.Parse(endpointId); err == nil {
		q := `SELECT * FROM ` + e.viewName + ` WHERE endpoint_id = UUID_TO_BIN(?);`
		errDB := e.Select(&rows, q, i)
		return rows, errDB
	} else {
		return nil, err
	}
}

func (e *EndpointTopicRelsTableEngine) GetByTopicID(topicId string) ([]*EndpointTopicRel, error) {
	rows := []*EndpointTopicRel{}

	if i, err := uuid.Parse(topicId); err == nil {
		q := `SELECT * FROM ` + e.viewName + ` WHERE endpoint_id = UUID_TO_BIN(?);`
		errDB := e.Select(&rows, q, i)
		return rows, errDB
	} else {
		return nil, err
	}
}

func (e *EndpointTopicRelsTableEngine) GetByEndpointPath(path string) ([]*EndpointTopicRel, error) {
	rows := []*EndpointTopicRel{}

	q := `SELECT * FROM ` + e.viewName + ` WHERE path = ?;`
	errDB := e.Select(&rows, q, path)
	return rows, errDB
}
