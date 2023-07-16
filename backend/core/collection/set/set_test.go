package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDifference(t *testing.T) {
	type args struct {
		setA Set[string]
		setB Set[string]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				setA: FromSlice([]string{}),
				setB: FromSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "empty setA",
			args: args{
				setA: FromSlice([]string{}),
				setB: FromSlice([]string{"a", "b", "c"}),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty setB",
			args: args{
				setA: FromSlice([]string{"a", "b", "c"}),
				setB: FromSlice([]string{}),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "no difference",
			args: args{
				setA: FromSlice([]string{"a", "b", "c"}),
				setB: FromSlice([]string{"a", "b", "c"}),
			},
			want: nil,
		},
		{
			name: "difference in setA",
			args: args{
				setA: FromSlice([]string{"a", "b", "d"}),
				setB: FromSlice([]string{"a", "b"}),
			},
			want: []string{"d"},
		},
		{
			name: "difference in setB",
			args: args{
				setA: FromSlice([]string{"a", "b"}),
				setB: FromSlice([]string{"a", "b", "DIFFERENT"}),
			},
			want: []string{"DIFFERENT"},
		},
	}
	for _, tt := range tests {
		got := Difference(tt.args.setA, tt.args.setB)
		require.Equal(t, FromSlice(tt.want), FromSlice(got), tt.name)
	}
}

func TestFindDuplicate(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				items: []string{},
			},
			want: nil,
		},
		{
			name: "unique",
			args: args{
				items: []string{"a", "b", "c"},
			},
			want: nil,
		},
		{
			name: "duplicates",
			args: args{
				items: []string{"a", "a", "b", "b", "c"},
			},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		got := FindDuplicates(tt.args.items)
		require.Equal(t, FromSlice(tt.want), FromSlice(got), tt.name)
	}
}

func TestFromStrSlice(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "empty",
			args: args{items: []string{}},
			want: Set[string]{},
		},
		{
			name: "not empty",
			args: args{items: []string{"a", "b", "a"}},
			want: Set[string]{
				"a": {},
				"b": {},
			},
		},
	}
	for _, tt := range tests {
		got := FromSlice(tt.args.items)
		require.Equal(t, tt.want, got, tt.name)
	}
}

func TestFromStrSlicePtr(t *testing.T) {
	type args struct {
		items *[]string
	}
	tests := []struct {
		name string
		args args
		want Set[string]
	}{
		{
			name: "empty",
			args: args{items: &[]string{}},
			want: Set[string]{},
		},
		{
			name: "not empty",
			args: args{items: &[]string{"a", "b", "a"}},
			want: Set[string]{
				"a": {},
				"b": {},
			},
		},
	}
	for _, tt := range tests {
		got := FromSlicePtr(tt.args.items)
		require.Equal(t, tt.want, got, tt.name)
	}
}

