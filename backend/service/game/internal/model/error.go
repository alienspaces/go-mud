package model

import (
	"fmt"
	"net/http"

	coreerror "gitlab.com/alienspaces/go-mud/backend/core/error"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

const (
	ErrorCodeActionInvalid          coreerror.ErrorCode = "action.invalid"
	ErrorCodeActionInvalidDirection coreerror.ErrorCode = "action.invalid_direction"
	ErrorCodeActionInvalidTarget    coreerror.ErrorCode = "action.invalid_target"
	ErrorCodeActionTooEarly         coreerror.ErrorCode = "action.too_early"
	ErrorCodeActionInvalidCharacter coreerror.ErrorCode = "action.invalid_character"
	ErrorCodeActionInvalidDungeon   coreerror.ErrorCode = "action.invalid_dungeon"
	ErrorCodeCharacterNameTaken     coreerror.ErrorCode = "character.name_taken"
)

func NewInternalError(message string, args ...any) error {
	return coreerror.NewInternalError(message, args...)
}

func NewCharacterNameTakenError(rec *record.Character) error {
	msg := fmt.Sprintf("character name >%s< has been taken", rec.Name)
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeCharacterNameTaken,
		Message:        msg,
	}
}

func NewInvalidActionError(message string, args ...any) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalid,
		Message:        message,
	}
}

func NewInvalidDirectionError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidDirection,
		Message:        message,
	}
}

func NewInvalidTargetError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidTarget,
		Message:        message,
	}
}

func NewActionInvalidCharacterError(characterID string) error {
	msg := fmt.Sprintf("character ID >%s< is dead or missing", characterID)
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidCharacter,
		Message:        msg,
	}
}

func NewActionInvalidDungeonError(dungeonID string) error {
	msg := fmt.Sprintf("dungeon ID >%s< is dead or missing", dungeonID)
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidDungeon,
		Message:        msg,
	}
}

func NewActionTooEarlyError(dungeonInstanceTurnNumber, entityInstanceTurnNumber int) error {
	msg := fmt.Sprintf("dungeon instance turn >%d< is less than or equal to entity instance turn >%d<", dungeonInstanceTurnNumber, entityInstanceTurnNumber)
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionTooEarly,
		Message:        msg,
	}
}
