package server

// RunDaemon - Starts the daemon process. Override to implement a custom daemon run function.
// The daemon process is a long running background process intended to listen or poll for events
// and then process those events.
func (rnr *Runner) RunDaemon(args map[string]interface{}) error {

	rnr.Log.Debug("** RunDaemon **")

	return nil
}
