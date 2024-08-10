package model

type RefreshTokens struct {
	User_id string `json:"user_id"`
	Token   string `json:"token"`
	Exp     int64  `json:"exp"`
}
