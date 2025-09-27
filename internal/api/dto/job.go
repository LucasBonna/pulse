package dto

import "time"

type CreateJobRequest struct {
	Name            string `json:"name" validate:"required,min=1,max=100"`
	URL             string `json:"url" validate:"required,url"`
	Method          string `json:"method" validate:"required,oneof=GET POST PUT PATCH DELETE"`
	Headers         string `json:"headers,omitempty" validate:"max=1000"`
	IntervalSeconds int64  `json:"interval_seconds" validate:"required,min=1,max=86400"`
	Active          bool   `json:"active,omitempty"`
}

type CreateJobResponse struct {
	Id              int64      `json:"id"`
	Name            string     `json:"name"`
	Url             string     `json:"url"`
	Method          string     `json:"method"`
	Headers         *string    `json:"headers"`
	IntervalSeconds int64      `json:"interval_seconds"`
	NextRunAt       *time.Time `json:"next_runt_at"`
	Active          *bool      `json:"active"`
}

type UpdateJobRequest struct {
	Name            string `json:"name,omitempty" validate:"omitempty,min=1,max=100"`
	URL             string `json:"url,omitempty" validate:"omitempty,url"`
	Method          string `json:"method,omitempty" validate:"omitempty,oneof=GET POST PUT PATCH DELETE"`
	Headers         string `json:"headers,omitempty" validate:"max=1000"`
	IntervalSeconds *int64 `json:"interval_seconds,omitempty" validate:"omitempty,min=1,max=86400"`
	Active          *bool  `json:"active,omitempty"`
}

type JobIDRequest struct {
	ID int64 `json:"id" validate:"required,min=1"`
}
