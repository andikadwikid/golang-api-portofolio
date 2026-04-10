package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type JobHistory struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty"`
	CompanyName   string               `bson:"company_name"`
	Position      string               `bson:"position"`
	JobDescription string              `bson:"job_description"`
	Duration      string               `bson:"duration"`
	UserID        primitive.ObjectID   `bson:"user_id"`
	StartDate     time.Time            `bson:"start_date"`
	EndDate       time.Time            `bson:"end_date,omitempty"`
	IsCurrent     bool                 `bson:"is_current"`
	SkillIDs      []primitive.ObjectID `bson:"skill_ids,omitempty"`
	CreatedAt     time.Time            `bson:"created_at"`
	UpdatedAt     time.Time            `bson:"updated_at"`
}

type CreateJobHistoryInput struct {
	CompanyName   string    `json:"company_name" binding:"required"`
	Position      string    `json:"position" binding:"required"`
	JobDescription string   `json:"job_description" binding:"required"`
	Duration      string    `json:"duration" binding:"required"`
	StartDate     time.Time `json:"start_date" binding:"required"`
	EndDate       time.Time `json:"end_date,omitempty"`
	IsCurrent     bool      `json:"is_current"`
	SkillIDs      []string  `json:"skill_ids,omitempty"`
}

type UpdateJobHistoryInput struct {
	CompanyName   *string    `json:"company_name,omitempty"`
	Position      *string    `json:"position,omitempty"`
	JobDescription *string   `json:"job_description,omitempty"`
	Duration      *string    `json:"duration,omitempty"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	IsCurrent     *bool      `json:"is_current,omitempty"`
	SkillIDs      []string   `json:"skill_ids,omitempty"`
}
