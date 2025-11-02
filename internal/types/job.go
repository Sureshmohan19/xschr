// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."

package types

import "time"

// JobState represents the current state of a job
type JobState string

const (
	StatePending JobState = "pending"
	StateRunning JobState = "running"
	StateDone    JobState = "done"
	StateFailed  JobState = "failed"
)

// Job represents a compute job to be scheduled and executed
type Job struct {
	ID          int64
	Command     string
	CPUs        int
	State       JobState
	PID         int // Process ID when running
	SubmittedAt time.Time
	StartedAt   *time.Time // nil if not started
	FinishedAt  *time.Time // nil if not finished
	ExitCode    *int       // nil if not finished
}
