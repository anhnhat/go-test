package project

import "time"

type Health struct {
	Health       string `json:"Health"`
	HealthReason string `json:"HealthReason"`
}

type ProjectResponse struct {
	Name        string `json:"Name" binding:"required"`
	Description string `json:"Description"`
	StatusID    uint   `json:"StatusID,omitempty"`
	Priority    string `json:"Priority,omitempty"`
}

type CreateProjectRequest struct {
	Name           string `json:"Name" binding:"required"`
	Description    string `json:"Description"`
	StatusID       uint   `json:"StatusID,omitempty"`
	Priority       string `json:"Priority,omitempty"`
	Members        []uint `json:"Members,omitempty"`
	Health         Health
	ClientName     string    `json:"ClientName"`
	Budget         float32   `json:"Budget"`
	ActualReceived float32   `json:"ActualReceived"`
	StartDate      time.Time `json:"StartDate"`
	EndDate        time.Time `json:"EndDate"`
}
