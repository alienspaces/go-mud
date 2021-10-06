package prepare

import (
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/repository"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
)

func setupRepository(l logger.Logger, p preparer.Preparer, db *sqlx.DB) (preparable.Preparable, func() error, error) {

	sql := `
CREATE TABLE test (
	id                UUID CONSTRAINT test_pk PRIMARY KEY DEFAULT gen_random_uuid(),
	"name"            TEXT NOT NULL,
	created_at TIMESTAMP WITH TIME ZONE DEFAULT current_timestamp,
	updated_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
	deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);
`
	_, err := db.Exec(sql)
	if err != nil {
		return nil, nil, err
	}

	teardown := func() error {
		sql := `
		DROP TABLE test;
		`
		_, err := db.Exec(sql)
		if err != nil {
			return err
		}
		return nil
	}

	r := Repository{
		repository.Repository{
			Log:     l,
			Prepare: p,
			Tx:      nil,

			// Config
			Config: repository.Config{
				TableName: "test",
			},
		},
	}

	return &r, teardown, nil
}

type Record struct {
	repository.Record
	Name string `db:"name"`
}

type Repository struct {
	repository.Repository
}

// NewRecord -
func (r *Repository) NewRecord() *Record {
	return &Record{}
}

// NewRecordArray -
func (r *Repository) NewRecordArray() []*Record {
	return []*Record{}
}

// GetOne -
func (r *Repository) GetOne(id string, forUpdate bool) (*Record, error) {
	rec := r.NewRecord()
	if err := r.GetOneRec(id, rec, forUpdate); err != nil {
		r.Log.Warn("Failed statement execution >%v<", err)
		return nil, err
	}
	return rec, nil
}

// GetMany -
func (r *Repository) GetMany(
	params map[string]interface{},
	paramOperators map[string]string,
	forUpdate bool) ([]*Record, error) {

	recs := r.NewRecordArray()

	rows, err := r.GetManyRecs(params, paramOperators)
	if err != nil {
		r.Log.Warn("Failed statement execution >%v<", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rec := r.NewRecord()
		err := rows.StructScan(rec)
		if err != nil {
			r.Log.Warn("Failed executing struct scan >%v<", err)
			return nil, err
		}
		recs = append(recs, rec)
	}

	r.Log.Debug("Fetched >%d< records", len(recs))

	return recs, nil
}

// CreateOne -
func (r *Repository) CreateOne(rec *Record) error {

	if rec.ID == "" {
		rec.ID = repository.NewRecordID()
	}
	rec.CreatedAt = repository.NewCreatedAt()

	err := r.CreateOneRec(rec)
	if err != nil {
		rec.CreatedAt = time.Time{}
		r.Log.Warn("Failed statement execution >%v<", err)
		return err
	}

	return nil
}

// UpdateOne -
func (r *Repository) UpdateOne(rec *Record) error {

	origUpdatedAt := rec.UpdatedAt
	rec.UpdatedAt = repository.NewUpdatedAt()

	err := r.UpdateOneRec(rec)
	if err != nil {
		rec.UpdatedAt = origUpdatedAt
		r.Log.Warn("Failed statement execution >%v<", err)
		return err
	}

	return nil
}

// CreateTestRecord - creates a record for testing
func (r *Repository) CreateTestRecord(rec *Record) error {
	return r.CreateOne(rec)
}

var createOneSQL = `
INSERT INTO test (
	id,
	name,
	created_at
) VALUES (
	:id,
	:name,
	:created_at
)
RETURNING *
`

var updateOneSQL = `
UPDATE test SET
  name       = :name,
  updated_at = :updated_at
WHERE id = :id
AND   deleted_at IS NULL
RETURNING *
`

// CreateOneSQL -
func (r *Repository) CreateOneSQL() string {
	return createOneSQL
}

// UpdateOneSQL -
func (r *Repository) UpdateOneSQL() string {
	return updateOneSQL
}

// OrderBy -
func (r *Repository) OrderBy() string {
	return "created_at desc"
}
