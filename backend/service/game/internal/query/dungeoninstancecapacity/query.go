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

	rows, err := q.Rows(opts)
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
WITH 
"dungeon_capacity" AS (
    SELECT 
        d.id        AS dungeon_id, 
        count(l.id) AS dungeon_location_count
    FROM dungeon d
    JOIN location l 
        ON l.dungeon_id = d.id
    GROUP BY d.id
), 
"dungeon_instance_capacity" AS (
    SELECT 
        di.id         AS dungeon_instance_id,
        di.dungeon_id AS dungeon_id,
        count(ci.id)  AS dungeon_instance_character_count
    FROM dungeon_instance di
    LEFT JOIN character_instance ci 
        ON ci.dungeon_instance_id = di.id        
    GROUP BY di.id        
)
SELECT 
    dic.dungeon_instance_id, 
    dic.dungeon_instance_character_count,
    dc.dungeon_id, 
    dc.dungeon_location_count
FROM dungeon_instance_capacity dic
JOIN dungeon_capacity dc 
    ON dc.dungeon_id = dic.dungeon_id
    AND dc.dungeon_location_count > dic.dungeon_instance_character_count
WHERE 1 = 1
`
}
