// Package repository provides methods for interacting with the database
package repository

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/repositor"
)

// Repository -
type Repository struct {
	Config           Config
	Log              logger.Logger
	Tx               *sqlx.Tx
	Prepare          preparer.Repository
	createAttributes []string
	updateAttributes []string

	attributeIndex         set.Set[string]
	readOnlyAttributeIndex set.Set[string]
}

var _ repositor.Repositor = &Repository{}

// Config -
type Config struct {
	TableName          string
	Attributes         []string
	ReadonlyAttributes []string
	ArrayFields        set.Set[string]
}

// Init -
func (r *Repository) Init() error {
	r.Log.Debug("Initialising repository %s", r.TableName())

	if r.Tx == nil {
		return errors.New("repository Tx is nil, cannot initialise")
	}

	if r.Prepare == nil {
		return errors.New("repository Prepare is nil, cannot initialise")
	}

	if r.TableName() == "" {
		return errors.New("repository TableName is empty, cannot initialise")
	}

	if len(r.Attributes()) == 0 {
		return errors.New("repository Attributes are empty, cannot initialise")
	}

	if r.ArrayFields() == nil {
		return errors.New("repository ArrayFields is nil, cannot initialise")
	}

	readonlyAttributeIndex := map[string]struct{}{}
	for _, attribute := range r.ReadonlyAttributes() {
		readonlyAttributeIndex[attribute] = struct{}{}
	}
	r.readOnlyAttributeIndex = readonlyAttributeIndex

	createAttributes := []string{}
	updateAttributes := []string{}
	attributeIndex := map[string]struct{}{}
	for _, attribute := range r.Attributes() {
		attributeIndex[attribute] = struct{}{}
		if _, ok := readonlyAttributeIndex[attribute]; !ok {
			createAttributes = append(createAttributes, attribute)
		}
		if attribute == "created_at" || attribute == "deleted_at" {
			continue
		}
		if _, ok := readonlyAttributeIndex[attribute]; !ok {
			updateAttributes = append(updateAttributes, attribute)
		}
	}

	r.createAttributes = createAttributes
	r.updateAttributes = updateAttributes
	r.attributeIndex = attributeIndex

	return nil
}

// TableName -
func (r *Repository) TableName() string {
	return r.Config.TableName
}

func (r *Repository) Attributes() []string {
	return r.Config.Attributes
}

func (r *Repository) ReadonlyAttributes() []string {
	return r.Config.ReadonlyAttributes
}

func (r *Repository) ArrayFields() set.Set[string] {
	return r.Config.ArrayFields
}

// GetOneRec -
func (r *Repository) GetOneRec(recordID any, rec any, lock *coresql.Lock) error {

	// preparer
	p := r.Prepare

	// stmt
	querySQL := p.GetOneSQL(r)

	opts := &coresql.Options{
		Lock: lock,
	}

	querySQL, _, err := coresql.From(querySQL, opts)
	if err != nil {
		r.Log.Debug("failed generating query >%v<", err)
		return err
	}

	r.Log.Debug("Get record ID >%s<", recordID)

	err = r.Tx.QueryRowx(querySQL, recordID).StructScan(rec)
	if err != nil {
		r.Log.Warn("Failed executing query >%v<", err)
		r.Log.Warn("SQL: >%s<", querySQL)

		var recID string
		switch id := recordID.(type) {
		case string:
			recID = id
		case []byte:
			recID = base64.URLEncoding.EncodeToString(id)
		default:
			r.Log.Warn("unknown record ID type >%v<", id)
			return err
		}

		r.Log.Warn("recordID: >%v<", recID)

		rec = nil

		return err
	}

	r.Log.Debug("Record fetched")

	return nil
}

// GetManyRecs -
func (r *Repository) GetManyRecs(opts *coresql.Options) (rows *sqlx.Rows, err error) {

	// preparer
	p := r.Prepare

	// stmt
	querySQL := p.GetManySQL(r)

	// tx
	tx := r.Tx

	// params
	opts, err = r.resolveOptions(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve opts: sql >%s< opts >%#v< >%v<", querySQL, opts, err)
	}

	querySQL, queryArgs, err := coresql.From(querySQL, opts)
	if err != nil {
		r.Log.Warn("failed generating query >%v<", err)
		return nil, err
	}

	r.Log.Debug("querySQL >%s<", querySQL)
	r.Log.Debug("queryArgs >%+v<", queryArgs)

	rows, err = tx.NamedQuery(querySQL, queryArgs)
	if err != nil {
		r.Log.Warn("SQL: >%s<", querySQL)
		r.Log.Warn("Failed querying row >%v<", err)
		return nil, err
	}

	return rows, nil
}

