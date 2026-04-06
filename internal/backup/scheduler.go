package backup

import (
	"context"
	"log"
	"sync"
	"time"

	"Darvin-Resume/internal/crypto"
)

// Scheduler manages automatic backup scheduling.
type Scheduler struct {
	mu       sync.Mutex
	interval time.Duration
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	running  bool
	started  bool
}

// NewScheduler creates a new backup scheduler.
func NewScheduler(interval time.Duration) *Scheduler {
	if interval < 5*time.Minute {
		interval = 5 * time.Minute
	}
	return &Scheduler{
		interval: interval,
	}
}

// Start starts the automatic backup scheduler.
// It runs in a goroutine and performs backups at the configured interval.
func (s *Scheduler) Start(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running || s.started {
		return
	}
	s.started = true
	s.running = true

	schedCtx, cancel := context.WithCancel(ctx)
	s.cancel = cancel

	s.wg.Add(1)
	go s.run(schedCtx)

	log.Printf("[BackupScheduler] started with interval: %v", s.interval)
}

// run is the main scheduler loop.
func (s *Scheduler) run(ctx context.Context) {
	defer s.wg.Done()

	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	// Perform initial backup after a short delay (so app startup is not blocked)
	select {
	case <-ctx.Done():
		return
	case <-time.After(10 * time.Second):
		s.performBackup()
	}

	for {
		select {
		case <-ctx.Done():
			s.mu.Lock()
			s.running = false
			s.mu.Unlock()
			log.Println("[BackupScheduler] stopped")
			return
		case <-ticker.C:
			s.performBackup()
		}
	}
}

// performBackup executes a single automatic backup using device key encryption.
func (s *Scheduler) performBackup() {
	deviceKey, err := crypto.DeviceKey()
	if err != nil {
		log.Printf("[BackupScheduler] get device key failed: %v", err)
		return
	}

	path, err := CreateBackup(deviceKey)
	if err != nil {
		log.Printf("[BackupScheduler] backup failed: %v", err)
		return
	}

	log.Printf("[BackupScheduler] auto backup created: %s", path)
}

// Stop stops the scheduler gracefully.
// It waits for any in-progress backup to complete before returning.
func (s *Scheduler) Stop() {
	s.mu.Lock()
	if !s.running {
		s.mu.Unlock()
		return
	}

	if s.cancel != nil {
		s.cancel()
	}
	s.mu.Unlock()

	// Wait for the goroutine to finish
	s.wg.Wait()

	s.mu.Lock()
	s.running = false
	s.mu.Unlock()

	log.Println("[BackupScheduler] stopped gracefully")
}

// IsRunning returns whether the scheduler is currently running.
func (s *Scheduler) IsRunning() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.running
}

// SetInterval updates the backup interval.
// Note: The scheduler must be restarted for the new interval to take effect.
func (s *Scheduler) SetInterval(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if interval < 5*time.Minute {
		interval = 5 * time.Minute
	}
	s.interval = interval
}

// GetInterval returns the current interval.
func (s *Scheduler) GetInterval() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.interval
}
