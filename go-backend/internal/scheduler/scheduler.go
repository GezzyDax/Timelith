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
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type Scheduler struct {
	cron           *cron.Cron
	db             *database.DB
	sessionManager *telegram.SessionManager
	jobs           map[uuid.UUID]cron.EntryID
	dispatcher     *Dispatcher
}

func NewScheduler(db *database.DB, sessionManager *telegram.SessionManager) *Scheduler {
	return &Scheduler{
		cron:           cron.New(cron.WithSeconds()),
		db:             db,
		sessionManager: sessionManager,
		jobs:           make(map[uuid.UUID]cron.EntryID),
		dispatcher:     NewDispatcher(db, sessionManager),
	}
}

func (s *Scheduler) Start(ctx context.Context) error {
	logger.Log.Info("Starting scheduler")

	// Load all active schedules
	schedules, err := s.db.ListActiveSchedules()
	if err != nil {
		return fmt.Errorf("failed to load schedules: %w", err)
	}

	// Add each schedule to cron
	for _, schedule := range schedules {
		// Create a copy to avoid implicit memory aliasing
		sched := schedule
		if err := s.AddSchedule(&sched); err != nil {
			logger.Log.Error("Failed to add schedule",
				zap.String("schedule_id", sched.ID.String()),
				zap.Error(err))
			continue
		}
	}

	// Start cron
	s.cron.Start()

	logger.Log.Info("Scheduler started",
		zap.Int("schedules_loaded", len(schedules)))

	// Run dispatcher in background
	go s.dispatcher.Run(ctx)

	return nil
}