func (r *Repository) resolveOptions(opts *coresql.Options) (*coresql.Options, error) {
	if opts == nil {
		return opts, nil
	}

	params := []coresql.Param{}

	for _, p := range opts.Params {

		// Skip parameters that aren't valid attributes for the record
		if _, ok := r.attributeIndex[p.Col]; !ok {
			continue
		}

		switch t := p.Val.(type) {
		case []string:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		case []int:
			p.Array = convert.GenericSlice(t)
			p.Val = nil
		}

		// if Op is specified, it is assumed you know what you're doing
		if p.Op != "" {
			params = append(params, p)
			continue
		}

		isArrayField := r.ArrayFields().Contains(p.Col)
		if isArrayField {
			if len(p.Array) > 0 {
				p.Op = coresql.OpContains
			} else {
				p.Op = coresql.OpAny
			}
		} else {
			if len(p.Array) > 0 {
				p.Op = coresql.OpIn
			} else {
				p.Op = coresql.OpEqual
			}
		}

		params = append(params, p)
	}

	opts.Params = params

	return opts, nil
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
func (r *Repository) DeleteOne(id any) error {
	return r.deleteOneRec(id)
}

func (r *Repository) deleteOneRec(recordID any) error {

	params := map[string]interface{}{
		"id":         recordID,
		"deleted_at": NewRecordNullTimestamp(),
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
		return fmt.Errorf("expecting to delete exactly one row but deleted >%d<", raf)
	}

	r.Log.Debug("Deleted >%d< records", raf)

	return nil
}

// RemoveOne -
func (r *Repository) RemoveOne(id any) error {
	return r.removeOneRec(id)
}

func (r *Repository) removeOneRec(recordID any) error {

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
		return fmt.Errorf("expecting to remove exactly one row but removed >%d<", raf)
	}

	r.Log.Debug("Removed >%d< records", raf)

	return nil
}

// GetOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) GetOneSQL() string {
	return fmt.Sprintf(`
SELECT %s FROM %s WHERE id = $1 AND deleted_at IS NULL
`,
		strings.Join(r.Attributes(), ", "),
		r.TableName())
}

// GetManySQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) GetManySQL() string {
	return fmt.Sprintf(`
SELECT %s FROM %s WHERE deleted_at IS NULL
`,
		strings.Join(r.Attributes(), ", "),
		r.TableName())
}

// CreateOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) CreateOneSQL() string {
	return fmt.Sprintf(`
INSERT INTO %s (
%s
) VALUES (
%s
)
RETURNING %s
`,
		r.TableName(),
		strings.Join(r.createAttributes, ",\n"),
		colonPrefixedCommaNewlineSeparated(r.createAttributes),
		strings.Join(r.Attributes(), ", "))
}

func colonPrefixedCommaNewlineSeparated(attributes []string) string {
	var strBuilder strings.Builder

	for i, a := range attributes {
		strBuilder.WriteString("\t:")
		strBuilder.WriteString(a)

		if i != len(attributes)-1 {
			strBuilder.WriteString(",\n")
		}
	}

	return strBuilder.String()
}

// UpdateOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) UpdateOneSQL() string {

	return fmt.Sprintf(`
UPDATE %s SET
%s
WHERE id = :id
AND   deleted_at IS NULL
RETURNING %s
`,
		r.TableName(),
		equalsAndNewlineSeparated(r.updateAttributes),
		strings.Join(r.Attributes(), ", "))
}

func equalsAndNewlineSeparated(attributes []string) string {
	var strBuilder strings.Builder

	for i, a := range attributes {
		strBuilder.WriteString("\t")
		strBuilder.WriteString(a)
		strBuilder.WriteString(" = :")
		strBuilder.WriteString(a)

		if i != len(attributes)-1 {
			strBuilder.WriteString(",\n")
		}
	}

	return strBuilder.String()
}

// UpdateManySQL -
func (r *Repository) UpdateManySQL() string {
	return ""
}

// DeleteOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) DeleteOneSQL() string {
	return fmt.Sprintf(`
UPDATE %s SET deleted_at = :deleted_at WHERE id = :id AND deleted_at IS NULL RETURNING %s
`,
		r.TableName(),
		strings.Join(r.Attributes(), ", "),
	)
}

// DeleteManySQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) DeleteManySQL() string {
	sql := `
UPDATE %s SET deleted_at = :deleted_at WHERE deleted_at IS NULL
`
	return fmt.Sprintf(sql, r.TableName())
}

// RemoveOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) RemoveOneSQL() string {
	return fmt.Sprintf(`
DELETE FROM %s WHERE id = :id
`, r.TableName())
}

// RemoveManySQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) RemoveManySQL() string {
	sql := `
DELETE FROM %s WHERE 1 = 1
`
	return fmt.Sprintf(sql, r.TableName())
}
