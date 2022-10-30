// Package repository provides methods for interacting with the database
package repository

import (
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	coresql "gitlab.com/alienspaces/go-mud/server/core/sql"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/repositor"
)

// Repository -
type Repository struct {
	Config             Config
	Log                logger.Logger
	Tx                 *sqlx.Tx
	Prepare            preparer.Repository
	RecordParams       map[string]*RecordParam
	computedAttributes []string
}

var _ repositor.Repositor = &Repository{}

// Config -
type Config struct {
	TableName  string
	Attributes []string
}

// RecordParam -
type RecordParam struct {
	TypeInt        bool
	TypeString     bool
	TypeNullString bool
}

// Init -
func (r *Repository) Init() error {
	l := r.Logger("Init")
	l.Debug("Initialising repository %s", r.TableName())

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

	computedAttributes := []string{}
	for _, attribute := range r.Attributes() {
		if attribute == "created_at" ||
			attribute == "deleted_at" {
			continue
		}
		computedAttributes = append(computedAttributes, attribute)
	}
	r.computedAttributes = computedAttributes

	return nil
}

// TableName -
func (r *Repository) TableName() string {
	return r.Config.TableName
}

func (r *Repository) Attributes() []string {
	return r.Config.Attributes
}

// GetOneRec -
func (r *Repository) GetOneRec(recordID string, rec interface{}, forUpdate bool) error {
	l := r.Logger("GetOneRec")
	tx := r.Tx

	// preparer
	p := r.Prepare

	// stmt
	var stmt *sqlx.Stmt

	if forUpdate {
		stmt = p.GetOneForUpdateStmt(r)
	} else {
		stmt = p.GetOneStmt(r)
	}

	l.Debug("Get record ID >%s<", recordID)

	stmt = tx.Stmtx(stmt)

	err := stmt.QueryRowx(recordID).StructScan(rec)
	if err != nil {
		l.Warn("failed executing query >%v<", err)
		l.Warn("SQL: >%s<", p.GetOneSQL(r))
		l.Warn("recordID: >%v<", recordID)

		rec = nil

		return err
	}

	l.Debug("Record fetched")

	return nil
}

// GetManyRecs -
func (r *Repository) GetManyRecs(params map[string]interface{}, operators map[string]string, forUpdate bool) (rows *sqlx.Rows, err error) {
	l := r.Logger("GetManyRecs")
	tx := r.Tx

	querySQL := r.GetManySQL()

	// params
	querySQL, queryParams, err := coresql.FromParamsAndOperators("", querySQL, params, operators)
	if err != nil {
		l.Debug("Failed generating query >%v<", err)
		return nil, err
	}

	if forUpdate {
		querySQL += "FOR UPDATE SKIP LOCKED"
	}

	l.Debug("SQL >%s< Params >%#v<", querySQL, queryParams)

	rows, err = tx.NamedQuery(querySQL, queryParams)
	if err != nil {
		l.Warn("failed querying rows >%v<", err)
		return nil, err
	}

	return rows, nil
}

// CreateOneRec -
func (r *Repository) CreateOneRec(rec interface{}) error {
	l := r.Logger("CreateOneRec")
	tx := r.Tx

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.CreateOneStmt(r)

	stmt = tx.NamedStmt(stmt)

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		l.Warn("failed executing create >%v<", err)
		return err
	}

	return nil
}

// UpdateOneRec -
func (r *Repository) UpdateOneRec(rec interface{}) error {
	l := r.Logger("UpdateOneRec")
	tx := r.Tx

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.UpdateOneStmt(r)

	stmt = tx.NamedStmt(stmt)

	err := stmt.QueryRowx(rec).StructScan(rec)
	if err != nil {
		l.Warn("failed executing update >%v<", err)
		return err
	}

	return nil
}

// DeleteOne -
func (r *Repository) DeleteOne(id string) error {
	return r.deleteOneRec(id)
}

func (r *Repository) deleteOneRec(recordID string) error {
	l := r.Logger("deleteOneRec")
	tx := r.Tx

	params := map[string]interface{}{
		"id":         recordID,
		"deleted_at": NewDeletedAt(),
	}

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.DeleteOneStmt(r)

	stmt = tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		l.Warn("failed executing delete >%v<", err)
		return err
	}

	// rows affected
	raf, err := res.RowsAffected()
	if err != nil {
		l.Warn("failed executing rows affected >%v<", err)
		return err
	}

	// expect a single row
	if raf != 1 {
		return fmt.Errorf("expecting to delete exactly one row but deleted >%d<", raf)
	}

	l.Debug("Deleted >%d< records", raf)

	return nil
}

// RemoveOne -
func (r *Repository) RemoveOne(id string) error {
	return r.removeOneRec(id)
}

func (r *Repository) removeOneRec(recordID string) error {
	l := r.Logger("removeOneRec")
	tx := r.Tx

	// preparer
	p := r.Prepare

	// stmt
	stmt := p.RemoveOneStmt(r)

	params := map[string]interface{}{
		"id": recordID,
	}

	stmt = tx.NamedStmt(stmt)

	res, err := stmt.Exec(params)
	if err != nil {
		l.Warn("failed executing remove >%v<", err)
		return err
	}

	// rows affected
	raf, err := res.RowsAffected()
	if err != nil {
		l.Warn("failed executing rows affected >%v<", err)
		return err
	}

	// expect a single row
	if raf != 1 {
		return fmt.Errorf("expecting to remove exactly one row but removed >%d<", raf)
	}

	l.Debug("Removed >%d< records", raf)

	return nil
}

// GetOneSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) GetOneSQL() string {
	return fmt.Sprintf(`
SELECT %s FROM %s WHERE id = $1 AND deleted_at IS NULL
`,
		commaSeparated(r.Attributes()),
		r.TableName())
}

// GetOneForUpdateSQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) GetOneForUpdateSQL() string {
	return fmt.Sprintf(`
SELECT %s FROM %s WHERE id = $1 AND deleted_at IS NULL FOR UPDATE SKIP LOCKED
`,
		commaSeparated(r.Attributes()),
		r.TableName())
}

// GetManySQL - This SQL statement ends with a newline so that any parameters can be easily appended.
func (r *Repository) GetManySQL() string {
	return fmt.Sprintf(`
SELECT %s FROM %s WHERE deleted_at IS NULL `,
		commaSeparated(r.Attributes()),
		r.TableName())
}

func commaSeparated(attributes []string) string {
	var strBuilder strings.Builder

	for i, a := range attributes {
		strBuilder.WriteString(a)

		if i != len(attributes)-1 {
			strBuilder.WriteString(", ")
		}
	}

	return strBuilder.String()
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
		commaNewlineSeparated(r.Attributes()),
		colonPrefixedCommaNewlineSeparated(r.Attributes()),
		commaSeparated(r.Attributes()))
}

func commaNewlineSeparated(attributes []string) string {
	var strBuilder strings.Builder

	for i, a := range attributes {
		strBuilder.WriteString("\t")
		strBuilder.WriteString(a)

		if i != len(attributes)-1 {
			strBuilder.WriteString(",\n")
		}
	}

	return strBuilder.String()
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
		equalsAndNewlineSeparated(r.computedAttributes),
		commaSeparated(r.Attributes()))
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
		commaSeparated(r.Attributes()),
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

// Logger -
func (r *Repository) Logger(functionName string) logger.Logger {
	return r.Log.WithPackageContext("core/repository").WithInstanceContext(r.TableName()).WithFunctionContext(functionName)
}
