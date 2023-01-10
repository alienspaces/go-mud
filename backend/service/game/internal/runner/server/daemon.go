package runner

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

func (rnr *Runner) getDungeonInstanceRecs() ([]*record.DungeonInstance, error) {
	l := loggerWithContext(rnr.Log, "getDungeonInstanceRecs")

	m, err := rnr.initModeller(l)
	if err != nil {
		l.Warn("failed initialising modeller >%v<", err)
		return nil, err
	}

	diRecs, err := m.GetDungeonInstanceRecs(nil, nil, false)
	if err != nil {
		l.Warn("failed getting dungeon instance records >%v<", err)
		return nil, err
	}

	err = m.Rollback()
	if err != nil {
		l.Warn("failed model rollback >%v<", err)
		return nil, err
	}

	return diRecs, nil
}

type processState string

const (
	processStatePending processState = "pending"
	processStateRunning processState = "running"
	processStateDone    processState = "done"
	processStateError   processState = "error"
)

type dungeonInstanceState struct {
	turn  int
	state processState
	err   error
}

type dungeonInstanceProcessingResult struct {
	Error             error
	DungeonInstanceID string
	Turn              int
}

func mergeDungeonInstanceStates(dungeonInstanceStates map[string]*dungeonInstanceState, diRecs []*record.DungeonInstance) map[string]*dungeonInstanceState {
	for idx := range diRecs {
		if _, ok := dungeonInstanceStates[diRecs[idx].ID]; !ok {
			dungeonInstanceStates[diRecs[idx].ID] = &dungeonInstanceState{
				state: processStatePending,
			}
		}
	}
	return dungeonInstanceStates
}

// RunDaemon - Starts the daemon process. Override to implement a custom daemon run function.
// The daemon process is a long running background process intended to listen or poll for events
// and then process those events.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {
	l := loggerWithContext(rnr.Log, "RunDaemon")

	dungeonInstanceStates := make(map[string]*dungeonInstanceState)

	c := make(chan dungeonInstanceProcessingResult, 1)

	for keepRunning() {

		// Fetch all current dungeon instance records
		diRecs, err := rnr.getDungeonInstanceRecs()
		if err != nil {
			l.Warn("failed getting dungeon instance recs >%v<", err)
			return err
		}

		// Merge any new dungeon instance records with existing dungeon instance states
		dungeonInstanceStates = mergeDungeonInstanceStates(dungeonInstanceStates, diRecs)

		// When there is nothing to process, wait and check for new instances again
		if len(dungeonInstanceStates) == 0 {
			time.Sleep(1000 * time.Millisecond)
			continue
		}

		runningCount := 0
		for dungeonInstanceID := range dungeonInstanceStates {
			switch dungeonInstanceStates[dungeonInstanceID].state {
			case processStatePending:
				l.Debug("(pending) Kicking off dungeon instance ID >%s< turn >%d<", dungeonInstanceID, dungeonInstanceStates[dungeonInstanceID].turn)
				dungeonInstanceStates[dungeonInstanceID].state = processStateRunning

				go func(dungeonInstanceID string) {

					var m *model.Model

					handleErr := func(err error) {
						l.Warn("failed with error >%v<", err)
						if m != nil {
							rerr := m.Rollback()
							if rerr != nil {
								l.Warn("failed model rollback >%v<", rerr)
								err = fmt.Errorf("%v with %v", rerr, err)
							}
						}
						c <- dungeonInstanceProcessingResult{
							DungeonInstanceID: dungeonInstanceID,
							Error:             err,
						}
					}

					m, err := rnr.initModeller(l)
					if err != nil {
						handleErr(err)
						return
					}

					result, err := processDungeonInstanceTurn(l, m, dungeonInstanceID)
					if err != nil {
						handleErr(err)
						return
					}

					if result.Incremented {
						err = m.Commit()
						if err != nil {
							handleErr(err)
							return
						}
					} else {
						err = m.Rollback()
						if err != nil {
							handleErr(err)
							return
						}
					}

					c <- dungeonInstanceProcessingResult{
						DungeonInstanceID: dungeonInstanceID,
						Turn:              result.Record.TurnCount,
					}

				}(dungeonInstanceID)

				runningCount++

			case processStateError:
				l.Warn("(error) Removing dungeon instance ID >%s< from processing >%v<", dungeonInstanceID, dungeonInstanceStates[dungeonInstanceID].err)
				delete(dungeonInstanceStates, dungeonInstanceID)
			case processStateDone:
				l.Debug("(done) Enqueuing dungeon instance ID >%s< turn >%d<", dungeonInstanceID, dungeonInstanceStates[dungeonInstanceID].turn)
				dungeonInstanceStates[dungeonInstanceID].state = processStatePending
			default:
				// no-op
			}
		}

		// Wait for a result from one of the routines
		if runningCount > 0 {
			result := <-c
			if result.Error != nil {
				dungeonInstanceStates[result.DungeonInstanceID].state = processStateError
				dungeonInstanceStates[result.DungeonInstanceID].err = result.Error
				continue
			}

			dungeonInstanceStates[result.DungeonInstanceID].state = processStateDone
			dungeonInstanceStates[result.DungeonInstanceID].turn = result.Turn
		}
	}

	return nil
}

func processDungeonInstanceTurn(l logger.Logger, m *model.Model, dungeonInstanceID string) (*model.IncrementDungeonInstanceTurnResult, error) {
	l = loggerWithContext(l, "processDungeonInstanceTurn")

	var iditResult *model.IncrementDungeonInstanceTurnResult
	var err error

WHILE_RESULT_NOT_INCREMENTED:
	for iditResult == nil || !iditResult.Incremented {
		// Increment turn
		iditResult, err = m.IncrementDungeonInstanceTurn(dungeonInstanceID)
		if err != nil {
			l.Warn("failed incrementing dungeon ID >%s< instance turn >%v<", dungeonInstanceID, err)
			return nil, err
		}

		if !iditResult.Incremented && iditResult.WaitMilliseconds > 0 {
			l.Debug("Sleeping for >%d< milliseconds", iditResult.WaitMilliseconds)
			time.Sleep(time.Duration(iditResult.WaitMilliseconds) * time.Millisecond)
			continue WHILE_RESULT_NOT_INCREMENTED
		}

		// Process monster instances
		recs, err := m.GetMonsterInstanceRecs(map[string]interface{}{
			"dungeon_instance_id": dungeonInstanceID,
		}, nil, false)
		if err != nil {
			l.Warn("failed getting dungeon ID >%s< monster instance records >%v<", dungeonInstanceID, err)
			return nil, err
		}

		l.Info("Processing turn >%d< with >%d< monster instance records", iditResult.Record.TurnCount, len(recs))

		// TODO: Process monster command
		for idx := range recs {
			l.Info("Processing monster instance ID >%s< monster ID >%s<", recs[idx].ID, recs[idx].MonsterID)
			dmicResult, err := m.DecideMonsterAction(recs[idx].ID)
			if err != nil {
				l.Warn("failed deciding monster instance ID >%s< action >%v<", recs[idx].ID, err)
				return nil, err
			}

			l.Info("Monster instance ID >%s< Sentence >%s<", recs[idx].ID, dmicResult.Sentence)
		}
	}

	l.Debug("Processed dungeon instance ID >%s< turn >%d<", dungeonInstanceID, iditResult.Record.TurnCount)

	return iditResult, nil
}

func keepRunning() bool {
	return true
}