func (s *Scheduler) AddSchedule(schedule *models.Schedule) error {
	// Parse timezone
	loc, err := time.LoadLocation(schedule.Timezone)
	if err != nil {
		return fmt.Errorf("invalid timezone: %w", err)
	}

	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	cronSchedule, err := parser.Parse(schedule.CronExpr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	job := cron.NewChain().Then(cron.FuncJob(func() {
		s.executeSchedule(schedule.ID)
	}))

	entryID := s.cron.Schedule(cronSchedule, job)

	s.jobs[schedule.ID] = entryID

	// Update next run time
	nextRun := cronSchedule.Next(time.Now().In(loc))
	schedule.NextRunAt = models.NewNullTime(nextRun)
	if err := s.db.UpdateSchedule(schedule); err != nil {
		logger.Log.Error("Failed to update schedule next_run_at",
			zap.String("schedule_id", schedule.ID.String()),
			zap.Error(err))
	}

	logger.Log.Info("Added schedule to cron",
		zap.String("schedule_id", schedule.ID.String()),
		zap.String("cron_expr", schedule.CronExpr),
		zap.Time("next_run", nextRun))

	return nil
}

func (s *Scheduler) RemoveSchedule(scheduleID uuid.UUID) {
	if entryID, exists := s.jobs[scheduleID]; exists {
		s.cron.Remove(entryID)
		delete(s.jobs, scheduleID)

		logger.Log.Info("Removed schedule from cron",
			zap.String("schedule_id", scheduleID.String()))
	}
}

func (s *Scheduler) executeSchedule(scheduleID uuid.UUID) {
	logger.Log.Info("Executing schedule",
		zap.String("schedule_id", scheduleID.String()))

	// Get schedule details
	schedule, err := s.db.GetSchedule(scheduleID)
	if err != nil {
		logger.Log.Error("Failed to get schedule",
			zap.String("schedule_id", scheduleID.String()),
			zap.Error(err))
		return
	}

	// Check day filter
	if !s.shouldRunToday(schedule) {
		logger.Log.Info("Schedule skipped due to day filter",
			zap.String("schedule_id", scheduleID.String()),
			zap.String("day_filter", schedule.DayFilter.String))
		return
	}

	// Get account (with load balancing if enabled)
	account, err := s.getAccountForSchedule(schedule)
	if err != nil {
		logger.Log.Error("Failed to get account",
			zap.String("account_id", schedule.AccountID.String()),
			zap.Error(err))
		s.logJobExecution(scheduleID, "failed", "", fmt.Sprintf("Account not found: %v", err))
		return
	}

	// Get template
	template, err := s.db.GetTemplate(schedule.TemplateID)
	if err != nil {
		logger.Log.Error("Failed to get template",
			zap.String("template_id", schedule.TemplateID.String()),
			zap.Error(err))
		s.logJobExecution(scheduleID, "failed", "", fmt.Sprintf("Template not found: %v", err))
		return
	}

	// Get all channels for this schedule
	var channelUUIDs []uuid.UUID
	for _, cidStr := range schedule.ChannelIDs {
		if cid, err := uuid.Parse(cidStr); err == nil {
			channelUUIDs = append(channelUUIDs, cid)
		}
	}

	if len(channelUUIDs) == 0 {
		logger.Log.Error("No valid channels in schedule",
			zap.String("schedule_id", scheduleID.String()))
		s.logJobExecution(scheduleID, "failed", "", "No channels configured")
		return
	}

	// Queue messages for each channel with delays
	for i, channelID := range channelUUIDs {
		channel, err := s.db.GetChannel(channelID)
		if err != nil {
			logger.Log.Error("Failed to get channel",
				zap.String("channel_id", channelID.String()),
				zap.Error(err))
			continue
		}

		// Calculate delay for this message
		delay := s.calculateDelay(schedule, i)

		job := &MessageJob{
			ScheduleID: scheduleID,
			Account:    account,
			Template:   template,
			Channel:    channel,
			Message:    template.Content, // TODO: Process variables
			Delay:      delay,
		}

		s.dispatcher.Enqueue(job)
	}

	// Update last run time
	schedule.LastRunAt = models.NewNullTime(time.Now())
	if err := s.db.UpdateSchedule(schedule); err != nil {
		logger.Log.Error("Failed to update schedule last_run_at",
			zap.String("schedule_id", scheduleID.String()),
			zap.Error(err))
	}
}

func (s *Scheduler) shouldRunToday(schedule *models.Schedule) bool {
	now := time.Now()
	weekday := int(now.Weekday())

	if !schedule.DayFilter.Valid || schedule.DayFilter.String == "all" {
		return true
	}

	switch schedule.DayFilter.String {
	case "weekdays":
		return weekday >= 1 && weekday <= 5
	case "weekends":
		return weekday == 0 || weekday == 6
	case "custom":
		for _, day := range schedule.CustomDays {
			if day == weekday {
				return true
			}
		}
		return false
	}

	return true
}

func (s *Scheduler) getAccountForSchedule(schedule *models.Schedule) (*models.Account, error) {
	// If load balancing is not enabled, use the specified account
	if !schedule.LoadBalance {
		return s.db.GetAccount(schedule.AccountID)
	}

	// Get least used active account
	return s.db.GetLeastUsedAccount()
}

func (s *Scheduler) calculateDelay(schedule *models.Schedule, index int) time.Duration {
	if index == 0 || schedule.DelayMaxSeconds == 0 {
		return 0
	}

	minDelay := schedule.DelayMinSeconds
	maxDelay := schedule.DelayMaxSeconds

	if minDelay >= maxDelay {
		return time.Duration(minDelay) * time.Second
	}

	// Random delay between min and max
	randomSeconds := minDelay + (maxDelay-minDelay)*index/10 // Simple distribution
	return time.Duration(randomSeconds) * time.Second
}

func (s *Scheduler) logJobExecution(scheduleID uuid.UUID, status, message, errorMsg string) {
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

	if err := s.db.CreateJobLog(log); err != nil {
		logger.Log.Error("Failed to create job log",
			zap.String("schedule_id", scheduleID.String()),
			zap.Error(err))
	}
}

func (s *Scheduler) Stop() {
	logger.Log.Info("Stopping scheduler")
	s.cron.Stop()
	s.dispatcher.Stop()
}
