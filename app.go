package main

import (
	"context"
	"log"

	"open-resume/internal/database"
	"open-resume/internal/model"
	"open-resume/internal/service"
)

// App struct
type App struct {
	ctx    context.Context
	svc    service.ResumeService
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		svc: service.NewResumeService(),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}

// shutdown is called when the app stops
func (a *App) shutdown(ctx context.Context) {
	// Close database connection
	if err := database.Close(); err != nil {
		log.Printf("Error closing database: %v", err)
	}
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return "Hello " + name + ", It's show time!"
}

// CreateResume creates a new resume with the given title
func (a *App) CreateResume(title string) (*model.Resume, error) {
	resume := &model.Resume{
		Title: title,
	}
	if err := a.svc.Create(context.Background(), resume); err != nil {
		return nil, err
	}
	return resume, nil
}

// GetResume retrieves a resume by its ID
func (a *App) GetResume(id string) (*model.Resume, error) {
	return a.svc.GetByID(context.Background(), id)
}

// ListResumes returns all non-deleted resumes
func (a *App) ListResumes() ([]*model.ResumeListItem, error) {
	return a.svc.List(context.Background())
}

// UpdateResume updates an existing resume with new JSON data
func (a *App) UpdateResume(id string, jsonData string) error {
	return a.svc.UpdateJSON(context.Background(), id, jsonData)
}

// DeleteResume soft-deletes a resume by its ID
func (a *App) DeleteResume(id string) error {
	return a.svc.Delete(context.Background(), id)
}

// UpdateResumeModule updates a specific module within a resume
func (a *App) UpdateResumeModule(id string, moduleType string, moduleData string) error {
	return a.svc.UpdateJSON(context.Background(), id, moduleData)
}
