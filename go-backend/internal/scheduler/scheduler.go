package scheduler

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/robfig/cron/v3"
	"github.com/timelith/backend/internal/models"
	"github.com/timelith/backend/internal/telegram"
	"github.com/timelith/backend/pkg/logger"
	"gorm.io/gorm"
)

type Scheduler struct {
	db              *gorm.DB
	redis           *redis.Client
	telegramManager *telegram.Manager
	log             *logger.Logger
	cron            *cron.Cron
	stopChan        chan struct{}
}

func New(db *gorm.DB, redisClient *redis.Client, tm *telegram.Manager, log *logger.Logger) *Scheduler {
	return &Scheduler{
		db:              db,
		redis:           redisClient,
		telegramManager: tm,
		log:             log,
		cron:            cron.New(cron.WithSeconds()),
		stopChan:        make(chan struct{}),
	}
}

func (s *Scheduler) Start() {
	s.log.Info("Starting scheduler...")

	s.cron.Start()

	// Start a goroutine to check for due schedules
	go s.checkSchedulesLoop()

	s.log.Info("Scheduler started successfully")
}

func (s *Scheduler) Stop() {
	s.log.Info("Stopping scheduler...")
	close(s.stopChan)
	s.cron.Stop()
	s.log.Info("Scheduler stopped")
}

func (s *Scheduler) checkSchedulesLoop() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkAndExecuteSchedules()
		case <-s.stopChan:
			return
		}
	}
}

func (s *Scheduler) checkAndExecuteSchedules() {
	ctx := context.Background()

	// Find all schedules that are due for execution
	var schedules []models.Schedule
	now := time.Now()

	err := s.db.Preload("TelegramAccount").
		Preload("MessageTemplate").
		Preload("ScheduleChannels.Channel").
		Where("active = ? AND next_run_at <= ?", true, now).
		Find(&schedules).Error

	if err != nil {
		s.log.Errorw("Failed to fetch schedules", "error", err)
		return
	}

	if len(schedules) == 0 {
		return
	}

	s.log.Infow("Found schedules to execute", "count", len(schedules))

	for _, schedule := range schedules {
		go s.executeSchedule(ctx, &schedule)
	}
}

func (s *Scheduler) executeSchedule(ctx context.Context, schedule *models.Schedule) {
	s.log.Infow("Executing schedule",
		"schedule_id", schedule.ID,
		"schedule_name", schedule.Name)

	// Get Telegram account
	account := schedule.TelegramAccount
	if account.Status != "authorized" {
		s.log.Warnw("Account not authorized",
			"schedule_id", schedule.ID,
			"account_id", account.ID)
		return
	}

	// Get message content
	message := schedule.MessageTemplate.Content

	// Send to all channels
	for _, sc := range schedule.ScheduleChannels {
		go s.sendToChannel(ctx, schedule, &account, &sc.Channel, message)
	}

	// Update schedule's next run time
	s.updateScheduleNextRun(schedule)
}

func (s *Scheduler) sendToChannel(ctx context.Context, schedule *models.Schedule, account *models.TelegramAccount, channel *models.Channel, message string) {
	log := &models.SendLog{
		ScheduleID:        schedule.ID,
		TelegramAccountID: account.ID,
		ChannelID:         channel.ID,
		Status:            "pending",
		MessageContent:    message,
	}

	// Create log entry
	if err := s.db.Create(log).Error; err != nil {
		s.log.Errorw("Failed to create send log", "error", err)
		return
	}

	// Update status to sending
	log.Status = "sending"
	s.db.Save(log)

	// Send message via Telegram
	messageID, err := s.telegramManager.SendMessage(ctx, account.ID, channel.TelegramID, message)

	if err != nil {
		s.log.Errorw("Failed to send message",
			"schedule_id", schedule.ID,
			"channel_id", channel.ID,
			"error", err)

		log.Status = "failed"
		log.ErrorMessage = err.Error()
	} else {
		s.log.Infow("Message sent successfully",
			"schedule_id", schedule.ID,
			"channel_id", channel.ID,
			"message_id", messageID)

		now := time.Now()
		log.Status = "sent"
		log.TelegramMessageID = messageID
		log.SentAt = &now
	}

	s.db.Save(log)
}

func (s *Scheduler) updateScheduleNextRun(schedule *models.Schedule) {
	now := time.Now()
	var nextRun time.Time

	switch schedule.ScheduleType {
	case "interval":
		nextRun = now.Add(time.Duration(schedule.IntervalMinutes) * time.Minute)
	case "cron":
		// In production, use a proper cron parser
		nextRun = now.Add(1 * time.Hour)
	case "once":
		// One-time schedules should be deactivated after execution
		schedule.Active = false
	}

	if schedule.Active {
		schedule.NextRunAt = &nextRun
	}

	schedule.LastRunAt = &now

	if err := s.db.Save(schedule).Error; err != nil {
		s.log.Errorw("Failed to update schedule", "schedule_id", schedule.ID, "error", err)
	}
}
