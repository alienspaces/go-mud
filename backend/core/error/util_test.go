package error

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestToError(t *testing.T) {
	type args struct {
		e error
	}
	tests := []struct {
		name    string
		args    args
		want    Error
		wantErr bool
	}{
		{
			name: "invalid JSON",
			args: args{
				e: NewInvalidJSONError(""),
			},
			want: invalidJSON,
		},
		{
			name: "invalid param",
			args: args{
				e: NewParamError(""),
			},
			want: invalidParam,
		},
		{
			name: "unauthenticated",
			args: args{
				e: NewUnauthenticatedError(""),
			},
			want: unauthenticated,
		},
		{

			name: "unauthorized",
			args: args{
				e: NewUnauthorizedError(),
			},
			want: unauthorized,
		},
		{
			name: "not found",
			args: args{
				e: NewNotFoundError("", ""),
			},
			want: notFound,
		},
		{
			name: "internal",
			args: args{
				e: NewInternalError(),
			},
			want: internal,
		},
		{
			name: "error",
			args: args{
				e: fmt.Errorf("error"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToError(tt.args.e)
			if tt.wantErr {
				require.NotNil(t, err, "should not be able to convert error type to coreerror.Error")
				require.Zero(t, got, "coreerror.Error should have zero value for error type that cannot be converted to coreerror.Error")
			} else {
				require.Nil(t, err, "should be able to convert error type to coreerror.Error")
				require.Equal(t, tt.want, got, "coreerror.Error converted should be the same as expected")
			}
		})
	}
}
