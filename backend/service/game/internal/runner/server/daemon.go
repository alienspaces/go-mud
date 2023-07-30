package runner

import (
	"fmt"
	"time"

	coresql "gitlab.com/alienspaces/go-mud/backend/core/sql"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

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

func daemonMergeDungeonInstanceStates(dungeonInstanceStates map[string]*dungeonInstanceState, diRecs []*record.DungeonInstance) map[string]*dungeonInstanceState {
	for idx := range diRecs {
		if _, ok := dungeonInstanceStates[diRecs[idx].ID]; !ok {
			dungeonInstanceStates[diRecs[idx].ID] = &dungeonInstanceState{
				state: processStatePending,
			}
		}
	}
	return dungeonInstanceStates
}

func daemonRemoveDungeonInstanceState(dungeonInstanceStates map[string]*dungeonInstanceState, diRec *record.DungeonInstance) map[string]*dungeonInstanceState {

	delete(dungeonInstanceStates, diRec.ID)

	return dungeonInstanceStates
}

func daemonGetDungeonInstanceRecs(l logger.Logger, m *model.Model) ([]*record.DungeonInstance, error) {
	l = loggerWithFunctionContext(l, "daemonGetDungeonInstanceRecs")

	diRecs, err := m.GetDungeonInstanceRecs(nil)
	if err != nil {
		l.Warn("failed getting dungeon instance records >%v<", err)
		return nil, err
	}

	return diRecs, nil
}

func daemonDungeonInstanceEmpty(l logger.Logger, m *model.Model, dir *record.DungeonInstance) (empty bool, err error) {
	l = loggerWithFunctionContext(l, "daemonDungeonInstanceEmpty")

	ciRecs, err := m.GetCharacterInstanceRecs(
		&coresql.Options{
			Params: []coresql.Param{
				{
					Col: "dungeon_instance_id",
					Val: dir.ID,
				},
			},
		},
	)
	if err != nil {
		l.Warn("failed getting many character instance records >%v<", err)
		return empty, err
	}

	l.Debug("Dungeon instance ID >%s< character instance count >%d<", dir.ID, len(ciRecs))

	return len(ciRecs) == 0, nil
}

func daemonShutdownDungeonInstance(l logger.Logger, m *model.Model, dir *record.DungeonInstance) error {
	l = loggerWithFunctionContext(l, "daemonShutdownDungeonInstance")

	err := m.DeleteDungeonInstance(dir.ID)
	if err != nil {
		l.Warn("failed deleting dungeon instance >%v<", err)
		return err
	}

	return nil
}

// daemonInitCycle initialises a new database transaction and must commit or
// rollback before returning. It takes an existing set of dungeon instance states,
// fetches all existing dungeon instances, checks whether they are empty, and updates
// the set of dungeon instance states by removing empty dungeon instances and adding
// new dungeon instances.
func (rnr *Runner) daemonInitCycle(l logger.Logger, dis map[string]*dungeonInstanceState) (map[string]*dungeonInstanceState, error) {
	l = loggerWithFunctionContext(l, "daemonInitCycle")

	m, err := rnr.initModeller(l)
	if err != nil {
		l.Warn("failed initialising modeller >%v<", err)
		return nil, err
	}

	// Fetch all current dungeon instance records
	diRecs, err := daemonGetDungeonInstanceRecs(l, m)
	if err != nil {
		err = m.Rollback()
		if err != nil {
			l.Warn("failed model rollback >%v<", err)
			return nil, err
		}
		l.Warn("failed getting dungeon instance recs >%v<", err)
		return nil, err
	}

	l.Debug("Dungeon instance count >%d<", len(diRecs))

	// When there are no characters instances in a particular dungeon instance
	// for a certain period of time, delete the dungeon instance.
	for idx := range diRecs {
		empty, err := daemonDungeonInstanceEmpty(l, m, diRecs[idx])
		if err != nil {
			err = m.Rollback()
			if err != nil {
				l.Warn("failed model rollback >%v<", err)
				return nil, err
			}
			l.Warn("failed check if dungeon instance ID >%s< is empty >%v<", diRecs[idx].ID, err)
			return nil, err
		}

		if empty {
			err := daemonShutdownDungeonInstance(l, m, diRecs[idx])
			if err != nil {
				err = m.Rollback()
				if err != nil {
					l.Warn("failed model rollback >%v<", err)
					return nil, err
				}
				l.Warn("failed shutting down dungeon instance ID >%s< >%v<", diRecs[idx].ID, err)
				return nil, err
			}
			dis = daemonRemoveDungeonInstanceState(dis, diRecs[idx])
		}
	}

	// Merge any new dungeon instance records with existing dungeon instance states
	dis = daemonMergeDungeonInstanceStates(dis, diRecs)

	err = m.Commit()
	if err != nil {
		l.Warn("failed model commit >%v<", err)
		return nil, err
	}

	return dis, nil
}

// RunDaemon is a long running background process that manages the server game loop.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {
	l := loggerWithFunctionContext(rnr.Log, "RunDaemon")

	dungeonInstanceStates := make(map[string]*dungeonInstanceState)

	c := make(chan dungeonInstanceProcessingResult, 1)

	cycles := 0

	for keepRunning() {

		l = loggerWithCycleContext(l, cycles)

		dungeonInstanceStates, err := rnr.daemonInitCycle(l, dungeonInstanceStates)
		if err != nil {
			l.Warn("failed daemon init cycle >%v<", err)
			return err
		}

		cycles++
		if cycles%20 == 0 {
			l.Info("Daemon cycle >%d<", cycles)
		}

		// When there is nothing to process, wait and check for new instances again
		if len(dungeonInstanceStates) == 0 {
			time.Sleep(3000 * time.Millisecond)
			continue
		}

		runningCount := 0
		for dungeonInstanceID := range dungeonInstanceStates {

			l = loggerWithInstanceContext(l, dungeonInstanceID)

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

					if result.incrementTurnResult.Incremented {
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
						Turn:              result.incrementTurnResult.Record.TurnNumber,
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

type processDungeonInstanceTurnResult struct {
	incrementTurnResult             *model.IncrementDungeonInstanceTurnResult
	monsterInstanceActionRecordSets []*record.ActionRecordSet
}

func processDungeonInstanceTurn(l logger.Logger, m *model.Model, dungeonInstanceID string) (*processDungeonInstanceTurnResult, error) {
	l = loggerWithFunctionContext(l, "processDungeonInstanceTurn")
	l = loggerWithInstanceContext(l, dungeonInstanceID)

	// TODO: 10-implement-effects:
	// Process any active effects that are still applied to the character.

	pditr := processDungeonInstanceTurnResult{}

WHILE_RESULT_NOT_INCREMENTED:
	for pditr.incrementTurnResult == nil || !pditr.incrementTurnResult.Incremented {
		// Increment turn
		iditr, err := m.IncrementDungeonInstanceTurn(&model.IncrementDungeonInstanceTurnArgs{
			DungeonInstanceID: dungeonInstanceID,
		})
		if err != nil {
			l.Warn("failed incrementing dungeon ID >%s< instance turn >%v<", dungeonInstanceID, err)
			return nil, err
		}

		if !iditr.Incremented && iditr.WaitMilliseconds > 0 {
			l.Debug("Sleeping for >%d< milliseconds", iditr.WaitMilliseconds)
			time.Sleep(time.Duration(iditr.WaitMilliseconds) * time.Millisecond)
			continue WHILE_RESULT_NOT_INCREMENTED
		}

		pditr.incrementTurnResult = iditr

		// TODO: 12-implement-death: Remove character instance when dead
		crecs, err := m.GetCharacterInstanceRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "dungeon_instance_id",
						Val: dungeonInstanceID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting dungeon ID >%s< character instance records >%v<", dungeonInstanceID, err)
			return nil, err
		}

		for idx := range crecs {
			l.Info("Processing character instance ID >%s< character ID >%s<", crecs[idx].ID, crecs[idx].CharacterID)

			// 12-implement-death:
			// Remove character instance when dead character has been decaying for 10 turns
			if crecs[idx].Health <= 0 {
				err := m.CharacterExitDungeon(crecs[idx].CharacterID)
				if err != nil {
					l.Warn("failed character instance ID >%s< exit dungeon >%v<", crecs[idx].ID, err)
					return nil, err
				}
			}
		}

		// Process monster instances
		mrecs, err := m.GetMonsterInstanceRecs(
			&coresql.Options{
				Params: []coresql.Param{
					{
						Col: "dungeon_instance_id",
						Val: dungeonInstanceID,
					},
				},
			},
		)
		if err != nil {
			l.Warn("failed getting dungeon ID >%s< monster instance records >%v<", dungeonInstanceID, err)
			return nil, err
		}

		l.Debug("Processing turn >%d< with >%d< monster instance records", iditr.Record.TurnNumber, len(mrecs))

		for idx := range mrecs {
			l.Info("Processing monster instance ID >%s< monster ID >%s< Name >%s<", mrecs[idx].ID, mrecs[idx].MonsterID)
			dmar, err := m.DecideMonsterAction(mrecs[idx].ID)
			if err != nil {
				l.Warn("failed deciding monster instance ID >%s< action >%v<", mrecs[idx].ID, err)
				return nil, err
			}

			l.Info("Monster instance ID >%s< Sentence >%s<", dmar.MonsterInstanceID, dmar.Sentence)

			if dmar.Sentence == "" {
				l.Info("Monster instance ID >%s< not doing anything this turn", dmar.MonsterInstanceID)
				continue
			}

			ars, err := m.ProcessMonsterAction(dmar.DungeonInstanceID, dmar.MonsterInstanceID, dmar.Sentence)
			if err != nil {
				l.Warn("failed processing monster action >%s< action >%v<", dmar.Sentence, err)
				return nil, err
			}

			l.Debug("Processed monster instance ID >%s< action >%#v<", dmar.MonsterInstanceID, ars.ActionRec)
			pditr.monsterInstanceActionRecordSets = append(pditr.monsterInstanceActionRecordSets, ars)
		}
	}

	l.Debug("Processed dungeon instance ID >%s< turn >%d<", dungeonInstanceID, pditr.incrementTurnResult.Record.TurnNumber)

	return &pditr, nil
}

func (rnr *Runner) initModeller(l logger.Logger) (*model.Model, error) {
	m, err := rnr.InitTx(l)
	if err != nil {
		l.Warn("failed initialising database transaction, cannot authen >%v<", err)
		return nil, err
	}
	return m.(*model.Model), nil
}

// keepRunning decides whether the server should continue
// to run based on current state etc..
func keepRunning() bool {
	return true
}

// loggerWithCycleContext provides a logger with function context
func loggerWithCycleContext(l logger.Logger, cycle int) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithContext("cycle", fmt.Sprintf("%d", cycle))
}

// loggerWithInstanceContext provides a logger with function context
func loggerWithInstanceContext(l logger.Logger, instanceID string) logger.Logger {
	if l == nil {
		return nil
	}
	return l.WithContext("instance", instanceID)
}
