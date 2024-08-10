package handler

import (
	"auth/logs"
	"auth/service"
	"log/slog"
)

type Handler struct {
	Log  *slog.Logger
	User service.UserService
}

func NewHandler(auth service.UserService) *Handler {
	logs.InitLogger()

	return &Handler{
		Log:  logs.Logger,
		User: auth,
	}
}
