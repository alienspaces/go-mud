package actioncharacterobject

import (
	"time"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	"gitlab.com/alienspaces/go-mud/backend/core/tag"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/repositor"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	// TableName - underlying database table name used for configuration
	TableName string = "action_character_object"
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
				Attributes: tag.GetValues(record.ActionCharacterObject{}, "db"),
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
func (r *Repository) NewRecord() *record.ActionCharacterObject {
	return &record.ActionCharacterObject{}
}

// NewRecordArray -
func (r *Repository) NewRecordArray() []*record.ActionCharacterObject {
	return []*record.ActionCharacterObject{}
}

// GetOne -
func (r *Repository) GetOne(id string, forUpdate bool) (*record.ActionCharacterObject, error) {
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
	forUpdate bool) ([]*record.ActionCharacterObject, error) {

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
func (r *Repository) CreateOne(rec *record.ActionCharacterObject) error {

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
func (r *Repository) UpdateOne(rec *record.ActionCharacterObject) error {

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