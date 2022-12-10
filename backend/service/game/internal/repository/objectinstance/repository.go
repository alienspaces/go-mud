package objectinstance

import (
	"time"

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
	TableName string = "object_instance"
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
				Attributes: tag.GetValues(record.ObjectInstance{}, "db"),
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
func (r *Repository) NewRecord() *record.ObjectInstance {
	return &record.ObjectInstance{}
}

// NewRecordArray -
func (r *Repository) NewRecordArray() []*record.ObjectInstance {
	return []*record.ObjectInstance{}
}

// GetOne -
func (r *Repository) GetOne(id string, forUpdate bool) (*record.ObjectInstance, error) {
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
	forUpdate bool) ([]*record.ObjectInstance, error) {

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

// CreateOne -
func (r *Repository) CreateOne(rec *record.ObjectInstance) error {

	if rec.ID == "" {
		rec.ID = repository.NewRecordID()
	}
	rec.CreatedAt = repository.NewCreatedAt()

	err := r.CreateOneRec(rec)
	if err != nil {
		rec.CreatedAt = time.Time{}
		r.Log.Warn("failed statement execution >%v<", err)
		return err
	}

	return nil
}

// UpdateOne -
func (r *Repository) UpdateOne(rec *record.ObjectInstance) error {

	origUpdatedAt := rec.UpdatedAt
	rec.UpdatedAt = repository.NewUpdatedAt()

	err := r.UpdateOneRec(rec)
	if err != nil {
		rec.UpdatedAt = origUpdatedAt
		r.Log.Warn("failed statement execution >%v<", err)
		return err
	}

	return nil
}
