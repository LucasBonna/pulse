-- name: GetJobByID :one
SELECT * FROM jobs WHERE id = ? LIMIT 1;

-- name: UpdateJob :one
UPDATE jobs
SET name = ?, url = ?, method = ?, headers = ?, interval_seconds = ?, active = ?
WHERE id = ?
RETURNING *;

-- name: GetAllJobs :many
SELECT * FROM jobs
ORDER BY id;

-- name: GetDueJobs :many
SELECT * FROM jobs
WHERE active = 1
  AND next_run_at IS NOT NULL
  AND next_run_at <= datetime('now')
ORDER BY next_run_at ASC;

-- name: CreateJob :one
INSERT INTO jobs (
  name, url, method, headers, interval_seconds, next_run_at, active
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;

-- name: UpdateJobNextRun :exec
UPDATE jobs
SET next_run_at = ?
WHERE id = ?;

-- name: CreateJobRun :one
INSERT INTO job_runs (job_id, status, started_at)
VALUES (?, ?, ?)
RETURNING *;

-- name: UpdateJobRun :exec
UPDATE job_runs
SET status = ?, response_code = ?, response_body = ?, finished_at = ?
WHERE id = ?;

-- name: DeleteJob :exec
DELETE FROM jobs
WHERE id = ?;
