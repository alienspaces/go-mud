package repository

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-boilerplate/server/core/config"
	"gitlab.com/alienspaces/go-boilerplate/server/core/log"
	"gitlab.com/alienspaces/go-boilerplate/server/core/prepare"
	"gitlab.com/alienspaces/go-boilerplate/server/core/store"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/preparer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/storer"
)

func NewDependencies() (configurer.Configurer, logger.Logger, storer.Storer, preparer.Preparer, error) {

	// configurer
	c, err := config.NewConfig([]config.Item{}, false)
	if err != nil {
		return nil, nil, nil, nil, err
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
			return nil, nil, nil, nil, err
		}
	}

	// logger
	l, err := log.NewLogger(c)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// storer
	s, err := store.NewStore(c, l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	// preparer
	p, err := prepare.NewPrepare(l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return c, l, s, p, nil
}

func TestInit(t *testing.T) {

	_, l, s, p, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	err = s.Init()
	require.NoError(t, err, "Store Init returns without error")

	tx, err := s.GetTx()
	require.NoError(t, err, "Store NewTx returns without error")

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		// Config
		Config: Config{
			TableName: "some_table_i_need_to_build_as_part_of_setup",
		},
	}

	err = r.Init(p, tx)
	require.NoError(t, err, "Repository Init returns without error")
}

func TestPrepareable(t *testing.T) {

	_, l, s, p, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	err = s.Init()
	require.NoError(t, err, "Store Init returns without error")

	tx, err := s.GetTx()
	require.NoError(t, err, "Store NewTx returns without error")

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		// Config
		Config: Config{
			TableName: "not_physically_needed_to_test_init",
		},
	}

	err = r.Init(p, tx)
	require.NoError(t, err, "Repository Init returns without error")

	// test functions of this "prepareable" thing
	testFuncs := map[string]func() string{
		"GetOneSQL":          r.GetOneSQL,
		"GetOneForUpdateSQL": r.GetOneForUpdateSQL,
		"GetManySQL":         r.GetManySQL,
		"DeleteOneSQL":       r.DeleteOneSQL,
		"DeleteManySQL":      r.DeleteManySQL,
		"RemoveOneSQL":       r.RemoveOneSQL,
		"RemoveManySQL":      r.RemoveManySQL,
	}

	for testFuncName, testFunc := range testFuncs {
		sql := testFunc()
		// Base type should return SQL
		t.Logf("Function %s returns SQL", testFuncName)
		require.NotEmpty(t, sql, fmt.Sprintf("%s returns SQL", testFuncName))
	}

	// test functions of this "prepareable" thing
	testFuncs = map[string]func() string{
		"CreateOneSQL":  r.CreateOneSQL,
		"UpdateOneSQL":  r.UpdateOneSQL,
		"UpdateManySQL": r.UpdateManySQL,
	}

	for testFuncName, testFunc := range testFuncs {
		sql := testFunc()
		// Base type should not return SQL
		t.Logf("Function %s does not return SQL", testFuncName)
		require.Empty(t, sql, fmt.Sprintf("%s returns SQL", testFuncName))
	}
}
