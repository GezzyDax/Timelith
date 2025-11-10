package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/logger"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/GezzyDax/timelith/go-backend/internal/telegram"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type MessageJob struct {
	ScheduleID uuid.UUID
	Account    *models.Account
	Template   *models.Template
	Channel    *models.Channel
	Message    string
	Delay      time.Duration
	Retries    int
}

type Dispatcher struct {
	db             *database.DB
	sessionManager *telegram.SessionManager
	queue          chan *MessageJob
	workers        int
	stopCh         chan struct{}
}

func NewDispatcher(db *database.DB, sessionManager *telegram.SessionManager) *Dispatcher {
	return &Dispatcher{
		db:             db,
		sessionManager: sessionManager,
		queue:          make(chan *MessageJob, 100),
		workers:        5,
		stopCh:         make(chan struct{}),
	}
}

func (d *Dispatcher) Run(ctx context.Context) {
	logger.Log.Info("Starting dispatcher",
		zap.Int("workers", d.workers))

	// Start workers
	for i := 0; i < d.workers; i++ {
		go d.worker(ctx, i)
	}

	<-d.stopCh
	close(d.queue)

	logger.Log.Info("Dispatcher stopped")
}

func (d *Dispatcher) worker(ctx context.Context, workerID int) {
	logger.Log.Info("Dispatcher worker started",
		zap.Int("worker_id", workerID))

	for job := range d.queue {
		d.processJob(ctx, job, workerID)
	}

	logger.Log.Info("Dispatcher worker stopped",
		zap.Int("worker_id", workerID))
}

func (d *Dispatcher) processJob(ctx context.Context, job *MessageJob, workerID int) {
	// Apply delay if specified
	if job.Delay > 0 {
		logger.Log.Info("Delaying message",
			zap.Duration("delay", job.Delay),
			zap.String("channel", job.Channel.Name))
		time.Sleep(job.Delay)
	}

	logger.Log.Info("Processing message job",
		zap.Int("worker_id", workerID),
		zap.String("schedule_id", job.ScheduleID.String()),
		zap.String("account", job.Account.Phone),
		zap.String("channel", job.Channel.Name))

	// Load Telegram session if not already loaded
	if err := d.sessionManager.LoadSession(ctx, job.Account); err != nil {
		logger.Log.Error("Failed to load Telegram session",
			zap.String("account", job.Account.Phone),
			zap.Error(err))
		d.logJobResult(job.ScheduleID, "failed", "", fmt.Sprintf("Failed to load session: %v", err))
		return
	}

	// Send message (text or media based on template)
	var err error
	if job.Template.MediaType.Valid && job.Template.MediaType.String != "" {
		err = d.sessionManager.SendMediaMessage(ctx, job.Account.Phone, job.Channel.ChatID, job.Template)
	} else if job.Template.CopyFromChatID.Valid && job.Template.CopyFromMessageID.Valid {
		err = d.sessionManager.ForwardMessage(ctx, job.Account.Phone, job.Channel.ChatID,
			job.Template.CopyFromChatID.String, int(job.Template.CopyFromMessageID.Int64))
	} else {
		err = d.sessionManager.SendMessage(ctx, job.Account.Phone, job.Channel.ChatID, job.Message)
	}

	if err != nil {
		logger.Log.Error("Failed to send message",
			zap.String("account", job.Account.Phone),
			zap.String("channel", job.Channel.ChatID),
			zap.Error(err))

		// Retry logic
		if job.Retries < 3 {
			job.Retries++
			logger.Log.Info("Retrying message job",
				zap.Int("retry", job.Retries),
				zap.String("schedule_id", job.ScheduleID.String()))

			time.Sleep(time.Second * time.Duration(job.Retries*2)) // Exponential backoff
			d.Enqueue(job)
			return
		}

		d.logJobResult(job.ScheduleID, "failed", "", fmt.Sprintf("Failed after %d retries: %v", job.Retries, err))
		return
	}

	// Success - update account statistics
	if err := d.db.IncrementAccountMessageCount(job.Account.ID); err != nil {
		logger.Log.Error("Failed to update account stats",
			zap.String("account_id", job.Account.ID.String()),
			zap.Error(err))
	}

	logger.Log.Info("Message sent successfully",
		zap.String("account", job.Account.Phone),
		zap.String("channel", job.Channel.ChatID))

	d.logJobResult(job.ScheduleID, "success", "Message sent successfully", "")
}

func (d *Dispatcher) Enqueue(job *MessageJob) {
	select {
	case d.queue <- job:
		logger.Log.Debug("Job enqueued",
			zap.String("schedule_id", job.ScheduleID.String()))
	default:
		logger.Log.Warn("Dispatcher queue is full, dropping job",
			zap.String("schedule_id", job.ScheduleID.String()))
	}
}

func (d *Dispatcher) logJobResult(scheduleID uuid.UUID, status, message, errorMsg string) {
	log := &models.JobLog{
		ScheduleID: scheduleID,
		Status:     status,
		ExecutedAt: time.Now(),
	}

	if message != "" {
		log.Message = models.NewNullString(message)
	}
	if errorMsg != "" {
		log.Error = models.NewNullString(errorMsg)
	}

	if err := d.db.CreateJobLog(log); err != nil {
		logger.Log.Error("Failed to create job log",
			zap.String("schedule_id", scheduleID.String()),
			zap.Error(err))
	}
}

func (d *Dispatcher) Stop() {
	close(d.stopCh)
}
