package repository

import (
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/server/core/config"
	"gitlab.com/alienspaces/go-mud/server/core/log"
	"gitlab.com/alienspaces/go-mud/server/core/prepare"
	"gitlab.com/alienspaces/go-mud/server/core/store"
	"gitlab.com/alienspaces/go-mud/server/core/tag"
	"gitlab.com/alienspaces/go-mud/server/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/server/core/type/logger"
	"gitlab.com/alienspaces/go-mud/server/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/server/core/type/storer"
)

type testNestedMultiTag struct {
	A                      string                 `db:"db_A" json:"json_A"`
	B                      string                 `db:"db_B" json:"json_B"`
	TestNested             testNested             `db:"nested" json:"NESTED"`
	TestEmptyNested        testEmptyNested        `db:"empty_nested" json:"EMPTY_NESTED"`
	TestSingleNestedDB     testSingleNestedDB     `db:"single_nested_db" json:"SINGLE_NESTED_DB"`
	TestSingleNestedJSON   testSingleNestedJSON   `db:"single_nested_json" json:"SINGLE_NESTED_JSON"`
	SingleNestedNullTime   singleNestedNullTime   `db:"single_nested_null_time" json:"SINGLE_NESTED_NULL_TIME"`
	SingleNestedNullString singleNestedNullString `db:"single_nested_null_string" json:"SINGLE_NESTED_NULL_STRING"`
	F                      int                    `db:"db_f"`
	G                      int                    `json:"json_G"`
	H                      int
	Time                   time.Time `db:"db_time"`
}

type testNested struct {
	C  string `json:"json_C"`
	D  string `db:"db_d"`
	E  string `db:"db_E" json:"json_E"`
	E2 int
}

type testEmptyNested struct{}

type testSingleNestedDB struct {
	E3 int `db:"db_E3"`
}

type testSingleNestedJSON struct {
	E4 int `json:"json_E4"`
}

type singleNestedNullTime struct {
	E5 sql.NullTime `db:"db_E5"`
}

type singleNestedNullString struct {
	E6 sql.NullString `db:"db_E6"`
}

func NewDependencies() (configurer.Configurer, logger.Logger, storer.Storer, preparer.Repository, error) {

	// configurer
	c, err := config.NewConfigWithDefaults([]config.Item{}, false)
	if err != nil {
		return nil, nil, nil, nil, err
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
	p, err := prepare.NewRepositoryPreparer(l)
	if err != nil {
		return nil, nil, nil, nil, err
	}

	return c, l, s, p, nil
}

func TestInit(t *testing.T) {
	type fields struct {
		Config Config
	}

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ok",
			fields: fields{
				Config: Config{
					TableName:  "test",
					Attributes: []string{"test"},
				},
			},
		},
		{
			name: "err - no TableName",
			fields: fields{
				Config: Config{
					Attributes: []string{"test"},
				},
			},
			wantErr: true,
		},
		{
			name: "err - no Attributes",
			fields: fields{
				Config: Config{
					TableName: "test",
				},
			},
			wantErr: true,
		},
	}

	l, p, tx := newDependencies(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Repository{
				Config:  tt.fields.Config,
				Log:     l,
				Tx:      tx,
				Prepare: p,
			}

			err := r.Init()
			if tt.wantErr {
				require.Error(t, err, "Repository Init returns with error")
			} else {
				require.NoError(t, err, "Repository Init returns without error")
			}
		})
	}
}

func TestPrepareable(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		// Config
		Config: Config{
			TableName:  "not_physically_needed_to_test_init",
			Attributes: []string{"does_not_matter"},
		},
	}

	err := r.Init()
	require.NoError(t, err, "Repository Init returns without error")

	// test functions of this "prepareable" thing
	testFuncs := map[string]func() string{
		"GetOneSQL":          r.GetOneSQL,
		"GetOneForUpdateSQL": r.GetOneForUpdateSQL,
		"GetManySQL":         r.GetManySQL,
		"CreateOneSQL":       r.CreateOneSQL,
		"UpdateOneSQL":       r.UpdateOneSQL,
		"DeleteOneSQL":       r.DeleteOneSQL,
		"RemoveOneSQL":       r.RemoveOneSQL,
	}

	for testFuncName, testFunc := range testFuncs {
		sql := testFunc()
		// Base type should return SQL
		t.Logf("Function %s returns SQL", testFuncName)
		require.NotEmpty(t, sql, fmt.Sprintf("%s returns SQL", testFuncName))
	}

	// test functions of this "prepareable" thing
	testFuncs = map[string]func() string{
		"UpdateManySQL": r.UpdateManySQL,
	}

	for testFuncName, testFunc := range testFuncs {
		t.Run(testFuncName, func(t *testing.T) {
			sql := testFunc()
			// Base type should not return SQL
			t.Logf("Function %s does not return SQL", testFuncName)
			require.Empty(t, sql, fmt.Sprintf("%s returns SQL", testFuncName))
		})
	}
}

func TestGetOneSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	require.Equal(t, fmt.Sprintf(`
SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL
`, r.Config.TableName), r.GetOneSQL())
}

func TestGetOneForUpdateSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	require.Equal(t, fmt.Sprintf(`
SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL FOR UPDATE SKIP LOCKED
`, r.Config.TableName), r.GetOneForUpdateSQL())
}

func TestGetManySQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	require.Equal(t, fmt.Sprintf(`
SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL
`, r.Config.TableName), r.GetManySQL())
}

func TestCreateOneSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	r.Init()

	require.Equal(t, fmt.Sprintf(`
INSERT INTO %s (
	db_A,
	db_B,
	db_d,
	db_E,
	db_E3,
	db_E5,
	db_E6,
	db_f,
	db_time
) VALUES (
	:db_A,
	:db_B,
	:db_d,
	:db_E,
	:db_E3,
	:db_E5,
	:db_E6,
	:db_f,
	:db_time
)
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time
`, r.Config.TableName), r.CreateOneSQL())
}

func TestUpdateOneSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	r.Init()

	require.Equal(t, fmt.Sprintf(`
UPDATE %s SET
	db_A = :db_A,
	db_B = :db_B,
	db_d = :db_d,
	db_E = :db_E,
	db_E3 = :db_E3,
	db_E5 = :db_E5,
	db_E6 = :db_E6,
	db_f = :db_f,
	db_time = :db_time
WHERE id = :id
AND   deleted_at IS NULL
RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time
`, r.Config.TableName), r.UpdateOneSQL())
}

func newDependencies(t *testing.T) (logger.Logger, preparer.Repository, *sqlx.Tx) {
	_, l, s, p, err := NewDependencies()
	require.NoError(t, err, "NewDependencies returns without error")

	err = s.Init()
	require.NoError(t, err, "Store Init returns without error")

	tx, err := s.GetTx()
	require.NoError(t, err, "Store NewTx returns without error")
	return l, p, tx
}

func TestDeleteOneSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:  "test",
			Attributes: tag.GetValues(testNestedMultiTag{}, "db"),
		},
	}

	require.Equal(t, fmt.Sprintf(`
UPDATE %s SET deleted_at = :deleted_at WHERE id = :id AND deleted_at IS NULL RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time
`, r.Config.TableName), r.DeleteOneSQL())
}
