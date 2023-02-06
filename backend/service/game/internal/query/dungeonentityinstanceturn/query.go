package dungeonentityinstanceturn

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/query"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/querier"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	QueryName string = "dungeonentityinstanceturn_query"
)

type Query struct {
	query.Query
}

var _ querier.Querier = &Query{}

func NewQuery(l logger.Logger, p preparer.Query, tx *sqlx.Tx) (*Query, error) {

	q := &Query{
		query.Query{
			Log:      l,
			Tx:       tx,
			Preparer: p,
			Config: query.Config{
				Name: QueryName,
			},
		},
	}

	err := q.Init()
	if err != nil {
		l.Error("Failed new query >%v<", err)
		return nil, err
	}

	err = q.Preparer.Prepare(q)
	if err != nil {
		l.Error("failed preparing query >%v<", err)
		return nil, err
	}

	return q, nil
}

// NewRecord -
func (q *Query) NewRecord() *record.DungeonEntityInstanceTurn {
	return &record.DungeonEntityInstanceTurn{}
}

// NewRecordArray -
func (q *Query) NewRecordArray() []*record.DungeonEntityInstanceTurn {
	return []*record.DungeonEntityInstanceTurn{}
}

// GetMany -
func (q *Query) GetMany(
	params map[string]interface{},
	paramOperators map[string]string) ([]*record.DungeonEntityInstanceTurn, error) {

	recs := q.NewRecordArray()

	rows, err := q.GetRows(q.SQL(), params, paramOperators)
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
	dungeon_name,
	dungeon_instance_turn_number,
	entity_type,
	entity_instance_id,
	entity_name,
	entity_instance_turn_number
FROM dungeon_entity_instance_turn_view
WHERE 1 = 1
`
}
