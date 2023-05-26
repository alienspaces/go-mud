package dungeoninstancecapacity

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/query"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/tag"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	QueryName string = "dungeoninstancecapacity_query"
)

type Query struct {
	query.Query
}

var _ querier.Querier = &Query{}

func NewQuery(l logger.Logger, p preparer.Query, tx *sqlx.Tx) (*Query, error) {

	q := &Query{
		query.Query{
			Log:     l,
			Tx:      tx,
			Prepare: p,
			Config: query.Config{
				Name:        QueryName,
				ArrayFields: tag.GetArrayFieldTagValues(record.DungeonInstanceCapacity{}, "db"),
			},
		},
	}

	err := q.Init()
	if err != nil {
		l.Error("Failed new query >%v<", err)
		return nil, err
	}

	err = q.Prepare.Prepare(q)
	if err != nil {
		l.Error("failed preparing query >%v<", err)
		return nil, err
	}

	return q, nil
}

// NewRecord -
func (q *Query) NewRecord() *record.DungeonInstanceCapacity {
	return &record.DungeonInstanceCapacity{}
}

// NewRecordArray -
func (q *Query) NewRecordArray() []*record.DungeonInstanceCapacity {
	return []*record.DungeonInstanceCapacity{}
}

// GetMany -
func (q *Query) GetMany(opts *coresql.Options) ([]*record.DungeonInstanceCapacity, error) {

	recs := q.NewRecordArray()

	rows, err := q.GetRows(opts)
	if err != nil {
		q.Log.Warn("failed query execution >%v<", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rec := q.NewRecord()
		err := rows.StructScan(rec)
		if err != nil {
			q.Log.Warn("failed executing struct scan >%v<", err)
			return nil, err
		}
		recs = append(recs, rec)
	}

	q.Log.Debug("fetched >%d< records", len(recs))

	return recs, nil
}

func (q *Query) SQL() string {
	return `
SELECT 
    dungeon_instance_id, 
    dungeon_instance_character_count,
    dungeon_id, 
    dungeon_location_count
FROM dungeon_instance_capacity_view
WHERE 1 = 1
`
}
