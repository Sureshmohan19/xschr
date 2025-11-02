// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."

package db

import (
	"database/sql"
	"time"

	"github.com/Sureshmohan19/xschr/internal/types"
)

// InsertJob adds a new job to the database and returns its ID
func (c *DBConn) InsertJob(command string, cpus int) (int64, error) {
	query := `
		INSERT INTO jobs (command, cpus, state, submitted_at)
		VALUES (?, ?, ?, ?)
	`

	now := time.Now().Unix()
	return c.InsertReturnID(query, command, cpus, types.StatePending, now)
}

// GetJobByID retrieves a single job by its ID
func (c *DBConn) GetJobByID(id int64) (*types.Job, error) {
	query := `
		SELECT id, command, cpus, state, pid, submitted_at, started_at, finished_at, exit_code
		FROM jobs
		WHERE id = ?
	`

	row := c.QueryRow(query, id)
	return scanJob(row)
}

// GetPendingJobs returns all jobs in pending state, ordered by submission time
func (c *DBConn) GetPendingJobs() ([]*types.Job, error) {
	query := `
		SELECT id, command, cpus, state, pid, submitted_at, started_at, finished_at, exit_code
		FROM jobs
		WHERE state = ?
		ORDER BY submitted_at ASC
	`

	rows, err := c.Query(query, types.StatePending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// GetRunningJobs returns all jobs currently running
func (c *DBConn) GetRunningJobs() ([]*types.Job, error) {
	query := `
		SELECT id, command, cpus, state, pid, submitted_at, started_at, finished_at, exit_code
		FROM jobs
		WHERE state = ?
		ORDER BY started_at ASC
	`

	rows, err := c.Query(query, types.StateRunning)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// GetAllJobs returns all jobs in the database
func (c *DBConn) GetAllJobs() ([]*types.Job, error) {
	query := `
		SELECT id, command, cpus, state, pid, submitted_at, started_at, finished_at, exit_code
		FROM jobs
		ORDER BY submitted_at DESC
	`

	rows, err := c.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanJobs(rows)
}

// UpdateJobState updates a job's state
func (c *DBConn) UpdateJobState(id int64, state types.JobState) error {
	query := "UPDATE jobs SET state = ? WHERE id = ?"
	return c.Exec(query, state, id)
}

// MarkJobRunning marks a job as running with PID and start time
func (c *DBConn) MarkJobRunning(id int64, pid int) error {
	query := `
		UPDATE jobs 
		SET state = ?, pid = ?, started_at = ?
		WHERE id = ?
	`

	now := time.Now().Unix()
	return c.Exec(query, types.StateRunning, pid, now, id)
}

// MarkJobFinished marks a job as done or failed with exit code and finish time
func (c *DBConn) MarkJobFinished(id int64, exitCode int) error {
	state := types.StateDone
	if exitCode != 0 {
		state = types.StateFailed
	}

	query := `
		UPDATE jobs 
		SET state = ?, exit_code = ?, finished_at = ?
		WHERE id = ?
	`

	now := time.Now().Unix()
	return c.Exec(query, state, exitCode, now, id)
}

// DeleteJob removes a job from the database
func (c *DBConn) DeleteJob(id int64) error {
	query := "DELETE FROM jobs WHERE id = ?"
	return c.Exec(query, id)
}

// scanJob scans a single row into a Job struct
func scanJob(row *sql.Row) (*types.Job, error) {
	var job types.Job
	var pid sql.NullInt64
	var startedAt sql.NullInt64
	var finishedAt sql.NullInt64
	var exitCode sql.NullInt64
	var submittedAtUnix int64

	err := row.Scan(
		&job.ID,
		&job.Command,
		&job.CPUs,
		&job.State,
		&pid,
		&submittedAtUnix,
		&startedAt,
		&finishedAt,
		&exitCode,
	)

	if err != nil {
		return nil, err
	}

	// Convert nullable fields
	if pid.Valid {
		job.PID = int(pid.Int64)
	}

	job.SubmittedAt = time.Unix(submittedAtUnix, 0)

	if startedAt.Valid {
		t := time.Unix(startedAt.Int64, 0)
		job.StartedAt = &t
	}

	if finishedAt.Valid {
		t := time.Unix(finishedAt.Int64, 0)
		job.FinishedAt = &t
	}

	if exitCode.Valid {
		ec := int(exitCode.Int64)
		job.ExitCode = &ec
	}

	return &job, nil
}

// scanJobs scans multiple rows into Job structs
func scanJobs(rows *sql.Rows) ([]*types.Job, error) {
	var jobs []*types.Job

	for rows.Next() {
		var job types.Job
		var pid sql.NullInt64
		var startedAt sql.NullInt64
		var finishedAt sql.NullInt64
		var exitCode sql.NullInt64
		var submittedAtUnix int64

		err := rows.Scan(
			&job.ID,
			&job.Command,
			&job.CPUs,
			&job.State,
			&pid,
			&submittedAtUnix,
			&startedAt,
			&finishedAt,
			&exitCode,
		)

		if err != nil {
			return nil, err
		}

		// Convert nullable fields
		if pid.Valid {
			job.PID = int(pid.Int64)
		}

		job.SubmittedAt = time.Unix(submittedAtUnix, 0)

		if startedAt.Valid {
			t := time.Unix(startedAt.Int64, 0)
			job.StartedAt = &t
		}

		if finishedAt.Valid {
			t := time.Unix(finishedAt.Int64, 0)
			job.FinishedAt = &t
		}

		if exitCode.Valid {
			ec := int(exitCode.Int64)
			job.ExitCode = &ec
		}

		jobs = append(jobs, &job)
	}

	return jobs, rows.Err()
}
