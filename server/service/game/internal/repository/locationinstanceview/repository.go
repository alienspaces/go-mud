package locationinstanceview

import (
	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/server/core/repository"
	"gitlab.com/alienspaces/go-mud/server/core/tag"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/server/service/game/internal/record"
)

const (
	// TableName - underlying database table name used for configuration
	TableName string = "location_instance_view"
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
				TableName:  TableName,
				Attributes: tag.GetValues(record.LocationInstanceView{}, "db"),
			},
		},
	}

	err := r.Init()
	if err != nil {
		l.Warn("failed new repository >%v<", err)
		return nil, err
	}

	// prepare
	err = p.Prepare(r, preparer.ExcludePreparation{})
	if err != nil {
		l.Warn("failed preparing repository >%v<", err)
		return nil, err
	}

	return r, nil
}

// NewRecord -
func (r *Repository) NewRecord() *record.LocationInstanceView {
	return &record.LocationInstanceView{}
}

// NewRecordArray -
func (r *Repository) NewRecordArray() []*record.LocationInstanceView {
	return []*record.LocationInstanceView{}
}

// GetOne -
func (r *Repository) GetOne(id string, forUpdate bool) (*record.LocationInstanceView, error) {
	rec := r.NewRecord()
	if err := r.GetOneRec(id, rec, forUpdate); err != nil {
		r.Log.Warn("failed statement execution >%v<", err)
		return nil, err
	}
	return rec, nil
}

// GetMany -
func (r *Repository) GetMany(
	params map[string]interface{},
	paramOperators map[string]string,
	forUpdate bool) ([]*record.LocationInstanceView, error) {

	recs := r.NewRecordArray()

	rows, err := r.GetManyRecs(params, paramOperators, forUpdate)
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
