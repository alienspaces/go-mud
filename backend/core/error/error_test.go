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
				e: GetRegistryError(ErrorCodeValidationJSON),
			},
			want: GetRegistryError(ErrorCodeValidationJSON),
		},
		{
			name: "invalid query param",
			args: args{
				e: GetRegistryError(ErrorCodeValidationQueryParam),
			},
			want: GetRegistryError(ErrorCodeValidationQueryParam),
		},
		{
			name: "invalid path param",
			args: args{
				e: GetRegistryError(ErrorCodeValidationPathParam),
			},
			want: GetRegistryError(ErrorCodeValidationPathParam),
		},
		{
			name: "unauthenticated",
			args: args{
				e: GetRegistryError(ErrorCodeClientUnauthenticated),
			},
			want: GetRegistryError(ErrorCodeClientUnauthenticated),
		},
		{

			name: "unauthorized",
			args: args{
				e: GetRegistryError(ErrorCodeClientUnauthorized),
			},
			want: GetRegistryError(ErrorCodeClientUnauthorized),
		},
		{
			name: "not found",
			args: args{
				e: GetRegistryError(ErrorCodeResourceNotFound),
			},
			want: GetRegistryError(ErrorCodeResourceNotFound),
		},
		{
			name: "internal",
			args: args{
				e: GetRegistryError(ErrorCodeServerInternal),
			},
			want: GetRegistryError(ErrorCodeServerInternal),
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
