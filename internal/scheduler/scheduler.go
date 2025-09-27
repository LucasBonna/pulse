package scheduler

import (
	"context"
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
		db:          database,
		httpClient:  &http.Client{},
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
	jobs, err := s.getDueJobs(ctx)
	if err != nil {
		log.Printf("error getting due jobs: %v", err)
		return
	}

	for _, job := range jobs {
		if s.isJobRunning(job.ID) {
			log.Printf("Job %d (%s) is already running, skipping", job.ID, job.Name)
			continue
		}

		s.markJobAsRunning(job.ID, true)
		go s.executeJob(ctx, job)
	}
}

func (s *Scheduler) isJobRunning(jobID int64) bool {
	s.mutex.Lock()
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

	statusCode, err := s.makeHTTPRequest(job)
	if err != nil {
		log.Printf("error executing job %d", err)
	}

	nextRun := time.Now().Add(time.Duration(job.IntervalSeconds) * time.Second)
	s.updateJobNextRun(ctx, job.ID, nextRun)

	log.Printf("job %d completed with statusCode %v", job.ID, statusCode)
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

	return resp.StatusCode, nil
}
