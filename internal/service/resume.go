package service

import (
	"context"
	"errors"

	"open-resume/internal/model"
)

// ResumeService defines the interface for resume operations
type ResumeService interface {
	Create(ctx context.Context, resume *model.Resume) error
	GetByID(ctx context.Context, id string) (*model.Resume, error)
	List(ctx context.Context) ([]*model.ResumeListItem, error)
	Update(ctx context.Context, resume *model.Resume) error
	Delete(ctx context.Context, id string) error
	UpdateJSON(ctx context.Context, id string, jsonData string) error
}

// ErrResumeNotFound is returned when a resume is not found
var ErrResumeNotFound = errors.New("resume not found")

// ErrInvalidData is returned when data is invalid
var ErrInvalidData = errors.New("invalid data")
