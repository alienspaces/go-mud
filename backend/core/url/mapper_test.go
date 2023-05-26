package url

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlugify(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{
			name: "single word",
			arg:  "word",
			want: "word",
		},
		{
			name: "trailing non-word chars",
			arg:  "word     $%^^    ",
			want: "word",
		},
		{
			name: "initial non-word chars",
			arg:  " $%^%$           word",
			want: "word",
		},
		{
			name: "initial & trailing non-word chars",
			arg:  "  $%^$%        word          $%$^",
			want: "word",
		},
		{
			name: "quotes",
			arg:  "''''\"''''''''''''' word        '''''''\"'''''''''",
			want: "word",
		},
		{
			name: "hyphens & underscores",
			arg:  "-___w--__ord_-_-_",
			want: "___w-__ord_-_-_",
		},
		{
			name: "multiword",
			arg:  "123asdf123 asdf123",
			want: "123asdf123-asdf123",
		},
		{
			name: "all",
			arg:  "$@#$@!#$-\"-__123asdf123 $#@!$@#!$  asdf123__--   -- '''\"' $#@!$@#",
			want: "__123asdf123-asdf123__",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Slugify(tt.arg); got != tt.want {
				require.Equal(t, tt.want, got, "slugify should equal")
			}
		})
	}
}
