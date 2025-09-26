-- name: GetAllJobs :many
SELECT * FROM jobs
ORDER BY id;

-- name: CreateJob :one
INSERT INTO jobs (
  name, url, method, headers, interval_seconds, next_run_at, active
) VALUES (
  ?, ?, ?, ?, ?, ?, ?
)
RETURNING *;
