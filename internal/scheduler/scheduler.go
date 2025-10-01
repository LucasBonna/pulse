package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"lucasbonna/pulse/db"
	"net/http"
	"sync"
	"time"
)

type Scheduler struct {
	db          *db.Queries
	httpClient  *http.Client
	runningJobs map[int64]bool
	mutex       sync.RWMutex
	ticker      *time.Ticker
	done        chan bool
}

func NewScheduler(database *db.Queries) *Scheduler {
	return &Scheduler{
		db: database,
		httpClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		runningJobs: make(map[int64]bool),
		done:        make(chan bool),
	}
}

func (s *Scheduler) Start(ctx context.Context) {
	log.Println("Starting job scheduler (checking every 1 second)...")

	s.ticker = time.NewTicker(1 * time.Second)

	go s.run(ctx)
}

func (s *Scheduler) Stop() {
	log.Println("Stopping job scheduler...")
	if s.ticker != nil {
		s.ticker.Stop()
	}
	s.done <- true
}

func (s *Scheduler) run(ctx context.Context) {
	for {
		select {
		case <-s.done:
			log.Println("Scheduler stopped")
			return
		case <-s.ticker.C:
			s.checkAndRunJobs(ctx)
		}
	}
}

func (s *Scheduler) checkAndRunJobs(ctx context.Context) {
	log.Println("checking for jobs to run")
	jobs, err := s.db.GetDueJobs(ctx)
	if err != nil {
		log.Printf("error getting due jobs: %v", err)
		return
	}

	for _, job := range jobs {
		log.Printf("job %d (%s) found, checking before running...", job.ID, job.Name)
		if s.isJobRunning(job.ID) {
			log.Printf("Job %d (%s) is already running, skipping", job.ID, job.Name)
			continue
		}

		s.markJobAsRunning(job.ID, true)
		go s.executeJob(ctx, job)
	}
}

func (s *Scheduler) isJobRunning(jobID int64) bool {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.runningJobs[jobID]
}

func (s *Scheduler) markJobAsRunning(jobID int64, running bool) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if running {
		s.runningJobs[jobID] = true
	} else {
		delete(s.runningJobs, jobID)
	}
}

func (s *Scheduler) executeJob(ctx context.Context, job db.Job) {
	defer s.markJobAsRunning(job.ID, false)

	log.Printf("executing job %d: %s %s", job.ID, job.Method, job.Url)

	startTime := time.Now()

	// jobRun, err := s.db.CreateJobRun(ctx, db.CreateJobRunParams{
	// 	JobID:     job.ID,
	// 	Status:    sql.NullString{String: "running", Valid: true},
	// 	StartedAt: sql.NullTime{Time: startTime, Valid: true},
	// })
	// if err != nil {
	// 	log.Printf("Failed to create job run record: %v", err)
	// 	return
	// }

	_, err := s.makeHTTPRequest(job)
	finishTime := time.Now()

	status := "success"
	if err != nil {
		status = "failed"
		log.Printf("error executing job %d: %v", job.ID, err)
	}

	// s.db.UpdateJobRun(ctx, db.UpdateJobRunParams{
	// 	ID:           jobRun.ID,
	// 	Status:       sql.NullString{String: status, Valid: true},
	// 	ResponseCode: sql.NullInt64{Int64: int64(statusCode), Valid: true},
	// 	ResponseBody: sql.NullString{String: "", Valid: false},
	// 	FinishedAt:   sql.NullTime{Time: finishTime, Valid: true},
	// })

	currentTime := time.Now()

	duration := time.Duration(job.IntervalSeconds) * time.Second

	nextRun := currentTime.Add(duration).UTC()
	s.db.UpdateJobNextRun(ctx, db.UpdateJobNextRunParams{
		ID:        job.ID,
		NextRunAt: sql.NullTime{Time: nextRun, Valid: true},
	})

	log.Printf("job %d completed with status %s in %v", job.ID, status, finishTime.Sub(startTime))
}

func (s *Scheduler) makeHTTPRequest(job db.Job) (int, error) {
	req, err := http.NewRequest(job.Method.(string), job.Url, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}
