package telegram

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"sync"

	"github.com/GezzyDax/timelith/go-backend/internal/config"
	"github.com/GezzyDax/timelith/go-backend/internal/logger"
	"github.com/GezzyDax/timelith/go-backend/internal/models"
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/tg"
	"go.uber.org/zap"
)

type SessionManager struct {
	cfg           *config.Config
	activeClients map[string]*clientEntry
	mu            sync.RWMutex
	gcm           cipher.AEAD
}

type clientEntry struct {
	client  *telegram.Client
	storage *session.StorageMemory
}

var ErrInvalidPassword = errors.New("telegram password invalid")

func NewSessionManager(cfg *config.Config) (*SessionManager, error) {
	// Initialize encryption
	key := []byte(cfg.EncryptionKey)
	if len(key) != 32 {
		return nil, fmt.Errorf("encryption key must be 32 bytes")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCM: %w", err)
	}

	return &SessionManager{
		cfg:           cfg,
		activeClients: make(map[string]*clientEntry),
		gcm:           gcm,
	}, nil
}

func (sm *SessionManager) getEntry(phone string) (*clientEntry, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	entry, exists := sm.activeClients[phone]
	if !exists {
		return nil, fmt.Errorf("telegram client not found for phone %s", phone)
	}
	return entry, nil
}

func (sm *SessionManager) newClientWithSession(ctx context.Context, sessionBytes []byte) (*session.StorageMemory, *telegram.Client, error) {
	storage := &session.StorageMemory{}
	if len(sessionBytes) > 0 {
		if err := storage.StoreSession(ctx, sessionBytes); err != nil {
			return nil, nil, fmt.Errorf("failed to load session: %w", err)
		}
	}

	client := telegram.NewClient(sm.cfg.TelegramAppID, sm.cfg.TelegramAppHash, telegram.Options{
		SessionStorage: storage,
	})

	return storage, client, nil
}

func getPhoneCodeHash(sent tg.AuthSentCodeClass) (string, error) {
	switch v := sent.(type) {
	case *tg.AuthSentCode:
		return v.PhoneCodeHash, nil
	default:
		return "", fmt.Errorf("unsupported sent code type %T", sent)
	}
}

