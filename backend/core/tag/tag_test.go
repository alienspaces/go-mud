package tag

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type nestedMultiTag struct {
	A      string `db:"db_A" json:"json_A"`
	B      string `db:"db_B" json:"json_B"`
	nested `db:"nested" json:"NESTED"`
	// emptyNested            `db:"empty_nested" json:"EMPTY_NESTED"`
	singleNestedDB         `db:"single_nested_db" json:"SINGLE_NESTED_DB"`
	singleNestedJSON       `db:"single_nested_json" json:"SINGLE_NESTED_JSON"`
	singleNestedNullTime   `db:"single_nested_null_time" json:"SINGLE_NESTED_NULL_TIME"`
	singleNestedNullString `db:"single_nested_null_string" json:"SINGLE_NESTED_NULL_STRING"`
	f                      int `db:"db_f"`
	G                      int `json:"json_G"`
	h                      int
	// time                   time.Time `db:"db_time"`
}

type nested struct {
	C  string `json:"json_C"`
	d  string `db:"db_d"`
	E  string `db:"db_E" json:"json_E"`
	E2 int
}

// type emptyNested struct{}

type singleNestedDB struct {
	E3 int `db:"db_E3"`
}

type singleNestedJSON struct {
	E4 int `json:"json_E4"`
}

type singleNestedNullTime struct {
	E5 sql.NullTime `db:"db_E5"`
}

type singleNestedNullString struct {
	E6 sql.NullString `db:"db_E6"`
}

func TestGetValues(t *testing.T) {
	entity := nestedMultiTag{
		A: "l",
		B: "m",
		nested: nested{
			C: "n",
			d: "o",
			E: "p",
		},
		f: 1,
		G: 2,
		h: 3,
	}

	type args struct {
		entity interface{}
		tag    string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nested multi-tag - db tag",
			args: args{
				entity: entity,
				tag:    "db",
			},
			want: []string{"db_A", "db_B", "db_d", "db_E", "db_E3", "db_E5", "db_E6", "db_f", "db_time"},
		},
		{
			name: "nested multi-tag - json tag",
			args: args{
				entity: entity,
				tag:    "json",
			},
			want: []string{"json_A", "json_B", "json_C", "json_E", "json_E4", "json_G"},
		},
		{
			name: "nested multi-tag - non-existent tag",
			args: args{
				entity: entity,
				tag:    "zzz",
			},
			want: []string{},
		},
		{
			name: "nested multi-tag - empty tag",
			args: args{
				entity: entity,
				tag:    "",
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFieldTagValues(tt.args.entity, tt.args.tag)
			require.Equal(t, tt.want, got, fmt.Sprintf("getValues - %s", tt.name))
		})
	}
}
