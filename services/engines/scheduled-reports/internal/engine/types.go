package engine

import (
	"context"
	"time"
)

type ReportFrequency string

const (
	FreqDaily   ReportFrequency = "daily"
	FreqWeekly  ReportFrequency = "weekly"
	FreqMonthly ReportFrequency = "monthly"
	FreqQuarterly ReportFrequency = "quarterly"
	FreqYearly  ReportFrequency = "yearly"
)

type ReportSection struct {
	Title string `json:"title"`
	Type  string `json:"type"`
}

type ReportSchedule struct {
	ID              string           `json:"id"`
	UserID          string           `json:"user_id"`
	Name            string           `json:"name"`
	Description     string           `json:"description,omitempty"`
	Frequency       ReportFrequency  `json:"frequency"`
	Sections        []ReportSection  `json:"sections"`
	Enabled         bool             `json:"enabled"`
	LastGeneratedAt *string          `json:"last_generated_at,omitempty"`
	CreatedAt       string           `json:"created_at"`
	UpdatedAt       string           `json:"updated_at"`
}

type GeneratedReport struct {
	ID          string                 `json:"id"`
	ScheduleID  string                 `json:"schedule_id"`
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name"`
	Frequency   ReportFrequency        `json:"frequency"`
	Sections    []ReportSection        `json:"sections"`
	Data        map[string]interface{} `json:"data"`
	GeneratedAt string                 `json:"generated_at"`
}

type DataProvider interface {
	FetchSectionData(ctx context.Context, userID string, sectionType string) (map[string]interface{}, error)
}

var _ = time.UTC
