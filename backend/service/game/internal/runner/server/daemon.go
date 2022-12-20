package runner

import (
	"time"
)

// RunDaemon - Starts the daemon process. Override to implement a custom daemon run function.
// The daemon process is a long running background process intended to listen or poll for events
// and then process those events.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {

	rnr.Log.Debug("** RunDaemon **")

	turn := 0
	for keepRunning() {

		// TODO: Cycle through dungeon instances incrementing turns and processing all monster actions in a goroutine
		waitForDungeonInstanceTurn()

		rnr.Log.Info("** RUN >%d< **", turn)

		turn++
	}

	return nil
}

func waitForDungeonInstanceTurn() {
	// TODO: Update dungeon instance turn counter
	time.Sleep(1500 * time.Millisecond)
}

func keepRunning() bool {
	return true
}
