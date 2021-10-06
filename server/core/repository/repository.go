// Package repository provides methods for interacting with the database
package repository

import (
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/repositor"
)

// Repository -
type Repository struct {
	Config       Config
	Log          logger.Logger
	Tx           *sqlx.Tx
	Prepare      preparer.Preparer
	RecordParams map[string]*RecordParam
}

var _ repositor.Repositor = &Repository{}

// Config -
type Config struct {
	TableName string
}

// RecordParam -
type RecordParam struct {
	TypeInt        bool
	TypeString     bool
	TypeNullString bool
}

// Init -
func (r *Repository) Init(p preparer.Preparer, tx *sqlx.Tx) error {

	r.Log.Debug("Initialising repository %s", r.TableName())

	if p != nil {
		r.Prepare = p
	}

	if tx != nil {
		r.Tx = tx
	}

	if r.Tx == nil {
		return errors.New("Tx is nil, cannot initialise")
	}

	if r.Prepare == nil {
		return errors.New("Prepare is nil, cannot initialise")
	}

	return nil
}

// TableName -
func (r *Repository) TableName() string {
	return r.Config.TableName
}

// GetOneRec -
func (r *Repository) GetOneRec(recordID string, rec interface{}, forUpdate bool) error {

	// preparer
	p := r.Prepare

	// stmt
	var stmt *sqlx.Stmt

	if forUpdate {
		stmt = p.GetOneForUpdateStmt(r)
	} else {
		stmt = p.GetOneStmt(r)
	}

	r.Log.Debug("Get record ID >%s<", recordID)

	stmt = r.Tx.Stmtx(stmt)

	err := stmt.QueryRowx(recordID).StructScan(rec)
	if err != nil {
		r.Log.Warn("Failed executing query >%v<", err)
		r.Log.Warn("SQL: >%s<", p.GetOneSQL(r))
		r.Log.Warn("recordID: >%v<", recordID)

		rec = nil

		return err
	}

	r.Log.Debug("Record fetched")

	return nil
}

// GetManyRecs -
func (r *Repository) GetManyRecs(params map[string]interface{}, operators map[string]string) (rows *sqlx.Rows, err error) {

	// preparer
	p := r.Prepare

	// stmt
	querySQL := p.GetManySQL(r)

	// tx
	tx := r.Tx

	// params
	querySQL, queryParams, err := r.sqlFromParamsAndOperator(querySQL, params, operators)
	if err != nil {
		r.Log.Debug("Failed generating query >%v<", err)
		return nil, err
	}

	r.Log.Debug("Query >%s<", querySQL)
	r.Log.Debug("Parameters >%+v<", queryParams)

	rows, err = tx.NamedQuery(querySQL, queryParams)
	if err != nil {
		r.Log.Warn("Failed querying row >%v<", err)
		return nil, err
	}

	return rows, nil
}

// CreateOneRec -
func (r *Repository) CreateOneRec(rec interface{}) error {

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.CreateOneStmt(r)

	stmt = r.Tx.NamedStmt(stmt)

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		r.Log.Warn("Failed executing create >%v<", err)
		return err
	}

	return nil
}

// UpdateOneRec -
func (r *Repository) UpdateOneRec(rec interface{}) error {

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.UpdateOneStmt(r)

	stmt = r.Tx.NamedStmt(stmt)

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		r.Log.Warn("Failed executing update >%v<", err)
		return err
	}

	return nil
}

// DeleteOne -
func (r *Repository) DeleteOne(id string) error {
	return r.deleteOneRec(id)
}

func (r *Repository) deleteOneRec(recordID string) error {

	params := map[string]interface{}{
		"id":         recordID,
		"deleted_at": NewDeletedAt(),
	}

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.DeleteOneStmt(r)

	stmt = r.Tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		r.Log.Warn("Failed executing delete >%v<", err)
		return err
	}

	// rows affected
	raf, err := res.RowsAffected()
	if err != nil {
		r.Log.Warn("Failed executing rows affected >%v<", err)
		return err
	}

	// expect a single row
	if raf != 1 {
		return fmt.Errorf("Expecting to delete exactly one row but deleted >%d<", raf)
	}

	r.Log.Debug("Deleted >%d< records", raf)

	return nil
}

// RemoveOne -
func (r *Repository) RemoveOne(id string) error {
	return r.removeOneRec(id)
}

func (r *Repository) removeOneRec(recordID string) error {

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.RemoveOneStmt(r)

	params := map[string]interface{}{
		"id": recordID,
	}

	stmt = r.Tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		r.Log.Warn("Failed executing remove >%v<", err)
		return err
	}

	// rows affected
	raf, err := res.RowsAffected()
	if err != nil {
		r.Log.Warn("Failed executing rows affected >%v<", err)
		return err
	}

	// expect a single row
	if raf != 1 {
		return fmt.Errorf("Expecting to remove exactly one row but removed >%d<", raf)
	}

	r.Log.Debug("Removed >%d< records", raf)

	return nil
}

// GetOneSQL -
func (r *Repository) GetOneSQL() string {
	return fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deleted_at IS NULL", r.TableName())
}

// GetOneForUpdateSQL -
func (r *Repository) GetOneForUpdateSQL() string {
	return fmt.Sprintf("SELECT * FROM %s WHERE id = $1 AND deleted_at IS NULL FOR UPDATE SKIP LOCKED", r.TableName())
}

// GetManySQL -
func (r *Repository) GetManySQL() string {
	sql := `
SELECT * FROM %s WHERE deleted_at IS NULL
`
	return fmt.Sprintf(sql, r.TableName())
}

// CreateOneSQL -
func (r *Repository) CreateOneSQL() string {
	return ""
}

// UpdateOneSQL -
func (r *Repository) UpdateOneSQL() string {
	return ""
}

// UpdateManySQL -
func (r *Repository) UpdateManySQL() string {
	return ""
}

// DeleteOneSQL -
func (r *Repository) DeleteOneSQL() string {
	return fmt.Sprintf("UPDATE %s SET deleted_at = :deleted_at WHERE id = :id AND deleted_at IS NULL RETURNING *", r.TableName())
}

// DeleteManySQL -
func (r *Repository) DeleteManySQL() string {
	sql := `
UPDATE %s SET deleted_at = :deleted_at WHERE deleted_at IS NULL
`
	return fmt.Sprintf(sql, r.TableName())
}

// RemoveOneSQL -
func (r *Repository) RemoveOneSQL() string {
	return fmt.Sprintf("DELETE FROM %s WHERE id = :id", r.TableName())
}

// RemoveManySQL -
func (r *Repository) RemoveManySQL() string {
	sql := `
DELETE FROM %s WHERE 1 = 1
`
	return fmt.Sprintf(sql, r.TableName())
}