func TestIntersection(t *testing.T) {
	type args struct {
		setA Set[string]
		setB Set[string]
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				setA: FromSlice([]string{}),
				setB: FromSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "empty setA",
			args: args{
				setA: FromSlice([]string{}),
				setB: FromSlice([]string{"a", "b", "c"}),
			},
			want: nil,
		},
		{
			name: "empty setB",
			args: args{
				setA: FromSlice([]string{"a", "b", "c"}),
				setB: FromSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "no intersection",
			args: args{
				setA: FromSlice([]string{"a", "b"}),
				setB: FromSlice([]string{"c", "d"}),
			},
			want: nil,
		},
		{
			name: "intersection",
			args: args{
				setA: FromSlice([]string{"a", "b", "DIFFERENT"}),
				setB: FromSlice([]string{"a", "b"}),
			},
			want: []string{"a", "b"},
		},
		{
			name: "intersection",
			args: args{
				setA: FromSlice([]string{"a", "b"}),
				setB: FromSlice([]string{"a", "b", "DIFFERENT"}),
			},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		got := Intersection(tt.args.setA, tt.args.setB)
		require.Equal(t, FromSlice(tt.want), FromSlice(got), tt.name)
	}
}

func TestIsSubset(t *testing.T) {
	type args struct {
		source []string
		target []string
	}

	commonTargetSlice := []string{"a", "b", "c", "d"}

	tests := []struct {
		name           string
		args           args
		expectedResult bool
	}{
		{
			name: "source is a subset of target",
			args: args{
				source: []string{"a", "b", "c"},
				target: commonTargetSlice,
			},
			expectedResult: true,
		},
		{
			name: "source & target are empty slice",
			args: args{
				source: []string{},
				target: []string{},
			},
			expectedResult: true,
		},
		{
			name: "source contains all element in target",
			args: args{
				source: commonTargetSlice,
				target: commonTargetSlice,
			},
			expectedResult: true,
		},
		{
			name: "source only have some elements from target",
			args: args{
				source: []string{"a", "b", "e"},
				target: commonTargetSlice,
			},
			expectedResult: false,
		},
		{
			name: "source doesnt have any elements from target",
			args: args{
				source: []string{"e", "f"},
				target: commonTargetSlice,
			},
			expectedResult: false,
		},
	}

	for _, test := range tests {
		result := IsSubset(test.args.source, test.args.target)
		require.Equal(t, result, test.expectedResult, test.name)
	}
}

func TestIsCompleteBipartiteGraph(t *testing.T) {
	type args struct {
		setA map[string]Set[string]
		setB Set[string]
	}
	tests := []struct {
		name                        string
		args                        args
		wantDisconnectedNodeALabels []string
		wantDisconnectedNodeBLabels []string
	}{
		{
			name: "complete",
			args: args{
				setA: map[string]Set[string]{
					"1": FromSlice([]string{"a", "b"}),
					"2": FromSlice([]string{"a"}),
					"3": FromSlice([]string{"b", "d"}),
					"4": FromSlice([]string{"b"}),
				},
				setB: FromSlice([]string{"a", "b"}),
			},
			wantDisconnectedNodeALabels: nil,
			wantDisconnectedNodeBLabels: nil,
		},
		{
			name: "disconnected node A",
			args: args{
				setA: map[string]Set[string]{
					"1":              FromSlice([]string{"a", "b"}),
					"2":              FromSlice([]string{"a"}),
					"3":              FromSlice([]string{"b", "d"}),
					"4":              FromSlice([]string{"b"}),
					"DISCONNECTED A": FromSlice([]string{"e"}),
				},
				setB: FromSlice([]string{"a", "b"}),
			},
			wantDisconnectedNodeALabels: []string{"DISCONNECTED A"},
			wantDisconnectedNodeBLabels: nil,
		},
		{
			name: "disconnected node B",
			args: args{
				setA: map[string]Set[string]{
					"1": FromSlice([]string{"a", "b"}),
					"2": FromSlice([]string{"a"}),
					"3": FromSlice([]string{"b", "d"}),
					"4": FromSlice([]string{"b"}),
				},
				setB: FromSlice([]string{"a", "b", "DISCONNECTED B"}),
			},
			wantDisconnectedNodeALabels: nil,
			wantDisconnectedNodeBLabels: []string{"DISCONNECTED B"},
		},
		{
			name: "disconnected nodes A and B",
			args: args{
				setA: map[string]Set[string]{
					"1":               FromSlice([]string{"a", "b"}),
					"2":               FromSlice([]string{"a"}),
					"3":               FromSlice([]string{"b", "d"}),
					"4":               FromSlice([]string{"b"}),
					"DISCONNECTED A1": FromSlice([]string{"e"}),
					"DISCONNECTED A2": FromSlice([]string{"e"}),
				},
				setB: FromSlice([]string{"a", "b", "DISCONNECTED B1", "DISCONNECTED B2"}),
			},
			wantDisconnectedNodeALabels: []string{"DISCONNECTED A1", "DISCONNECTED A2"},
			wantDisconnectedNodeBLabels: []string{"DISCONNECTED B1", "DISCONNECTED B2"},
		},
	}
	for _, tt := range tests {
		gotDisconnectedNodeALabels, gotDisconnectedNodeBLabels := IsCompleteBipartiteGraph(tt.args.setA, tt.args.setB)
		require.Equal(t, FromSlice(tt.wantDisconnectedNodeALabels), FromSlice(gotDisconnectedNodeALabels), tt.name)
		require.Equal(t, FromSlice(tt.wantDisconnectedNodeBLabels), FromSlice(gotDisconnectedNodeBLabels), tt.name)
	}
}

func TestSeparate(t *testing.T) {
	type testStruct struct {
		a string
		b string
		c int
	}

	type args struct {
		setA []testStruct
		setB []testStruct
	}
	tests := []struct {
		name        string
		args        args
		wantUniqueA []testStruct
		wantCommon  []testStruct
		wantUniqueB []testStruct
	}{
		{
			name: "identical",
			args: args{
				setA: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "a",
						b: "c",
						c: 0,
					},
				},
				setB: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "a",
						b: "c",
						c: 0,
					},
				},
			},
			wantUniqueA: nil,
			wantCommon: []testStruct{
				{
					a: "a",
					b: "b",
					c: 0,
				},
				{
					a: "a",
					b: "c",
					c: 0,
				},
			},
			wantUniqueB: nil,
		},
		{
			name: "disjoint",
			args: args{
				setA: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "a",
						b: "c",
						c: 0,
					},
				},
				setB: []testStruct{
					{
						a: "a",
						b: "d",
						c: 0,
					},
					{
						a: "a",
						b: "e",
						c: 0,
					},
				},
			},
			wantUniqueA: []testStruct{
				{
					a: "a",
					b: "b",
					c: 0,
				},
				{
					a: "a",
					b: "c",
					c: 0,
				},
			},
			wantCommon: nil,
			wantUniqueB: []testStruct{
				{
					a: "a",
					b: "d",
					c: 0,
				},
				{
					a: "a",
					b: "e",
					c: 0,
				},
			},
		},
		{
			name: "common and unique",
			args: args{
				setA: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "zzzz",
						b: "c",
						c: 1,
					},
				},
				setB: []testStruct{
					{
						a: "a",
						b: "c",
						c: 0,
					},
					{
						a: "a",
						b: "e",
						c: 0,
					},
				},
			},
			wantUniqueA: []testStruct{
				{
					a: "a",
					b: "b",
					c: 0,
				},
			},
			wantCommon: []testStruct{
				{
					a: "zzzz",
					b: "c",
					c: 1,
				},
			},
			wantUniqueB: []testStruct{
				{
					a: "a",
					b: "e",
					c: 0,
				},
			},
		},
		{
			name: "setA empty, setB populated",
			args: args{
				setA: nil,
				setB: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "a",
						b: "c",
						c: 0,
					},
				},
			},
			wantUniqueA: nil,
			wantCommon:  nil,
			wantUniqueB: []testStruct{
				{
					a: "a",
					b: "b",
					c: 0,
				},
				{
					a: "a",
					b: "c",
					c: 0,
				},
			},
		},
		{
			name: "setA populated, setB empty",
			args: args{
				setA: []testStruct{
					{
						a: "a",
						b: "b",
						c: 0,
					},
					{
						a: "a",
						b: "c",
						c: 0,
					},
				},
				setB: nil,
			},
			wantUniqueA: []testStruct{
				{
					a: "a",
					b: "b",
					c: 0,
				},
				{
					a: "a",
					b: "c",
					c: 0,
				},
			},
			wantCommon:  nil,
			wantUniqueB: nil,
		},
	}

	keyFn := func(t testStruct) string {
		return t.b
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uniqueA, common, uniqueB := Separate(tt.args.setA, tt.args.setB, keyFn)
			require.Equal(t, FromSlice(tt.wantUniqueA), FromSlice(uniqueA), "uniqueA return value should equal")
			require.Equal(t, FromSlice(tt.wantCommon), FromSlice(common), "common return value should equal")
			require.Equal(t, FromSlice(tt.wantUniqueB), FromSlice(uniqueB), "uniqueB return value should equal")
		})
	}
}
