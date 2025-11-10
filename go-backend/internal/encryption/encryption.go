package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
)

var masterKey []byte

// InitMasterKey инициализирует master encryption key.
// Ключ генерируется при первом запуске и сохраняется в .env.master
// При последующих запусках ключ загружается из файла
func InitMasterKey() error {
	keyFile := ".env.master"

	// Пробуем загрузить существующий ключ из файла
	if data, err := os.ReadFile(keyFile); err == nil {
		decoded, err := base64.StdEncoding.DecodeString(string(data))
		if err != nil {
			return fmt.Errorf("failed to decode master key: %w", err)
		}
		masterKey = decoded
		return nil
	}

	// Пробуем загрузить из переменной окружения
	if keyStr := os.Getenv("MASTER_ENCRYPTION_KEY"); keyStr != "" {
		decoded, err := base64.StdEncoding.DecodeString(keyStr)
		if err != nil {
			return fmt.Errorf("failed to decode MASTER_ENCRYPTION_KEY: %w", err)
		}
		masterKey = decoded
		return nil
	}

	// Генерируем новый ключ (32 байта для AES-256)
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return fmt.Errorf("failed to generate random key: %w", err)
	}
	masterKey = key

	// Сохраняем в файл
	encoded := base64.StdEncoding.EncodeToString(key)
	if err := os.WriteFile(keyFile, []byte(encoded), 0600); err != nil {
		return fmt.Errorf("failed to save master key: %w", err)
	}

	fmt.Printf("✓ Generated new master encryption key and saved to %s\n", keyFile)
	return nil
}

// Encrypt шифрует plaintext используя AES-256-GCM
func Encrypt(plaintext string) (string, error) {
	if len(masterKey) == 0 {
		return "", errors.New("master key not initialized")
	}

	if plaintext == "" {
		return "", errors.New("plaintext cannot be empty")
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	// Создаем nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to create nonce: %w", err)
	}

	// Шифруем
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	// Кодируем в base64 для хранения в БД
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt расшифровывает ciphertext используя AES-256-GCM
func Decrypt(ciphertext string) (string, error) {
	if len(masterKey) == 0 {
		return "", errors.New("master key not initialized")
	}

	if ciphertext == "" {
		return "", errors.New("ciphertext cannot be empty")
	}

	// Декодируем из base64
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintext), nil
}

// GetMasterKey возвращает копию master key для тестирования (не использовать в продакшне)
func GetMasterKey() []byte {
	keyCopy := make([]byte, len(masterKey))
	copy(keyCopy, masterKey)
	return keyCopy
}