// Encrypt session data before storing in database
func (sm *SessionManager) EncryptSession(data []byte) ([]byte, error) {
	nonce := make([]byte, sm.gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := sm.gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

// Decrypt session data from database
func (sm *SessionManager) DecryptSession(data []byte) ([]byte, error) {
	nonceSize := sm.gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := sm.gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// AuthenticatePhone initiates phone authentication
func (sm *SessionManager) AuthenticatePhone(ctx context.Context, phone string) (string, []byte, error) {
	storage, client, err := sm.newClientWithSession(ctx, nil)
	if err != nil {
		return "", nil, err
	}

	var phoneCodeHash string

	err = client.Run(ctx, func(ctx context.Context) error {
		sent, err := client.Auth().SendCode(ctx, phone, auth.SendCodeOptions{})
		if err != nil {
			return err
		}
		phoneCodeHash, err = getPhoneCodeHash(sent)
		return err
	})

	if err != nil {
		return "", nil, err
	}

	sessionBytes, dumpErr := storage.Bytes(nil)
	if dumpErr != nil {
		return "", nil, fmt.Errorf("failed to dump session after sending code: %w", dumpErr)
	}

	encrypted, encErr := sm.EncryptSession(sessionBytes)
	if encErr != nil {
		return "", nil, encErr
	}

	// No need to keep this client around; drop reference so GC can clean storage.
	return phoneCodeHash, encrypted, nil
}

// VerifyCode verifies the authentication code and completes login if no password is required.
func (sm *SessionManager) VerifyCode(ctx context.Context, phone, code, phoneCodeHash string, pendingSession []byte) (finalSession []byte, nextPending []byte, requiresPassword bool, passwordHint string, err error) {
	storage, client, err := sm.newClientWithSession(ctx, pendingSession)
	if err != nil {
		return nil, nil, false, "", err
	}

	runErr := client.Run(ctx, func(ctx context.Context) error {
		_, err := client.Auth().SignIn(ctx, phone, code, phoneCodeHash)
		if errors.Is(err, auth.ErrPasswordAuthNeeded) {
			pw, pwErr := client.API().AccountGetPassword(ctx)
			if pwErr != nil {
				logger.Log.Warn("Failed to fetch password hint",
					zap.String("phone", phone),
					zap.Error(pwErr))
			} else {
				passwordHint = pw.Hint
			}
		}
		return err
	})

	sessionBytes, dumpErr := storage.Bytes(nil)
	if dumpErr != nil {
		return nil, nil, false, "", fmt.Errorf("failed to dump session after code verification: %w", dumpErr)
	}

	encryptedSession, encErr := sm.EncryptSession(sessionBytes)
	if encErr != nil {
		return nil, nil, false, "", encErr
	}

	if runErr != nil {
		if errors.Is(runErr, auth.ErrPasswordAuthNeeded) {
			return nil, encryptedSession, true, passwordHint, nil
		}
		return nil, nil, false, "", runErr
	}

	return encryptedSession, nil, false, "", nil
}

// VerifyPassword finalizes login when Telegram account has 2FA enabled.
func (sm *SessionManager) VerifyPassword(ctx context.Context, phone string, pendingSession []byte, password string) ([]byte, error) {
	rawPending, err := sm.DecryptSession(pendingSession)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt pending session: %w", err)
	}

	storage, client, err := sm.newClientWithSession(ctx, rawPending)
	if err != nil {
		return nil, err
	}

	if err := client.Run(ctx, func(ctx context.Context) error {
		_, err := client.Auth().Password(ctx, password)
		return err
	}); err != nil {
		if errors.Is(err, auth.ErrPasswordInvalid) {
			return nil, ErrInvalidPassword
		}
		return nil, err
	}

	finalSession, err := storage.Bytes(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to dump session after password verification: %w", err)
	}

	encryptedSession, err := sm.EncryptSession(finalSession)
	if err != nil {
		return nil, err
	}

	return encryptedSession, nil
}

// LoadSession loads a client from encrypted session data
func (sm *SessionManager) LoadSession(ctx context.Context, account *models.Account) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Check if already loaded
	if _, exists := sm.activeClients[account.Phone]; exists {
		return nil
	}

	// Decrypt session
	sessionData, err := sm.DecryptSession(account.SessionData)
	if err != nil {
		return fmt.Errorf("failed to decrypt session: %w", err)
	}

	storage, client, err := sm.newClientWithSession(ctx, sessionData)
	if err != nil {
		return err
	}

	sm.activeClients[account.Phone] = &clientEntry{
		client:  client,
		storage: storage,
	}

	logger.Log.Info("Loaded Telegram session",
		zap.String("phone", account.Phone),
		zap.String("account_id", account.ID.String()))

	return nil
}

// GetClient returns an existing client
func (sm *SessionManager) GetClient(phone string) (*telegram.Client, error) {
	entry, err := sm.getEntry(phone)
	if err != nil {
		return nil, err
	}
	return entry.client, nil
}

// SendMessage sends a message to a chat
func (sm *SessionManager) SendMessage(ctx context.Context, phone, chatID, message string) error {
	client, err := sm.GetClient(phone)
	if err != nil {
		return err
	}

	return client.Run(ctx, func(ctx context.Context) error {
		api := client.API()

		// Resolve peer (can be username, phone, or chat ID)
		peer, err := sm.resolvePeer(ctx, api, chatID)
		if err != nil {
			return fmt.Errorf("failed to resolve peer: %w", err)
		}

		// Send message
		_, err = api.MessagesSendMessage(ctx, &tg.MessagesSendMessageRequest{
			Peer:    peer,
			Message: message,
		})

		return err
	})
}

// resolvePeer resolves a chat ID/username to a Telegram peer
func (sm *SessionManager) resolvePeer(ctx context.Context, api *tg.Client, chatID string) (tg.InputPeerClass, error) {
	// Try to resolve as username
	if len(chatID) > 0 && chatID[0] == '@' {
		resolved, err := api.ContactsResolveUsername(ctx, chatID[1:])
		if err != nil {
			return nil, err
		}

		if len(resolved.Users) > 0 {
			user := resolved.Users[0].(*tg.User)
			return &tg.InputPeerUser{
				UserID:     user.ID,
				AccessHash: user.AccessHash,
			}, nil
		}
	}

	// Otherwise use as chat ID (simplified - in production handle different peer types)
	return &tg.InputPeerEmpty{}, fmt.Errorf("peer resolution not fully implemented")
}

// CloseClient closes and removes a client
func (sm *SessionManager) CloseClient(phone string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	_, exists := sm.activeClients[phone]
	if !exists {

		return nil
	}

	// Note: Client doesn't have Close method in gotd/td
	// We just remove it from the map
	delete(sm.activeClients, phone)

	logger.Log.Info("Closed Telegram client",
		zap.String("phone", phone))

	return nil
}

// SendMediaMessage sends a message with media attachments
func (sm *SessionManager) SendMediaMessage(ctx context.Context, phone, chatID string, template *models.Template) error {
	client, err := sm.GetClient(phone)
	if err != nil {
		return err
	}

	return client.Run(ctx, func(ctx context.Context) error {
		api := client.API()

		peer, err := sm.resolvePeer(ctx, api, chatID)
		if err != nil {
			return fmt.Errorf("failed to resolve peer: %w", err)
		}

		// Handle different media types
		switch template.MediaType.String {
		case "photo":
			return sm.sendPhoto(ctx, api, peer, template)
		case "video":
			return sm.sendVideo(ctx, api, peer, template)
		case "album":
			return sm.sendAlbum(ctx, api, peer, template)
		default:
			return fmt.Errorf("unsupported media type: %s", template.MediaType.String)
		}
	})
}

func (sm *SessionManager) sendPhoto(ctx context.Context, api *tg.Client, peer tg.InputPeerClass, template *models.Template) error {
	// TODO: Implement photo upload and send
	// This requires uploading the file first, then sending it
	logger.Log.Warn("Photo sending not fully implemented yet")
	return fmt.Errorf("photo sending not implemented")
}

func (sm *SessionManager) sendVideo(ctx context.Context, api *tg.Client, peer tg.InputPeerClass, template *models.Template) error {
	// TODO: Implement video upload and send
	logger.Log.Warn("Video sending not fully implemented yet")
	return fmt.Errorf("video sending not implemented")
}

func (sm *SessionManager) sendAlbum(ctx context.Context, api *tg.Client, peer tg.InputPeerClass, template *models.Template) error {
	// TODO: Implement album (media group) upload and send
	logger.Log.Warn("Album sending not fully implemented yet")
	return fmt.Errorf("album sending not implemented")
}

// ForwardMessage forwards a message from one chat to another
func (sm *SessionManager) ForwardMessage(ctx context.Context, phone, toChatID, fromChatID string, messageID int) error {
	client, err := sm.GetClient(phone)
	if err != nil {
		return err
	}

	return client.Run(ctx, func(ctx context.Context) error {
		api := client.API()

		// Resolve destination peer
		toPeer, err := sm.resolvePeer(ctx, api, toChatID)
		if err != nil {
			return fmt.Errorf("failed to resolve destination peer: %w", err)
		}

		// Resolve source peer
		fromPeer, err := sm.resolvePeer(ctx, api, fromChatID)
		if err != nil {
			return fmt.Errorf("failed to resolve source peer: %w", err)
		}

		// Forward message
		_, err = api.MessagesForwardMessages(ctx, &tg.MessagesForwardMessagesRequest{
			FromPeer: fromPeer,
			ToPeer:   toPeer,
			ID:       []int{messageID},
		})

		return err
	})
}

// Close closes all clients
func (sm *SessionManager) Close() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for phone := range sm.activeClients {
		delete(sm.activeClients, phone)
	}

	logger.Log.Info("Closed all Telegram clients")
}
