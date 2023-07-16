package objectinstanceview

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/tag"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	// TableName - underlying database table name used for configuration
	TableName string = "object_instance_view"
)

// Repository -
type Repository struct {
	repository.Repository
}

var _ repositor.Repositor = &Repository{}

// NewRepository -
func NewRepository(l logger.Logger, p preparer.Repository, tx *sqlx.Tx) (*Repository, error) {

	r := &Repository{
		repository.Repository{
			Log:     l,
			Prepare: p,
			Tx:      tx,

			// Config
			Config: repository.Config{
				TableName:   TableName,
				Attributes:  tag.GetFieldTagValues(record.ObjectInstanceView{}, "db"),
				ArrayFields: tag.GetArrayFieldTagValues(record.ObjectInstanceView{}, "db"),
			},
		},
	}

	err := r.Init()
	if err != nil {
		l.Warn("failed new repository >%v<", err)
		return nil, err
	}

	// prepare
	err = p.Prepare(r, preparer.ExcludePreparation{
		CreateOne:  true,
		CreateMany: true,
		UpdateOne:  true,
		UpdateMany: true,
		DeleteOne:  true,
		DeleteMany: true,
		RemoveOne:  true,
		RemoveMany: true,
	})
	if err != nil {
		l.Warn("failed preparing repository >%v<", err)
		return nil, err
	}

	return r, nil
}

// NewRecord -
func (r *Repository) NewRecord() *record.ObjectInstanceView {
	return &record.ObjectInstanceView{}
}

// NewRecordArray -
func (r *Repository) NewRecordArray() []*record.ObjectInstanceView {
	return []*record.ObjectInstanceView{}
}

// GetOne -
func (r *Repository) GetOne(id string, lock *coresql.Lock) (*record.ObjectInstanceView, error) {
	rec := r.NewRecord()
	if err := r.GetOneRec(id, rec, lock); err != nil {
		r.Log.Warn("failed statement execution >%v<", err)
		return nil, err
	}
	return rec, nil
}

// GetMany -
func (r *Repository) GetMany(opts *coresql.Options) ([]*record.ObjectInstanceView, error) {

	recs := r.NewRecordArray()

	rows, err := r.GetManyRecs(opts)
	if err != nil {
		r.Log.Warn("failed statement execution >%v<", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		rec := r.NewRecord()
		err := rows.StructScan(rec)
		if err != nil {
			r.Log.Warn("failed executing struct scan >%v<", err)
			return nil, err
		}
		recs = append(recs, rec)
	}

	r.Log.Debug("fetched >%d< records", len(recs))

	return recs, nil
}
