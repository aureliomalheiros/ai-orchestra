package agents

import (
	"context"
	"time"
)

type Planner interface {
	CreatePlan(ctx context.Context, feature, context string) (*Plan, error)
}

type Executor interface {
	ExecuteTask(ctx context.Context, task Task, planContext string) (ExecutionResult, error)
}

type Plan struct {
	Tasks           []Task    `json:"tasks"`
	Context         string    `json:"context"`
	EstimatedEffort string    `json:"estimated_effort"`
	CreatedAt       time.Time `json:"created_at"`
}

type Task struct {
	ID                 string   `json:"id"`
	Description        string   `json:"description"`
	Files              []string `json:"files"`
	AcceptanceCriteria string   `json:"acceptance_criteria"`
	Dependencies       []string `json:"dependencies"`
	Priority           string   `json:"priority"`
	Status             string   `json:"status"`
}

type ExecutionResult struct {
	Success      bool          `json:"success"`
	Output       string        `json:"output"`
	FilesChanged []string      `json:"files_changed"`
	Duration     time.Duration `json:"duration"`
	Error        string        `json:"error,omitempty"`
}

type ReviewResult struct {
	Approved     bool     `json:"approved"`
	Issues       []string `json:"issues"`
	Suggestions  []string `json:"suggestions"`
	QualityScore float64  `json:"quality_score"`
	Summary      string   `json:"summary"`
}
