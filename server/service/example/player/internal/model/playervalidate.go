package model

import (
	"fmt"

	"gitlab.com/alienspaces/go-boilerplate/server/service/player/internal/record"
)

// ValidatePlayerRec - validates creating and updating an player record
func (m *Model) ValidatePlayerRec(rec *record.Player) error {

	switch rec.Provider {
	case record.AccountProviderAnonymous:
		// We only require a provider player ID to create an anonymous local player
		if rec.ProviderAccountID == "" {
			msg := "missing ProviderAccountID, cannot create an anonymous player"
			m.Log.Warn(msg)
			return fmt.Errorf(msg)
		}
	case record.AccountProviderGoogle:
		// We require a provider player ID and email to create a Google local player
		if rec.Email == "" {
			msg := "missing Email, cannot create a Google player"
			m.Log.Warn(msg)
			return fmt.Errorf(msg)
		}
		if rec.ProviderAccountID == "" {
			msg := "missing ProviderAccountID, cannot create a Google player"
			m.Log.Warn(msg)
			return fmt.Errorf(msg)
		}
	default:
		// no-op
	}
	return nil
}

// ValidateDeletePlayerRec - validates it is okay to delete a player record
func (m *Model) ValidateDeletePlayerRec(recID string) error {

	return nil
}
