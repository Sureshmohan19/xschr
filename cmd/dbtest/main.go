// XSchr - Experimental Scheduler in Go
// "Make it work, make it right, make it fast — in that order."
//
// Comprehensive database test suite
// Run with: go run cmd/dbtest/main.go

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Sureshmohan19/xschr/internal/db"
	"github.com/Sureshmohan19/xschr/internal/types"
)

func main() {
	dbPath := "./xschr_test.db"
	os.Remove(dbPath)

	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║        XSchr Database Layer - Final Test Suite        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Test 1: Database Initialization
	fmt.Println("TEST 1: Database Initialization")
	fmt.Println("--------------------------------")
	conn, err := db.OpenDB(dbPath)
	if err != nil {
		log.Fatal("❌ Failed to open database:", err)
	}
	defer conn.Close()
	fmt.Println("✓ Database opened successfully")
	fmt.Println("✓ schema_versions table created")
	fmt.Println("✓ jobs table created")
	fmt.Println()

	// Test 2: Job Insertion
	fmt.Println("TEST 2: Job Insertion")
	fmt.Println("---------------------")
	testJobs := []struct {
		command string
		cpus    int
	}{
		{"python train_model.py --epochs 100", 8},
		{"bash preprocess.sh dataset.csv", 4},
		{"./simulation --particles 1000000", 16},
		{"python analyze.py results.json", 2},
		{"make -j8 compile_project", 8},
	}

	var jobIDs []int64
	for i, tj := range testJobs {
		id, err := conn.InsertJob(tj.command, tj.cpus)
		if err != nil {
			log.Fatalf("❌ Failed to insert job %d: %v", i+1, err)
		}
		jobIDs = append(jobIDs, id)
		fmt.Printf("✓ Job %d inserted (ID: %d, CPUs: %d)\n", i+1, id, tj.cpus)
	}
	fmt.Println()

	// Test 3: Query Pending Jobs
	fmt.Println("TEST 3: Query Pending Jobs")
	fmt.Println("--------------------------")
	pending, err := conn.GetPendingJobs()
	if err != nil {
		log.Fatal("❌ Failed to get pending jobs:", err)
	}
	fmt.Printf("✓ Found %d pending jobs\n", len(pending))
	if len(pending) != 5 {
		log.Fatalf("❌ Expected 5 pending jobs, got %d", len(pending))
	}
	for _, job := range pending {
		if job.State != types.StatePending {
			log.Fatalf("❌ Job %d has wrong state: %s", job.ID, job.State)
		}
	}
	fmt.Println("✓ All jobs have correct pending state")
	fmt.Println()

	// Test 4: Job State Transitions
	fmt.Println("TEST 4: Job State Transitions")
	fmt.Println("------------------------------")

	// Start job 1
	fmt.Printf("Starting Job %d...\n", jobIDs[0])
	err = conn.MarkJobRunning(jobIDs[0], 10001)
	if err != nil {
		log.Fatal("❌ Failed to mark job running:", err)
	}
	job1, err := conn.GetJobByID(jobIDs[0])
	if err != nil {
		log.Fatal("❌ Failed to get job:", err)
	}
	if job1.State != types.StateRunning || job1.PID != 10001 {
		log.Fatalf("❌ Job state incorrect: %s, PID: %d", job1.State, job1.PID)
	}
	if job1.StartedAt == nil {
		log.Fatal("❌ Job started_at is nil")
	}
	fmt.Printf("✓ Job %d: pending → running (PID: %d)\n", jobIDs[0], job1.PID)

	// Start job 2
	fmt.Printf("Starting Job %d...\n", jobIDs[1])
	err = conn.MarkJobRunning(jobIDs[1], 10002)
	if err != nil {
		log.Fatal("❌ Failed to mark job running:", err)
	}
	fmt.Printf("✓ Job %d: pending → running (PID: 10002)\n", jobIDs[1])

	// Finish job 1 successfully
	time.Sleep(100 * time.Millisecond)
	fmt.Printf("Finishing Job %d (success)...\n", jobIDs[0])
	err = conn.MarkJobFinished(jobIDs[0], 0)
	if err != nil {
		log.Fatal("❌ Failed to mark job finished:", err)
	}
	job1, err = conn.GetJobByID(jobIDs[0])
	if err != nil {
		log.Fatal("❌ Failed to get job:", err)
	}
	if job1.State != types.StateDone || *job1.ExitCode != 0 {
		log.Fatalf("❌ Job state incorrect: %s, exit: %d", job1.State, *job1.ExitCode)
	}
	if job1.FinishedAt == nil {
		log.Fatal("❌ Job finished_at is nil")
	}
	fmt.Printf("✓ Job %d: running → done (exit: 0)\n", jobIDs[0])

	// Fail job 2
	fmt.Printf("Finishing Job %d (failure)...\n", jobIDs[1])
	err = conn.MarkJobFinished(jobIDs[1], 1)
	if err != nil {
		log.Fatal("❌ Failed to mark job finished:", err)
	}
	job2, err := conn.GetJobByID(jobIDs[1])
	if err != nil {
		log.Fatal("❌ Failed to get job:", err)
	}
	if job2.State != types.StateFailed || *job2.ExitCode != 1 {
		log.Fatalf("❌ Job state incorrect: %s, exit: %d", job2.State, *job2.ExitCode)
	}
	fmt.Printf("✓ Job %d: running → failed (exit: 1)\n", jobIDs[1])
	fmt.Println()

	// Test 5: Query by State
	fmt.Println("TEST 5: Query by State")
	fmt.Println("----------------------")
	pending, err = conn.GetPendingJobs()
	if err != nil {
		log.Fatal("❌ Failed to get pending jobs:", err)
	}
	fmt.Printf("✓ Pending jobs: %d (expected 3)\n", len(pending))
	if len(pending) != 3 {
		log.Fatalf("❌ Expected 3 pending jobs, got %d", len(pending))
	}

	running, err := conn.GetRunningJobs()
	if err != nil {
		log.Fatal("❌ Failed to get running jobs:", err)
	}
	fmt.Printf("✓ Running jobs: %d (expected 0)\n", len(running))
	if len(running) != 0 {
		log.Fatalf("❌ Expected 0 running jobs, got %d", len(running))
	}

	allJobs, err := conn.GetAllJobs()
	if err != nil {
		log.Fatal("❌ Failed to get all jobs:", err)
	}
	fmt.Printf("✓ Total jobs: %d (expected 5)\n", len(allJobs))
	if len(allJobs) != 5 {
		log.Fatalf("❌ Expected 5 total jobs, got %d", len(allJobs))
	}
	fmt.Println()

	// Test 6: Get Job by ID
	fmt.Println("TEST 6: Get Job by ID")
	fmt.Println("---------------------")
	for _, id := range jobIDs {
		job, err := conn.GetJobByID(id)
		if err != nil {
			log.Fatalf("❌ Failed to get job %d: %v", id, err)
		}
		if job.ID != id {
			log.Fatalf("❌ Job ID mismatch: expected %d, got %d", id, job.ID)
		}
		fmt.Printf("✓ Job %d retrieved correctly\n", id)
	}
	fmt.Println()

	// Test 7: Timestamp Validation
	fmt.Println("TEST 7: Timestamp Validation")
	fmt.Println("----------------------------")
	job1, _ = conn.GetJobByID(jobIDs[0])
	if job1.SubmittedAt.IsZero() {
		log.Fatal("❌ submitted_at is zero")
	}
	fmt.Printf("✓ submitted_at: %s\n", job1.SubmittedAt.Format(time.RFC3339))

	if job1.StartedAt == nil || job1.StartedAt.IsZero() {
		log.Fatal("❌ started_at is nil or zero")
	}
	fmt.Printf("✓ started_at:   %s\n", job1.StartedAt.Format(time.RFC3339))

	if job1.FinishedAt == nil || job1.FinishedAt.IsZero() {
		log.Fatal("❌ finished_at is nil or zero")
	}
	fmt.Printf("✓ finished_at:  %s\n", job1.FinishedAt.Format(time.RFC3339))

	duration := job1.FinishedAt.Sub(*job1.StartedAt)
	fmt.Printf("✓ Job duration: %s\n", duration)
	fmt.Println()

	// Test 8: Delete Job
	fmt.Println("TEST 8: Delete Job")
	fmt.Println("------------------")
	deleteID := jobIDs[4]
	err = conn.DeleteJob(deleteID)
	if err != nil {
		log.Fatal("❌ Failed to delete job:", err)
	}
	fmt.Printf("✓ Job %d deleted\n", deleteID)

	_, err = conn.GetJobByID(deleteID)
	if err == nil {
		log.Fatal("❌ Deleted job still exists")
	}
	fmt.Printf("✓ Job %d confirmed deleted\n", deleteID)

	allJobs, _ = conn.GetAllJobs()
	if len(allJobs) != 4 {
		log.Fatalf("❌ Expected 4 jobs after delete, got %d", len(allJobs))
	}
	fmt.Printf("✓ Total jobs now: %d\n", len(allJobs))
	fmt.Println()

	// Test 9: Concurrent Access
	fmt.Println("TEST 9: Concurrent Access (Thread Safety)")
	fmt.Println("------------------------------------------")
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func(n int) {
			_, err := conn.InsertJob(fmt.Sprintf("concurrent_job_%d", n), 1)
			if err != nil {
				log.Printf("❌ Concurrent insert %d failed: %v", n, err)
			}
			done <- true
		}(i)
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	allJobs, _ = conn.GetAllJobs()
	if len(allJobs) != 14 {
		log.Fatalf("❌ Expected 14 jobs after concurrent inserts, got %d", len(allJobs))
	}
	fmt.Println("✓ 10 concurrent inserts completed")
	fmt.Printf("✓ Total jobs now: %d\n", len(allJobs))
	fmt.Println()

	// Final Summary
	fmt.Println("╔════════════════════════════════════════════════════════╗")
	fmt.Println("║                    TEST SUMMARY                        ║")
	fmt.Println("╚════════════════════════════════════════════════════════╝")
	fmt.Println()

	states := map[types.JobState]int{}
	for _, job := range allJobs {
		states[job.State]++
	}

	fmt.Println("Database Statistics:")
	fmt.Println("-------------------")
	fmt.Printf("Total Jobs:    %d\n", len(allJobs))
	fmt.Printf("Pending:       %d\n", states[types.StatePending])
	fmt.Printf("Running:       %d\n", states[types.StateRunning])
	fmt.Printf("Done:          %d\n", states[types.StateDone])
	fmt.Printf("Failed:        %d\n", states[types.StateFailed])
	fmt.Println()

	fmt.Println("✅ ALL TESTS PASSED!")
	fmt.Println()
	fmt.Println("Database layer is complete and ready for:")
	fmt.Println("  • CLI commands (submit, status, cancel)")
	fmt.Println("  • Scheduler daemon")
	fmt.Println("  • Job execution engine")
	fmt.Println()
	fmt.Println("Database file: " + dbPath)
}
