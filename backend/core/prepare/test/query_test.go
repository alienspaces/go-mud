package prepare

import (
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/prepare"
	"gitlab.com/alienspaces/go-mud/server/core/query"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparable"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
)

func setupQuery(l logger.Logger, pq preparer.Query, db *sqlx.DB) (preparable.Query, func() error, error) {

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

	q := Query{
		query.Query{
			Log:      l,
			Preparer: pq,
			Tx:       nil,
		},
	}

	return &q, teardown, nil
}

type Query struct {
	query.Query
}

func (q *Query) SQL() string {
	return "SELECT * FROM test;"
}

func TestQueryInit(t *testing.T) {

	_, l, s, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	// new preparer.Query
	q, err := prepare.NewQueryPreparer(l)
	require.NoError(t, err, "NewPrepare returns without error")
	require.NotNil(t, q, "NewPrepare returns a preparer")

	// get db
	db, err := s.GetDb()
	require.NoError(t, err, "GetDb returns without error")

	// init preparer.Query
	err = q.Init(db)
	require.NoError(t, err, "Init preparer returns without error")
}

func TestQueryPrepare(t *testing.T) {

	// NOTE: Following tests are testing function calls with a successfully
	// prepared "preparableQuery" thing

	// Run the following tests within a function so we can utilise
	// a deferred function to teardown any database setup
	func() {
		_, l, s, err := NewDependencies()
		require.NoError(t, err, "NewDependencies returns without error")

		// new preparer.Query
		p, err := prepare.NewQueryPreparer(l)
		require.NoError(t, err, "NewQuerier returns without error")
		require.NotNil(t, p, "NewQuerier returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer.Query
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		q, teardown, err := setupQuery(l, p, db)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Init preparer returns without error")

		err = p.Prepare(q)
		require.NoError(t, err, "Prepare returns without error")

		// sql := p.SQL()
		// assert.NotEmpty(t, sql, "Function query returns SQL")

		stmt := p.Stmt(q)
		assert.NotNil(t, stmt, "Function stmt returns NamedStmt")
	}()

	// NOTE: Following tests are testing function calls with an unprepared "preparable" thing

	// Run the following tests within a function so we can utilise
	// a deferred function to teardown any database setup
	func() {
		_, l, s, err := NewDependencies()
		require.NoError(t, err, "NewDependencies returns without error")

		// new preparer.Query
		p, err := prepare.NewQueryPreparer(l)
		require.NoError(t, err, "NewPrepare returns without error")
		require.NotNil(t, p, "NewPrepare returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		q, teardown, err := setupQuery(l, p, db)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Init preparer returns without error")
		require.NotNil(t, q, "Query is not nil")

		stmt := p.Stmt(q)
		assert.Nil(t, stmt, "Function stmt returns NamedStmt")
	}()
}
