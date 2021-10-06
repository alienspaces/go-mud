package record

import (
	"gitlab.com/alienspaces/go-boilerplate/server/core/repository"
)

// AccountProvider - Valid account providers
const (
	AccountProviderAnonymous string = "anonymous"
	AccountProviderApple     string = "apple"
	AccountProviderFacebook  string = "facebook"
	AccountProviderGithub    string = "github"
	AccountProviderGoogle    string = "google"
	AccountProviderTwitter   string = "twitter"
)

// Player -
type Player struct {
	repository.Record
	Name              string `db:"name"`
	Email             string `db:"email"`
	Provider          string `db:"provider"`
	ProviderAccountID string `db:"provider_account_id"`
}

// PlayerRole - Valid roles
const (
	PlayerRoleDefault       string = "default"
	PlayerRoleAdministrator string = "administrator"
)

// PlayerRole -
type PlayerRole struct {
	repository.Record
	PlayerID string `db:"player_id"`
	Role     string `db:"role"`
}
