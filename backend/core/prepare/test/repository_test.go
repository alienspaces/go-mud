package prepare

import (
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	"gitlab.com/alienspaces/go-mud/backend/core/repository"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/core/tag"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparable"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

func setupRepository(l logger.Logger, p preparer.Repository, db *sqlx.DB, tx *sqlx.Tx) (preparable.Repository, func() error, error) {

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
			Tx:      tx,

			// Config
			Config: repository.Config{
				TableName:  "test",
				Attributes: tag.GetValues(Record{}, "db"),
			},
		},
	}

	err = r.Init()
	if err != nil {
		return nil, nil, err
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

	rows, err := r.GetManyRecs(params, paramOperators, forUpdate)
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

func NewDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := config.NewConfigWithDefaults([]config.Item{}, false)
	if err != nil {
		return nil, nil, nil, err
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, err
	}

	return c, l, s, nil
}

func TestPrepareInit(t *testing.T) {

	_, l, s, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	// new preparer
	p, err := prepare.NewRepositoryPreparer(l)
	require.NoError(t, err, "NewRepositoryPreparer returns without error")
	require.NotNil(t, p, "NewRepositoryPreparer returns a preparer")

	// get db
	db, err := s.GetDb()
	require.NoError(t, err, "GetDb returns without error")

	// init preparer
	err = p.Init(db)
	require.NoError(t, err, "Init preparer returns without error")
}

func TestPreparePrepare(t *testing.T) {

	// NOTE: Following tests are testing function calls with a successfully
	// prepared "preparable" thing

	// Run the following tests within a function, so we can utilise
	// a deferred function to teardown any database setup
	func() {
		_, l, s, err := NewDependencies()
		require.NoError(t, err, "NewDependencies returns without error")

		// new preparer
		p, err := prepare.NewRepositoryPreparer(l)
		require.NoError(t, err, "NewRepositoryPreparer returns without error")
		require.NotNil(t, p, "NewRepositoryPreparer returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		// get tx
		tx, err := s.GetTx()
		require.NoError(t, err, "GetTx returns without error")

		r, teardown, err := setupRepository(l, p, db, tx)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Setup repository returns without error")
		require.NotNil(t, r, "Repository is not nil")

		err = p.Prepare(r, preparer.ExcludePreparation{})
		require.NoError(t, err, "Prepare returns without error")

		// Test SQL functions
		testSQLFuncs := map[string]func(p preparable.Repository) string{
			"GetOneSQL":    p.GetOneSQL,
			"CreateSQL":    p.CreateSQL,
			"UpdateOneSQL": p.UpdateOneSQL,
			"DeleteOneSQL": p.DeleteOneSQL,
			"RemoveOneSQL": p.RemoveOneSQL,
		}

		for testFuncName, testFunc := range testSQLFuncs {
			t.Logf("Function %s returns SQL", testFuncName)
			// Expecting SQL
			sql := testFunc(r)
			assert.NotEmpty(t, sql, fmt.Sprintf("Function %s returns SQL", testFuncName))
		}

		// Test Stmt functions
		testStmtFuncs := map[string]func(p preparable.Repository) *sqlx.Stmt{
			"GetOneStmt":          p.GetOneStmt,
			"GetOneForUpdateStmt": p.GetOneForUpdateStmt,
		}

		for testFuncName, testFunc := range testStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.NotNil(t, stmt, fmt.Sprintf("Function %s returns Stmt", testFuncName))
		}

		// Test NamedStmt functions
		testNamedStmtFuncs := map[string]func(p preparable.Repository) *sqlx.NamedStmt{
			"CreateOneStmt": p.CreateOneStmt,
			"UpdateOneStmt": p.UpdateOneStmt,
			"DeleteOneStmt": p.DeleteOneStmt,
			"RemoveOneStmt": p.RemoveOneStmt,
		}

		for testFuncName, testFunc := range testNamedStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.NotNil(t, stmt, fmt.Sprintf("Function %s returns NamedStmt", testFuncName))
		}
	}()

	// NOTE: Following tests are testing function calls with an unprepared "preparable" thing

	// run the following tests within a function to we can utilise
	// a deferred function to teardown any database setup
	func() {
		_, l, s, err := NewDependencies()
		require.NoError(t, err, "NewDependencies returns without error")

		// new preparer
		p, err := prepare.NewRepositoryPreparer(l)
		require.NoError(t, err, "NewRepositoryPreparer returns without error")
		require.NotNil(t, p, "NewRepositoryPreparer returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		// get tx
		tx, err := s.GetTx()
		require.NoError(t, err, "GetTx returns without error")

		r, teardown, err := setupRepository(l, p, db, tx)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Init preparer returns without error")
		require.NotNil(t, r, "Repository is not nil")

		// Test Stmt functions
		testStmtFuncs := map[string]func(p preparable.Repository) *sqlx.Stmt{
			"GetOneStmt":          p.GetOneStmt,
			"GetOneForUpdateStmt": p.GetOneForUpdateStmt,
		}

		for testFuncName, testFunc := range testStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.Nil(t, stmt, fmt.Sprintf("Function %s returns Stmt", testFuncName))
		}

		// Test NamedStmt functions
		testNamedStmtFuncs := map[string]func(p preparable.Repository) *sqlx.NamedStmt{
			"CreateOneStmt": p.CreateOneStmt,
			"UpdateOneStmt": p.UpdateOneStmt,
			"DeleteOneStmt": p.DeleteOneStmt,
			"RemoveOneStmt": p.RemoveOneStmt,
		}

		for testFuncName, testFunc := range testNamedStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.Nil(t, stmt, fmt.Sprintf("Function %s returns NamedStmt", testFuncName))
		}
	}()
}
