package runner

import (
	"time"

	"gitlab.com/alienspaces/go-mud/backend/core/errors"
)

// General runtime constants
const (
	// Maximum errors to capture before actually crashing out
	errorLimit int = 10000
	// Number of errors to capture before raising the alarm
	alarmErrorCount int = 10
	// Number of milliSeconds to wait between cycles
	waitMilliseconds int = 3000

	// Valid args
	ArgCycleLimit string = "cycleLimit"
)

// RunDaemon - Starts the daemon process.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {

	rnr.Log.Debug("(template) RunDaemon")

	cycleLimit := 0
	if argRunTimes, ok := args[ArgCycleLimit]; ok {
		cycleLimit, ok = argRunTimes.(int)
		if !ok {
			cycleLimit = 0
		}
	}

	hasExceededErrorLimit := func(errCount int) bool {
		return errCount >= errorLimit
	}

	hasExceededCycleLimit := func(cycleCount int) bool {
		return cycleLimit != 0 && cycleCount >= cycleLimit
	}

	cycleCount := 0
	errs := &errors.Error{}
	for !hasExceededErrorLimit(errs.Count()) && !hasExceededCycleLimit(cycleCount) {

		// Wait for a bit
		if cycleCount != 0 {
			time.Sleep(time.Duration(waitMilliseconds) * time.Millisecond)
		}

		cycleCount++

		// Cycle log
		if cycleCount%10 == 0 || errs.Count() != 0 {
			rnr.Log.Info("(template) Cycle >%d< Errors >%d<", cycleCount, errs.Count())
		}

		// Raise the alarm
		if errs.Count() > 0 && errs.Count()%alarmErrorCount == 0 {
			rnr.Log.Error("(template) Ahwhooga Ahwhooga >%v<", errs.Error())
		}

		m, err := rnr.InitTx(rnr.Log)
		if err != nil {
			rnr.Log.Warn("(template) Failed initialising database transaction >%v<", err)
			errs.Add(err)
			continue
		}

		m.Commit()

		// Clear error stack if we get here as it means we've managed
		// to self correct any previous issues encountered.
		if errs.Count() > 0 {
			rnr.Log.Info("(template) We feel better now!")
			errs = &errors.Error{}
		}
	}

	if errs.Count() != 0 {
		return errs
	}

	return nil
}
