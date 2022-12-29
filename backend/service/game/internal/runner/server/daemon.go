package runner

import (
	"fmt"
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/model"
	"gitlab.com/alienspaces/go-mud/backend/service/game/internal/record"
)

// RunDaemon - Starts the daemon process. Override to implement a custom daemon run function.
// The daemon process is a long running background process intended to listen or poll for events
// and then process those events.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {
	l := loggerWithContext(rnr.Log, "RunDaemon")

	for keepRunning() {

		m, err := rnr.initModeller(l)
		if err != nil {
			l.Warn("failed initialising modeller >%v<", err)
			return err
		}

		diRecs, err := m.GetDungeonInstanceRecs(nil, nil, false)
		if err != nil {
			l.Warn("failed getting dungeon instance records >%v<", err)
			return err
		}

		err = m.Rollback()
		if err != nil {
			l.Warn("failed model rollback >%v<", err)
			return err
		}

		type DungeonInstanceTurnRoutineResult struct {
			Error             error
			DungeonInstanceID string
			Turn              int
		}

		resChan := make(chan DungeonInstanceTurnRoutineResult, 1)

		for idx := range diRecs {
			go func(diRec *record.DungeonInstance) {

				m, err := rnr.initModeller(l)
				if err != nil {
					l.Warn("failed initialising modeller >%v<", err)
					resChan <- DungeonInstanceTurnRoutineResult{
						DungeonInstanceID: diRec.ID,
						Error:             err,
					}
					return
				}

				result, err := processDungeonInstanceTurn(l, m, diRec)
				if err != nil {
					l.Warn("failed processing dungeon instance ID >%s< turn >%v<", diRec.ID, err)
					rerr := m.Rollback()
					if rerr != nil {
						l.Warn("failed model rollback >%v<", rerr)
						err = fmt.Errorf("%v with %v", rerr, err)
					}

					resChan <- DungeonInstanceTurnRoutineResult{
						DungeonInstanceID: diRec.ID,
						Error:             err,
					}
					return
				}

				if result == nil {
					err := fmt.Errorf("result is nil, failed processing dungeon instance ID >%s< turn", diRec.ID)
					l.Warn(err.Error())
					rerr := m.Rollback()
					if rerr != nil {
						l.Warn("failed model rollback >%v<", rerr)
						err = fmt.Errorf("%v with %v", rerr, err)
					}
					resChan <- DungeonInstanceTurnRoutineResult{
						DungeonInstanceID: diRec.ID,
						Error:             err,
					}
					return
				}

				if result.Record == nil {
					err := fmt.Errorf("result record is nil, failed processing dungeon instance ID >%s< turn", diRec.ID)
					l.Warn(err.Error())
					rerr := m.Rollback()
					if rerr != nil {
						l.Warn("failed model rollback >%v<", rerr)
						err = fmt.Errorf("%v with %v", rerr, err)
					}
					resChan <- DungeonInstanceTurnRoutineResult{
						DungeonInstanceID: diRec.ID,
						Error:             err,
					}
					return
				}

				if result.Incremented {
					err = m.Commit()
					if err != nil {
						l.Warn("failed model commit >%v<", err)
						resChan <- DungeonInstanceTurnRoutineResult{
							DungeonInstanceID: diRec.ID,
							Error:             err,
						}
						return
					}
				} else {
					err = m.Rollback()
					if err != nil {
						l.Warn("failed model rollback >%v<", err)
						resChan <- DungeonInstanceTurnRoutineResult{
							DungeonInstanceID: diRec.ID,
							Error:             err,
						}
						return
					}
				}

				l.Info("Processed dungeon instance ID >%s< turn >%d< **", diRec.ID, result.Record.TurnCount)

				resChan <- DungeonInstanceTurnRoutineResult{
					DungeonInstanceID: diRec.ID,
					Turn:              result.Record.TurnCount,
				}
			}(diRecs[idx])
		}

		// TODO: Perhaps wait for responses from all dungeon instances?
		result := <-resChan

		if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

func processDungeonInstanceTurn(l logger.Logger, m *model.Model, rec *record.DungeonInstance) (*model.IncrementDungeonInstanceTurnResult, error) {
	l = loggerWithContext(l, "processDungeonInstanceTurn")

	var result *model.IncrementDungeonInstanceTurnResult
	var err error

	for result == nil || !result.Incremented {
		result, err = m.IncrementDungeonInstanceTurn(rec.ID)
		if err != nil {
			return nil, err
		}

		l.Warn("** Sleeping for milliseconds >%d<", result.WaitMilliseconds)

		if !result.Incremented && result.WaitMilliseconds > 0 {
			time.Sleep(time.Duration(result.WaitMilliseconds) * time.Millisecond)
		}
	}

	return result, nil
}

func keepRunning() bool {
	return true
}
