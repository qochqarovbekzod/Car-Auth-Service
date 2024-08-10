package storage

import (
	"auth/storage/postgres"
	"database/sql"
	"log/slog"
)

type IStorage interface {
	Auth() postgres.AuthRepository
}

type iStorageImpl struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewStorage(db *sql.DB, log *slog.Logger) IStorage {
	return &iStorageImpl{
		db:     db,
		logger: log,
	}
}

func (s *iStorageImpl) Auth() postgres.AuthRepository {
	return postgres.NewAuthRepository(s.db, s.logger)
}
