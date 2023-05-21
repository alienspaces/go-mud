package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/require"

	"gitlab.com/alienspaces/go-mud/backend/core/collection/set"
	"gitlab.com/alienspaces/go-mud/backend/core/config"
	"gitlab.com/alienspaces/go-mud/backend/core/convert"
	"gitlab.com/alienspaces/go-mud/backend/core/log"
	"gitlab.com/alienspaces/go-mud/backend/core/prepare"
	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/store"
	"gitlab.com/alienspaces/go-mud/backend/core/tag"
	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/core/type/preparer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/storer"
)

type testNestedMultiTag struct {
	A          string   `db:"db_A" json:"json_A"`
	B          []string `db:"db_B" json:"json_B"`
	testNested `db:"nested" json:"NESTED"`
	//lint:ignore U1000 intended to be unused?
	testEmptyNested        `db:"empty_nested" json:"EMPTY_NESTED"`
	testSingleNestedDB     `db:"single_nested_db" json:"SINGLE_NESTED_DB"`
	testSingleNestedJSON   `db:"single_nested_json" json:"SINGLE_NESTED_JSON"`
	singleNestedNullTime   `db:"single_nested_null_time" json:"SINGLE_NESTED_NULL_TIME"`
	singleNestedNullString `db:"single_nested_null_string" json:"SINGLE_NESTED_NULL_STRING"`
	//lint:ignore U1000 intended to be unused?
	f int `db:"db_f"`
	G int `json:"json_G"`
	//lint:ignore U1000 intended to be unused?
	h int
	//lint:ignore U1000 intended to be unused?
	time time.Time `db:"db_time"`
}

type testNested struct {
	C string `json:"json_C"`
	//lint:ignore U1000 intended to be unused?
	d  []string `db:"db_d"`
	E  string   `db:"db_E" json:"json_E"`
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
	l := log.NewLogger(c)

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
					TableName:   "test",
					Attributes:  []string{"test"},
					ArrayFields: set.Set[string]{},
				},
			},
		},
		{
			name: "err - no TableName",
			fields: fields{
				Config: Config{
					Attributes:  []string{"test"},
					ArrayFields: set.Set[string]{},
				},
			},
			wantErr: true,
		},
		{
			name: "err - no Attributes",
			fields: fields{
				Config: Config{
					TableName:   "test",
					ArrayFields: set.Set[string]{},
				},
			},
			wantErr: true,
		},
		{
			name: "err - no ArrayFields",
			fields: fields{
				Config: Config{
					TableName:  "test",
					Attributes: []string{"test"},
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
			TableName:   "not_physically_needed_to_test_init",
			Attributes:  []string{"does_not_matter"},
			ArrayFields: set.Set[string]{},
		},
	}

	err := r.Init()
	require.NoError(t, err, "Repository Init returns without error")

	// test functions of this "prepareable" thing
	testFuncs := map[string]func() string{
		"GetOneSQL":     r.GetOneSQL,
		"GetManySQL":    r.GetManySQL,
		"CreateOneSQL":  r.CreateOneSQL,
		"UpdateOneSQL":  r.UpdateOneSQL,
		"DeleteOneSQL":  r.DeleteOneSQL,
		"DeleteManySQL": r.DeleteManySQL,
		"RemoveOneSQL":  r.RemoveOneSQL,
		"RemoveManySQL": r.RemoveManySQL,
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
			TableName:   "test",
			Attributes:  tag.GetFieldTagValues(testNestedMultiTag{}, "db"),
			ArrayFields: tag.GetArrayFieldTagValues(testNestedMultiTag{}, "db"),
		},
	}

	err := r.Init()
	require.NoError(t, err, "Repository Init returns without error")

	tests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want:          set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL\nAND db_A IN ('1','2')\nAND db_B IN ('a','b')",
					r.Config.TableName),
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE id = $1 AND deleted_at IS NULL\nAND db_B IN ('a','b')\nAND db_A IN ('1','2')",
					r.Config.TableName),
			}),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r.IsRLSDisabled = tt.IsRLSDisabled
			r.SetRLS(tt.identifiers)
			require.Contains(t, tt.want, strings.TrimSpace(r.GetOneSQL()))
		})
	}
}

func TestGetManySQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:   "test",
			Attributes:  tag.GetFieldTagValues(testNestedMultiTag{}, "db"),
			ArrayFields: tag.GetArrayFieldTagValues(testNestedMultiTag{}, "db"),
		},
	}

	err := r.Init()
	require.NoError(t, err, "Repository Init returns without error")

	sqlTests := []struct {
		name          string
		identifiers   map[string][]string
		IsRLSDisabled bool
		want          set.Set[string]
	}{
		{
			name:          "RLS disabled, no identifiers",
			IsRLSDisabled: true,
			want:          set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name:          "RLS disabled, has identifiers",
			IsRLSDisabled: true,
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name: "RLS enabled, no identifiers",
			want: set.FromSlice([]string{fmt.Sprintf("SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL", r.Config.TableName)}),
		},
		{
			name: "RLS enabled, has identifiers",
			identifiers: map[string][]string{
				"db_A": {"1", "2"},
				"db_B": {"a", "b"},
			},
			// appending of the RLS constraints to the SQL query is not ordered because of map iteration
			want: set.FromSlice([]string{
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL\nAND db_A IN ('1','2')\nAND db_B IN ('a','b')",
					r.Config.TableName),
				fmt.Sprintf(
					"SELECT db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time FROM %s WHERE deleted_at IS NULL\nAND db_B IN ('a','b')\nAND db_A IN ('1','2')",
					r.Config.TableName),
			}),
		},
	}

	for _, tt := range sqlTests {
		t.Run(tt.name, func(t *testing.T) {
			r.IsRLSDisabled = tt.IsRLSDisabled
			r.SetRLS(tt.identifiers)
			require.Contains(t, tt.want, strings.TrimSpace(r.GetManySQL()))
		})
	}

	type args struct {
		opts *coresql.Options
	}
	optTests := []struct {
		name string
		args args
		want *coresql.Options
	}{
		{
			name: "@> for multi value, @> for nested slice field for single value, IN with multi value int slice, = for string",
			args: args{
				opts: &coresql.Options{
					Params: []coresql.Param{
						{
							Col: "db_B",
							Val: []string{"a", "b"},
						},
						{
							Col: "db_d",
							Val: []string{"c"},
						},
						{
							Col: "db_f",
							Val: []int{1, 2},
						},
						{
							Col: "db_E6",
							Val: "a",
						},
					},
				},
			},
			want: &coresql.Options{
				Params: []coresql.Param{
					{
						Col:   "db_B",
						Op:    coresql.OpContains,
						Array: convert.GenericSlice([]string{"a", "b"}),
					},
					{
						Col:   "db_d",
						Op:    coresql.OpContains,
						Array: convert.GenericSlice([]string{"c"}),
					},
					{
						Col:   "db_f",
						Op:    coresql.OpIn,
						Array: convert.GenericSlice([]int{1, 2}),
					},
					{
						Col: "db_E6",
						Op:  coresql.OpEqualTo,
						Val: "a",
					},
				},
			},
		},
		{
			name: "ANY, ANY for nested slice field, = for int, IN with single value string slice",
			args: args{
				opts: &coresql.Options{
					Params: []coresql.Param{
						{
							Col: "db_B",
							Val: "a",
						},
						{
							Col: "db_d",
							Val: "c",
						},
						{
							Col: "db_f",
							Val: 1,
						},
						{
							Col: "db_E6",
							Val: []string{"a"},
						},
					},
				},
			},
			want: &coresql.Options{
				Params: []coresql.Param{
					{
						Col: "db_B",
						Op:  coresql.OpAny,
						Val: "a",
					},
					{
						Col: "db_d",
						Op:  coresql.OpAny,
						Val: "c",
					},
					{
						Col: "db_f",
						Op:  coresql.OpEqualTo,
						Val: 1,
					},
					{
						Col:   "db_E6",
						Op:    coresql.OpIn,
						Array: convert.GenericSlice([]string{"a"}),
					},
				},
			},
		},
	}
	for _, tt := range optTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.resolveOpts(tt.args.opts)
			require.NoError(t, err, "resolveOpts should not err")
			require.Equal(t, tt.want, got, "resolveOpts should return coresql.Options as expected")
		})
	}
}

func TestCreateOneSQL(t *testing.T) {
	l, p, tx := newDependencies(t)

	r := Repository{
		Log:     l,
		Prepare: p,
		Tx:      tx,

		Config: Config{
			TableName:   "test",
			Attributes:  tag.GetFieldTagValues(testNestedMultiTag{}, "db"),
			ArrayFields: tag.GetArrayFieldTagValues(testNestedMultiTag{}, "db"),
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
			TableName:   "test",
			Attributes:  tag.GetFieldTagValues(testNestedMultiTag{}, "db"),
			ArrayFields: tag.GetArrayFieldTagValues(testNestedMultiTag{}, "db"),
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
			TableName:   "test",
			Attributes:  tag.GetFieldTagValues(testNestedMultiTag{}, "db"),
			ArrayFields: tag.GetArrayFieldTagValues(testNestedMultiTag{}, "db"),
		},
	}

	require.Equal(t, fmt.Sprintf(`
UPDATE %s SET deleted_at = :deleted_at WHERE id = :id AND deleted_at IS NULL RETURNING db_A, db_B, db_d, db_E, db_E3, db_E5, db_E6, db_f, db_time
`, r.Config.TableName), r.DeleteOneSQL())
}
