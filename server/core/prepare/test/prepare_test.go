package prepare

import (
	"fmt"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/prepare"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparable"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

func NewDependencies() (configurer.Configurer, logger.Logger, storer.Storer, error) {

	// configurer
	c, err := config.NewConfig([]config.Item{}, false)
	if err != nil {
		return nil, nil, nil, err
	}

	configVars := []string{
		// database
		"APP_SERVER_DB_HOST",
		"APP_SERVER_DB_PORT",
		"APP_SERVER_DB_NAME",
		"APP_SERVER_DB_USER",
		"APP_SERVER_DB_PASSWORD",
	}
	for _, key := range configVars {
		err := c.Add(key, true)
		if err != nil {
			return nil, nil, nil, err
		}
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

func TestInit(t *testing.T) {

	_, l, s, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	// new preparer
	p, err := prepare.NewPrepare(l)
	require.NoError(t, err, "NewPrepare returns without error")
	require.NotNil(t, p, "NewPrepare returns a preparer")

	// get db
	db, err := s.GetDb()
	require.NoError(t, err, "GetDb returns without error")

	// init preparer
	err = p.Init(db)
	require.NoError(t, err, "Init preparer returns without error")
}

func TestPrepare(t *testing.T) {

	// NOTE: Following tests are testing function calls with a successfully
	// prepared "preparable" thing

	// run the following tests within a function to we can utilise
	// a deferred function to teardown any database setup
	func() {
		_, l, s, err := NewDependencies()
		require.NoError(t, err, "NewDependencies returns without error")

		// new preparer
		p, err := prepare.NewPrepare(l)
		require.NoError(t, err, "NewPrepare returns without error")
		require.NotNil(t, p, "NewPrepare returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		r, teardown, err := setupRepository(l, p, db)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Init preparer returns without error")
		require.NotNil(t, r, "Repository is not nil")

		err = p.Prepare(r)
		require.NoError(t, err, "Prepare returns without error")

		// Test SQL functions
		testSQLFuncs := map[string]func(p preparable.Preparable) string{
			"GetOneSQL":     p.GetOneSQL,
			"GetManySQL":    p.GetManySQL,
			"CreateSQL":     p.CreateSQL,
			"UpdateOneSQL":  p.UpdateOneSQL,
			"DeleteOneSQL":  p.DeleteOneSQL,
			"DeleteManySQL": p.DeleteManySQL,
			"RemoveOneSQL":  p.RemoveOneSQL,
			"RemoveManySQL": p.RemoveManySQL,
		}

		for testFuncName, testFunc := range testSQLFuncs {
			t.Logf("Function %s returns SQL", testFuncName)
			// Expecting SQL
			sql := testFunc(r)
			assert.NotEmpty(t, sql, fmt.Sprintf("Function %s returns SQL", testFuncName))
		}

		testSQLFuncs = map[string]func(p preparable.Preparable) string{
			"UpdateManySQL": p.UpdateManySQL,
		}

		for testFuncName, testFunc := range testSQLFuncs {
			t.Logf("Function %s does NOT return SQL", testFuncName)
			// Not expecting SQL
			sql := testFunc(r)
			assert.Empty(t, sql, fmt.Sprintf("Function %s returns SQL", testFuncName))
		}

		// Test Stmt functions
		testStmtFuncs := map[string]func(p preparable.Preparable) *sqlx.Stmt{
			"GetOneStmt":          p.GetOneStmt,
			"GetOneForUpdateStmt": p.GetOneForUpdateStmt,
		}

		for testFuncName, testFunc := range testStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.NotNil(t, stmt, fmt.Sprintf("Function %s returns Stmt", testFuncName))
		}

		// Test NamedStmt functions
		testNamedStmtFuncs := map[string]func(p preparable.Preparable) *sqlx.NamedStmt{
			"GetManyStmt":    p.GetManyStmt,
			"CreateOneStmt":  p.CreateOneStmt,
			"UpdateOneStmt":  p.UpdateOneStmt,
			"UpdateManyStmt": p.UpdateManyStmt,
			"DeleteOneStmt":  p.DeleteOneStmt,
			"DeleteManyStmt": p.DeleteManyStmt,
			"RemoveOneStmt":  p.RemoveOneStmt,
			"RemoveManyStmt": p.RemoveManyStmt,
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
		p, err := prepare.NewPrepare(l)
		require.NoError(t, err, "NewPrepare returns without error")
		require.NotNil(t, p, "NewPrepare returns a preparer")

		// get db
		db, err := s.GetDb()
		require.NoError(t, err, "GetDb returns without error")

		// init preparer
		err = p.Init(db)
		require.NoError(t, err, "Init preparer returns without error")

		r, teardown, err := setupRepository(l, p, db)
		defer func() {
			if teardown != nil {
				teardown()
			}
		}()

		require.NoError(t, err, "Init preparer returns without error")
		require.NotNil(t, r, "Repository is not nil")

		// Test Stmt functions
		testStmtFuncs := map[string]func(p preparable.Preparable) *sqlx.Stmt{
			"GetOneStmt":          p.GetOneStmt,
			"GetOneForUpdateStmt": p.GetOneForUpdateStmt,
		}

		for testFuncName, testFunc := range testStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.Nil(t, stmt, fmt.Sprintf("Function %s returns Stmt", testFuncName))
		}

		// Test NamedStmt functions
		testNamedStmtFuncs := map[string]func(p preparable.Preparable) *sqlx.NamedStmt{
			"GetManyStmt":    p.GetManyStmt,
			"CreateOneStmt":  p.CreateOneStmt,
			"UpdateOneStmt":  p.UpdateOneStmt,
			"UpdateManyStmt": p.UpdateManyStmt,
			"DeleteOneStmt":  p.DeleteOneStmt,
			"DeleteManyStmt": p.DeleteManyStmt,
			"RemoveOneStmt":  p.RemoveOneStmt,
			"RemoveManyStmt": p.RemoveManyStmt,
		}

		for testFuncName, testFunc := range testNamedStmtFuncs {
			t.Logf("Function %s returns stmt", testFuncName)
			stmt := testFunc(r)
			assert.Nil(t, stmt, fmt.Sprintf("Function %s returns NamedStmt", testFuncName))
		}
	}()
}
