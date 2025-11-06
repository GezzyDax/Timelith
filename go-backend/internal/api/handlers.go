package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/timelith/backend/internal/models"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

func (s *Server) respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (s *Server) respondError(w http.ResponseWriter, status int, message string) {
	s.respondJSON(w, status, Response{
		Success: false,
		Error:   message,
	})
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	s.respondJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "timelith-backend",
	})
}

// Telegram Auth Handlers

type SendCodeRequest struct {
	PhoneNumber string `json:"phone_number"`
	AccountID   uint   `json:"account_id"`
}

func (s *Server) handleSendCode(w http.ResponseWriter, r *http.Request) {
	var req SendCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.PhoneNumber == "" || req.AccountID == 0 {
		s.respondError(w, http.StatusBadRequest, "Phone number and account ID are required")
		return
	}

	err := s.telegramManager.SendAuthCode(req.AccountID, req.PhoneNumber)
	if err != nil {
		s.log.Errorw("Failed to send auth code", "error", err)
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to send code: %v", err))
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
	})
}

type VerifyCodeRequest struct {
	AccountID uint   `json:"account_id"`
	Code      string `json:"code"`
}

func (s *Server) handleVerifyCode(w http.ResponseWriter, r *http.Request) {
	var req VerifyCodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Code == "" || req.AccountID == 0 {
		s.respondError(w, http.StatusBadRequest, "Code and account ID are required")
		return
	}

	account, err := s.telegramManager.VerifyCode(req.AccountID, req.Code)
	if err != nil {
		s.log.Errorw("Failed to verify code", "error", err)
		s.respondError(w, http.StatusBadRequest, fmt.Sprintf("Verification failed: %v", err))
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"user_id":    account.TelegramUserID,
			"first_name": account.FirstName,
			"last_name":  account.LastName,
			"username":   account.Username,
		},
	})
}

type DisconnectRequest struct {
	AccountID uint `json:"account_id"`
}

func (s *Server) handleDisconnect(w http.ResponseWriter, r *http.Request) {
	var req DisconnectRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	s.telegramManager.RemoveClient(req.AccountID)

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
	})
}

func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	accountIDStr := chi.URLParam(r, "account_id")
	accountID, err := strconv.ParseUint(accountIDStr, 10, 32)
	if err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid account ID")
		return
	}

	var account models.TelegramAccount
	if err := s.db.First(&account, uint(accountID)).Error; err != nil {
		s.respondError(w, http.StatusNotFound, "Account not found")
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"status": account.Status,
		},
	})
}

// Channels Handlers

type SyncChannelsRequest struct {
	AccountID uint `json:"account_id"`
}

func (s *Server) handleSyncChannels(w http.ResponseWriter, r *http.Request) {
	var req SyncChannelsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	channels, err := s.telegramManager.GetChannels(req.AccountID)
	if err != nil {
		s.log.Errorw("Failed to sync channels", "error", err)
		s.respondError(w, http.StatusInternalServerError, fmt.Sprintf("Failed to sync: %v", err))
		return
	}

	// Convert to response format
	channelsData := make([]map[string]interface{}, len(channels))
	for i, ch := range channels {
		channelsData[i] = map[string]interface{}{
			"telegram_id":   ch.TelegramID,
			"name":          ch.Name,
			"type":          ch.Type,
			"username":      ch.Username,
			"title":         ch.Title,
			"members_count": ch.MembersCount,
		}
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"channels": channelsData,
		},
	})
}

// Schedules Handlers

func (s *Server) handleGetSchedules(w http.ResponseWriter, r *http.Request) {
	var schedules []models.Schedule
	err := s.db.Preload("TelegramAccount").
		Preload("MessageTemplate").
		Preload("ScheduleChannels.Channel").
		Where("active = ?", true).
		Find(&schedules).Error

	if err != nil {
		s.log.Errorw("Failed to fetch schedules", "error", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to fetch schedules")
		return
	}

	s.respondJSON(w, http.StatusOK, Response{
		Success: true,
		Data: map[string]interface{}{
			"schedules": schedules,
		},
	})
}

// Send Logs Handlers

func (s *Server) handleCreateSendLog(w http.ResponseWriter, r *http.Request) {
	var log models.SendLog
	if err := json.NewDecoder(r.Body).Decode(&log); err != nil {
		s.respondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := s.db.Create(&log).Error; err != nil {
		s.log.Errorw("Failed to create send log", "error", err)
		s.respondError(w, http.StatusInternalServerError, "Failed to create log")
		return
	}

	s.respondJSON(w, http.StatusCreated, Response{
		Success: true,
		Data: map[string]interface{}{
			"log_id": log.ID,
		},
	})
}
