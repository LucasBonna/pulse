CREATE TABLE IF NOT EXISTS jobs (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  url TEXT NOT NULL,
  method TEXT NOT NULL DEFAULT "GET",
  headers TEXT,
  interval_seconds INTEGER NOT NULL,
  next_run_at DATETIME,
  active boolean DEFAULT 1
);

CREATE TABLE IF NOT EXISTS job_runs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    job_id INTEGER NOT NULL,
    status TEXT,
    response_code INTEGER,
    response_body TEXT,
    started_at DATETIME,
    finished_at DATETIME,
    FOREIGN KEY (job_id) REFERENCES jobs(id)
);
