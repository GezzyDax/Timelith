package settings

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/GezzyDax/timelith/go-backend/internal/database"
	"github.com/GezzyDax/timelith/go-backend/internal/encryption"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/google/uuid"
)

// Service предоставляет thread-safe доступ к настройкам с кешированием и hot reload
type Service struct {
	db       *database.DB
	cache    map[string]*models.Setting
	mu       sync.RWMutex
	stopChan chan struct{}
}

// NewService создает новый экземпляр сервиса настроек
func NewService(db *database.DB) (*Service, error) {
	s := &Service{
		db:       db,
		cache:    make(map[string]*models.Setting),
		stopChan: make(chan struct{}),
	}

	// Загружаем настройки из БД при инициализации
	if err := s.reload(); err != nil {
		log.Printf("Warning: Failed to load settings from database: %v", err)
		// Не возвращаем ошибку, так как настройки могут быть еще не созданы
	}

	// Запускаем фоновую горутину для hot reload (каждые 30 секунд)
	go s.reloadLoop()

	return s, nil
}

// reload загружает все настройки из БД в кеш
func (s *Service) reload() error {
	settings, err := s.db.GetAllSettings()
	if err != nil {
		return fmt.Errorf("failed to load settings: %w", err)
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// Очищаем старый кеш
	s.cache = make(map[string]*models.Setting)

	// Загружаем настройки в кеш
	for i := range settings {
		s.cache[settings[i].Key] = &settings[i]
	}

	log.Printf("✓ Loaded %d settings from database", len(settings))
	return nil
}

// reloadLoop периодически перезагружает настройки из БД
func (s *Service) reloadLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := s.reload(); err != nil {
				log.Printf("Warning: Failed to reload settings: %v", err)
			}
		case <-s.stopChan:
			return
		}
	}
}

// Get возвращает значение настройки по ключу
// Если настройка зашифрована, автоматически расшифровывает значение
// Если настройка не найдена в БД, пытается получить из переменной окружения
func (s *Service) Get(key string) (string, error) {
	s.mu.RLock()
	setting, exists := s.cache[key]
	s.mu.RUnlock()

	if !exists {
		// Fallback на переменную окружения
		if envValue := os.Getenv(key); envValue != "" {
			return envValue, nil
		}
		return "", fmt.Errorf("setting not found: %s", key)
	}

	// Если значение зашифровано, расшифровываем
	if setting.Encrypted {
		decrypted, err := encryption.Decrypt(setting.Value)
		if err != nil {
			return "", fmt.Errorf("failed to decrypt setting %s: %w", key, err)
		}
		return decrypted, nil
	}

	return setting.Value, nil
}

// GetWithDefault возвращает значение настройки или дефолтное значение
func (s *Service) GetWithDefault(key, defaultValue string) string {
	value, err := s.Get(key)
	if err != nil {
		return defaultValue
	}
	return value
}

// GetBool возвращает булево значение настройки
func (s *Service) GetBool(key string) (bool, error) {
	value, err := s.Get(key)
	if err != nil {
		return false, err
	}
	return value == "true" || value == "1" || value == "yes", nil
}

// Set сохраняет настройку в БД (с шифрованием если требуется)
// После сохранения автоматически обновляет кеш
func (s *Service) Set(key, value string, encrypted bool, category string, updatedBy *uuid.UUID) error {
	valueToStore := value

	// Шифруем значение если требуется
	if encrypted {
		encryptedValue, err := encryption.Encrypt(value)
		if err != nil {
			return fmt.Errorf("failed to encrypt setting: %w", err)
		}
		valueToStore = encryptedValue
	}

	// Сохраняем в БД с использованием UPSERT
	if err := s.db.UpsertSetting(key, valueToStore, encrypted, category, updatedBy); err != nil {
		return fmt.Errorf("failed to save setting: %w", err)
	}

	// Немедленно обновляем кеш после записи
	return s.reload()
}

// SetBulk сохраняет несколько настроек одной транзакцией
func (s *Service) SetBulk(settings map[string]string, encrypted bool, category string, updatedBy *uuid.UUID) error {
	for key, value := range settings {
		if err := s.Set(key, value, encrypted, category, updatedBy); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}
	return nil
}

// Delete удаляет настройку из БД и кеша
func (s *Service) Delete(key string) error {
	if err := s.db.DeleteSetting(key); err != nil {
		return fmt.Errorf("failed to delete setting: %w", err)
	}

	s.mu.Lock()
	delete(s.cache, key)
	s.mu.Unlock()

	return nil
}

// GetAll возвращает все настройки (с расшифровкой)
func (s *Service) GetAll() ([]models.Setting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	settings := make([]models.Setting, 0, len(s.cache))
	for _, setting := range s.cache {
		// Копируем настройку
		settingCopy := *setting

		// Расшифровываем если требуется
		if setting.Encrypted {
			decrypted, err := encryption.Decrypt(setting.Value)
			if err != nil {
				log.Printf("Warning: Failed to decrypt setting %s: %v", setting.Key, err)
				continue
			}
			settingCopy.Value = decrypted
		}

		settings = append(settings, settingCopy)
	}

	return settings, nil
}

// GetByCategory возвращает все настройки в определенной категории
func (s *Service) GetByCategory(category string) ([]models.Setting, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	settings := make([]models.Setting, 0)
	for _, setting := range s.cache {
		if setting.Category == category {
			settingCopy := *setting

			// Расшифровываем если требуется
			if setting.Encrypted {
				decrypted, err := encryption.Decrypt(setting.Value)
				if err != nil {
					log.Printf("Warning: Failed to decrypt setting %s: %v", setting.Key, err)
					continue
				}
				settingCopy.Value = decrypted
			}

			settings = append(settings, settingCopy)
		}
	}

	return settings, nil
}

// IsSetupCompleted проверяет, завершен ли процесс setup
func (s *Service) IsSetupCompleted() bool {
	completed, _ := s.GetBool("setup_completed")
	return completed
}

// MarkSetupCompleted отмечает setup как завершенный
func (s *Service) MarkSetupCompleted(updatedBy *uuid.UUID) error {
	return s.Set("setup_completed", "true", false, "system", updatedBy)
}

// Stop останавливает фоновый reload loop
func (s *Service) Stop() {
	close(s.stopChan)
}
