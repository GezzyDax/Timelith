package telegram

import (
	"context"
	"fmt"
	"sync"

	"github.com/timelith/backend/internal/config"
	"github.com/timelith/backend/internal/models"
	"github.com/timelith/backend/pkg/logger"
	"gorm.io/gorm"
)

type Manager struct {
	config   *config.Config
	db       *gorm.DB
	log      *logger.Logger
	clients  map[uint]*TelegramClient
	mu       sync.RWMutex
}

func NewManager(cfg *config.Config, db *gorm.DB, log *logger.Logger) *Manager {
	return &Manager{
		config:  cfg,
		db:      db,
		log:     log,
		clients: make(map[uint]*TelegramClient),
	}
}

func (m *Manager) Initialize() error {
	m.log.Info("Initializing Telegram Manager...")

	// Load all active accounts from database
	var accounts []models.TelegramAccount
	if err := m.db.Where("status = ?", "authorized").Find(&accounts).Error; err != nil {
		return fmt.Errorf("failed to load accounts: %w", err)
	}

	m.log.Infow("Loading Telegram accounts", "count", len(accounts))

	for _, account := range accounts {
		if err := m.CreateClient(account.ID); err != nil {
			m.log.Errorw("Failed to create client", "account_id", account.ID, "error", err)
		}
	}

	return nil
}

func (m *Manager) CreateClient(accountID uint) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if client already exists
	if _, exists := m.clients[accountID]; exists {
		return nil
	}

	// Load account from database
	var account models.TelegramAccount
	if err := m.db.First(&account, accountID).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Create new Telegram client
	client, err := NewClient(m.config, &account, m.log)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	m.clients[accountID] = client
	m.log.Infow("Telegram client created", "account_id", accountID)

	return nil
}

func (m *Manager) GetClient(accountID uint) (*TelegramClient, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	client, exists := m.clients[accountID]
	if !exists {
		return nil, fmt.Errorf("client not found for account %d", accountID)
	}

	return client, nil
}

func (m *Manager) RemoveClient(accountID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if client, exists := m.clients[accountID]; exists {
		client.Stop()
		delete(m.clients, accountID)
		m.log.Infow("Telegram client removed", "account_id", accountID)
	}
}

func (m *Manager) SendMessage(ctx context.Context, accountID uint, channelID int64, message string) (int64, error) {
	client, err := m.GetClient(accountID)
	if err != nil {
		return 0, err
	}

	return client.SendMessage(ctx, channelID, message)
}

func (m *Manager) SendAuthCode(accountID uint, phoneNumber string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Load account
	var account models.TelegramAccount
	if err := m.db.First(&account, accountID).Error; err != nil {
		return fmt.Errorf("account not found: %w", err)
	}

	// Create temporary client for auth
	client, err := NewClient(m.config, &account, m.log)
	if err != nil {
		return err
	}

	// Send code
	if err := client.SendCode(context.Background(), phoneNumber); err != nil {
		return err
	}

	// Store client temporarily
	m.clients[accountID] = client

	return nil
}

func (m *Manager) VerifyCode(accountID uint, code string) (*models.TelegramAccount, error) {
	client, err := m.GetClient(accountID)
	if err != nil {
		return nil, err
	}

	userInfo, err := client.VerifyCode(context.Background(), code)
	if err != nil {
		return nil, err
	}

	// Update account in database
	var account models.TelegramAccount
	if err := m.db.First(&account, accountID).Error; err != nil {
		return nil, err
	}

	account.Status = "authorized"
	account.TelegramUserID = userInfo.UserID
	account.FirstName = userInfo.FirstName
	account.LastName = userInfo.LastName
	account.Username = userInfo.Username

	if err := m.db.Save(&account).Error; err != nil {
		return nil, err
	}

	return &account, nil
}

func (m *Manager) GetChannels(accountID uint) ([]ChannelInfo, error) {
	client, err := m.GetClient(accountID)
	if err != nil {
		return nil, err
	}

	return client.GetChannels(context.Background())
}
