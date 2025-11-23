package models

import (
	"time"
)

// UserJson defines model for UserJson.
type UserJson struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	TeamName string `json:"team_name"`
}

type UserDB struct {
	UserId              string `gorm:"primaryKey"`
	IsActive            bool   `gorm:"index"`
	Username            string
	CurrTeamName        string          `gorm:"index"`
	Team                TeamDB          `gorm:"foreignKey:CurrTeamName"`
	PullRequestAssignee []PullRequestDB `gorm:"many2many:user_assignees"`
}

// PullRequest defines model for PullRequest.
type PullRequestJson struct {
	// AssignedReviewers user_id назначенных ревьюверов (0..2)
	AssignedReviewers []string              `json:"assigned_reviewers"`
	AuthorId          string                `json:"author_id"`
	CreatedAt         time.Time             `json:"createdAt"`
	MergedAt          time.Time             `json:"mergedAt"`
	PullRequestId     string                `json:"pull_request_id"`
	PullRequestName   string                `json:"pull_request_name"`
	Status            PullRequestEnumStatus `json:"status"`
}

type PullRequestDB struct {
	PullRequestId     string   `gorm:"primaryKey"`
	AssignedReviewers []UserDB `gorm:"many2many:user_assignees"`
	Author            UserDB   `gorm:"foreignKey:AuthorId"`
	AuthorId          string   `gorm:"index"`
	CreatedAt         time.Time
	MergedAt          time.Time
	PullRequestName   string
	Status            PullRequestEnumStatus
}

// TeamDB defines model for TeamDB.
type TeamDB struct {
	TeamName string `gorm:"primaryKey"`
}

type TeamJson struct {
	Members  []TeamMember `json:"members"`
	TeamName string       `json:"team_name"`
}

type NestedError struct {
	Code    ErrorResponseErrorCode `json:"code"`
	Message string                 `json:"message"`
}

// ErrorResponse defines model for ErrorResponse.
type ErrorResponse struct {
	Error NestedError `json:"error"`
}

// PullRequestShort defines model for PullRequestShort.
type PullRequestShort struct {
	AuthorId        string                `json:"author_id"`
	PullRequestId   string                `json:"pull_request_id"`
	PullRequestName string                `json:"pull_request_name"`
	Status          PullRequestEnumStatus `json:"status"`
}

// TeamMember defines model for TeamMember.
type TeamMember struct {
	IsActive bool   `json:"is_active"`
	UserId   string `json:"user_id"`
	Username string `json:"username"`
}

// TeamNameQuery defines model for TeamNameQuery.
type TeamNameQuery = string

// UserIdQuery defines model for UserIdQuery.
type UserIdQuery = string
