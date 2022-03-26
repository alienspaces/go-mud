package set

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDifference(t *testing.T) {
	type args struct {
		setA map[string]struct{}
		setB map[string]struct{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				setA: FromStrSlice([]string{}),
				setB: FromStrSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "empty setA",
			args: args{
				setA: FromStrSlice([]string{}),
				setB: FromStrSlice([]string{"a", "b", "c"}),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "empty setB",
			args: args{
				setA: FromStrSlice([]string{"a", "b", "c"}),
				setB: FromStrSlice([]string{}),
			},
			want: []string{"a", "b", "c"},
		},
		{
			name: "no difference",
			args: args{
				setA: FromStrSlice([]string{"a", "b", "c"}),
				setB: FromStrSlice([]string{"a", "b", "c"}),
			},
			want: nil,
		},
		{
			name: "difference in setA",
			args: args{
				setA: FromStrSlice([]string{"a", "b", "d"}),
				setB: FromStrSlice([]string{"a", "b"}),
			},
			want: []string{"d"},
		},
		{
			name: "difference in setB",
			args: args{
				setA: FromStrSlice([]string{"a", "b"}),
				setB: FromStrSlice([]string{"a", "b", "DIFFERENT"}),
			},
			want: []string{"DIFFERENT"},
		},
	}
	for _, tt := range tests {
		got := Difference(tt.args.setA, tt.args.setB)
		require.Equal(t, FromStrSlice(tt.want), FromStrSlice(got), tt.name)
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
		require.Equal(t, FromStrSlice(tt.want), FromStrSlice(got), tt.name)
	}
}

func TestFromStrSlice(t *testing.T) {
	type args struct {
		items []string
	}
	tests := []struct {
		name string
		args args
		want map[string]struct{}
	}{
		{
			name: "empty",
			args: args{items: []string{}},
			want: map[string]struct{}{},
		},
		{
			name: "not empty",
			args: args{items: []string{"a", "b", "a"}},
			want: map[string]struct{}{
				"a": {},
				"b": {},
			},
		},
	}
	for _, tt := range tests {
		got := FromStrSlice(tt.args.items)
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
		want map[string]struct{}
	}{
		{
			name: "empty",
			args: args{items: &[]string{}},
			want: map[string]struct{}{},
		},
		{
			name: "not empty",
			args: args{items: &[]string{"a", "b", "a"}},
			want: map[string]struct{}{
				"a": {},
				"b": {},
			},
		},
	}
	for _, tt := range tests {
		got := FromStrSlicePtr(tt.args.items)
		require.Equal(t, tt.want, got, tt.name)
	}
}

func TestIntersection(t *testing.T) {
	type args struct {
		setA map[string]struct{}
		setB map[string]struct{}
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "empty",
			args: args{
				setA: FromStrSlice([]string{}),
				setB: FromStrSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "empty setA",
			args: args{
				setA: FromStrSlice([]string{}),
				setB: FromStrSlice([]string{"a", "b", "c"}),
			},
			want: nil,
		},
		{
			name: "empty setB",
			args: args{
				setA: FromStrSlice([]string{"a", "b", "c"}),
				setB: FromStrSlice([]string{}),
			},
			want: nil,
		},
		{
			name: "no intersection",
			args: args{
				setA: FromStrSlice([]string{"a", "b"}),
				setB: FromStrSlice([]string{"c", "d"}),
			},
			want: nil,
		},
		{
			name: "intersection",
			args: args{
				setA: FromStrSlice([]string{"a", "b", "DIFFERENT"}),
				setB: FromStrSlice([]string{"a", "b"}),
			},
			want: []string{"a", "b"},
		},
		{
			name: "intersection",
			args: args{
				setA: FromStrSlice([]string{"a", "b"}),
				setB: FromStrSlice([]string{"a", "b", "DIFFERENT"}),
			},
			want: []string{"a", "b"},
		},
	}
	for _, tt := range tests {
		got := Intersection(tt.args.setA, tt.args.setB)
		require.Equal(t, FromStrSlice(tt.want), FromStrSlice(got), tt.name)
	}
}

func TestIsCompleteBipartiteGraph(t *testing.T) {
	type args struct {
		setA map[string]map[string]struct{}
		setB map[string]struct{}
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
				setA: map[string]map[string]struct{}{
					"1": FromStrSlice([]string{"a", "b"}),
					"2": FromStrSlice([]string{"a"}),
					"3": FromStrSlice([]string{"b", "d"}),
					"4": FromStrSlice([]string{"b"}),
				},
				setB: FromStrSlice([]string{"a", "b"}),
			},
			wantDisconnectedNodeALabels: nil,
			wantDisconnectedNodeBLabels: nil,
		},
		{
			name: "disconnected node A",
			args: args{
				setA: map[string]map[string]struct{}{
					"1":              FromStrSlice([]string{"a", "b"}),
					"2":              FromStrSlice([]string{"a"}),
					"3":              FromStrSlice([]string{"b", "d"}),
					"4":              FromStrSlice([]string{"b"}),
					"DISCONNECTED A": FromStrSlice([]string{"e"}),
				},
				setB: FromStrSlice([]string{"a", "b"}),
			},
			wantDisconnectedNodeALabels: []string{"DISCONNECTED A"},
			wantDisconnectedNodeBLabels: nil,
		},
		{
			name: "disconnected node B",
			args: args{
				setA: map[string]map[string]struct{}{
					"1": FromStrSlice([]string{"a", "b"}),
					"2": FromStrSlice([]string{"a"}),
					"3": FromStrSlice([]string{"b", "d"}),
					"4": FromStrSlice([]string{"b"}),
				},
				setB: FromStrSlice([]string{"a", "b", "DISCONNECTED B"}),
			},
			wantDisconnectedNodeALabels: nil,
			wantDisconnectedNodeBLabels: []string{"DISCONNECTED B"},
		},
		{
			name: "disconnected nodes A and B",
			args: args{
				setA: map[string]map[string]struct{}{
					"1":               FromStrSlice([]string{"a", "b"}),
					"2":               FromStrSlice([]string{"a"}),
					"3":               FromStrSlice([]string{"b", "d"}),
					"4":               FromStrSlice([]string{"b"}),
					"DISCONNECTED A1": FromStrSlice([]string{"e"}),
					"DISCONNECTED A2": FromStrSlice([]string{"e"}),
				},
				setB: FromStrSlice([]string{"a", "b", "DISCONNECTED B1", "DISCONNECTED B2"}),
			},
			wantDisconnectedNodeALabels: []string{"DISCONNECTED A1", "DISCONNECTED A2"},
			wantDisconnectedNodeBLabels: []string{"DISCONNECTED B1", "DISCONNECTED B2"},
		},
	}
	for _, tt := range tests {
		gotDisconnectedNodeALabels, gotDisconnectedNodeBLabels := IsCompleteBipartiteGraph(tt.args.setA, tt.args.setB)
		require.Equal(t, FromStrSlice(tt.wantDisconnectedNodeALabels), FromStrSlice(gotDisconnectedNodeALabels), tt.name)
		require.Equal(t, FromStrSlice(tt.wantDisconnectedNodeBLabels), FromStrSlice(gotDisconnectedNodeBLabels), tt.name)
	}
}
