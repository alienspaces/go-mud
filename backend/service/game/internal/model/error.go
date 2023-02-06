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
	ErrorCodeCharacterNameTaken     coreerror.ErrorCode = "character.name_taken"
)

func NewCharacterNameTakenError(rec *record.Character) error {
	msg := fmt.Sprintf("character name >%s< has been taken", rec.Name)
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeCharacterNameTaken,
		Message:        msg,
	}
}

func NewActionInvalidError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalid,
		Message:        message,
	}
}

func NewActionInvalidDirectionError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidDirection,
		Message:        message,
	}
}

func NewActionInvalidTargetError(message string) error {
	return coreerror.Error{
		HttpStatusCode: http.StatusBadRequest,
		ErrorCode:      ErrorCodeActionInvalidTarget,
		Message:        message,
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