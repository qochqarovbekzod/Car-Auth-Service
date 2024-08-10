package postgres

import (
	pb "auth/generated/auth"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	db, err := ConnectDB()
	if err != nil {
		t.Error("Failed connect to data base", "error", err)
		return
	}

	repo := NewAuthRepository(db, &slog.Logger{})

	req := pb.RegistrRequest{
		Email:       "jansdfkn",
		Password:    "1234",
		PhoneNumber: "32453t64324",
		FirstName:   "23regtrhf",
		LastName:    "efmdvdfkd",
		Role:        "Read",
	}
	res, err := repo.Register(&req)

	if err != nil {
		t.Error("Failed register user", "error", err)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, res.FirstName, req.FirstName)
	assert.Equal(t, req.Email, req.Email)
}
