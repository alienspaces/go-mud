package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {

	tests := map[string]struct {
		setup   func() func()
		dotEnv  bool
		items   []Item
		wantErr bool
	}{
		"NewConfig with items": {
			setup:  func() func() { return func() {} },
			dotEnv: false,
			items: []Item{
				Item{
					Key:      "HOME",
					Required: true,
				},
			},
			wantErr: false,
		},
		"NewConfig without items": {
			setup:   func() func() { return func() {} },
			dotEnv:  false,
			items:   nil,
			wantErr: false,
		},
		"NewConfig without dot env": {
			setup: func() func() {
				appHome := os.Getenv("APP_SERVER_HOME")
				os.Setenv("APP_SERVER_HOME", "./")
				return func() {
					os.Setenv("APP_SERVER_HOME", appHome)
				}
			},
			dotEnv:  true,
			items:   nil,
			wantErr: true,
		},
	}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		func() {

			teardown := tc.setup()
			defer func() {
				teardown()
			}()

			e, err := NewConfig(tc.items, tc.dotEnv)
			if tc.wantErr {
				require.Error(t, err, "NewConfig returns with error")
				return
			}
			require.NoError(t, err, "NewConfig returns without error")
			require.NotNil(t, e, "NewConfig returns environment object")
		}()
	}
}

func TestGet(t *testing.T) {

	tests := map[string]struct {
		items      []Item
		wantErr    bool
		wantValues []string
	}{
		"Get valid environment value": {
			items: []Item{
				Item{
					Key:      "HOME",
					Required: true,
				},
			},
			wantErr: false,
			wantValues: []string{
				os.Getenv("HOME"),
			},
		},
		"Get invalid environment value": {
			items: []Item{
				Item{
					Key:      "WORK",
					Required: true,
				},
			},
			wantErr:    true,
			wantValues: []string{},
		}}

	for tcName, tc := range tests {

		t.Logf("Running test >%s<", tcName)

		e, err := NewConfig(tc.items, false)
		if tc.wantErr {
			require.Error(t, err, "NewConfig returns with error")
			continue
		}
		for idx, item := range tc.items {
			value := e.Get(item.Key)
			require.Equal(t, tc.wantValues[idx], value, "Get returns expected value")
		}
	}
}
