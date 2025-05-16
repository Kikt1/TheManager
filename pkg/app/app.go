package app

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/Kikt1/TheManager/database"
	"github.com/Kikt1/TheManager/models"
)

// App holds your app state
type App struct {
  ctx context.Context
}

func NewApp() *App { return &App{} }

func (a *App) Startup(ctx context.Context) error {
  a.ctx = ctx
  if err := database.InitDB(); err != nil {
    return fmt.Errorf("failed to init DB: %w", err)
  }
  return nil
}

func (a *App) Shutdown(ctx context.Context) error {
  if err := database.CloseDB(); err != nil {
    return fmt.Errorf("failed to close DB: %w", err)
  }
  return nil
}

type LoginResponse struct {
  Success bool   `json:"success"`
  Message string `json:"message"`
  UserID  int    `json:"userId,omitempty"`
  Name    string `json:"name,omitempty"`
  Role    string `json:"role,omitempty"`
}

func (a *App) Login(rawPin string) LoginResponse {
  pinHash := hashPin(rawPin)
  user, err := models.ValidateUserPin(pinHash)
  if err != nil {
    fmt.Printf("Error validating PIN: %v\n", err)
    return LoginResponse{Success: false, Message: "Authentication error"}
  }
  if user == nil {
    return LoginResponse{Success: false, Message: "Invalid PIN"}
  }
  return LoginResponse{true, "Login successful", user.ID, user.Name, user.Role}
}

func hashPin(pin string) string {
  h := sha256.Sum256([]byte(pin))
  return hex.EncodeToString(h[:])
}
