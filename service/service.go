package service

import (
	"auth/api/handler/token"
	pb "auth/generated/auth"
	"auth/model"
	"auth/storage"
	"context"
	"log/slog"
)

type UserService interface {
	Registr(ctx context.Context, in *pb.RegistrRequest) (*pb.RegistrResponse, error)
	Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error)
	LogOut(ctx context.Context, in *pb.TokenRequest) (*pb.Void, error)
	GetUserProfile(ctx context.Context, in *pb.Id) (*pb.UserProfileResponse, error)
	UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error)
	VolideitToken(ctx context.Context, in *pb.Id) (*pb.Void, error)
	RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponce, error)
	WreateRefreshToken(token model.RefreshTokens) error
}

type AuthService struct {
	pb.UnimplementedAuthServiceServer
	Log  *slog.Logger
	User storage.IStorage
}

func NewService(user storage.IStorage, log *slog.Logger) UserService {
	return &AuthService{
		Log:  log,
		User: user,
	}
}

func (s *AuthService) Registr(ctx context.Context, in *pb.RegistrRequest) (*pb.RegistrResponse, error) {
	resp, err := s.User.Auth().Register(in)
	if err != nil {
		s.Log.Error("Error registering")
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	resp, err := s.User.Auth().Login(in)
	if err != nil {
		s.Log.Error("Error logging in")
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) LogOut(ctx context.Context, in *pb.TokenRequest) (*pb.Void, error) {
	err := s.User.Auth().LogOut(in)
	if err != nil {
		s.Log.Error("Error logging out")
		return nil, err
	}
	return &pb.Void{}, nil
}

func (s *AuthService) RefreshToken(ctx context.Context, in *pb.RefreshTokenRequest) (*pb.TokenResponce, error) {
	s.Log.Info("Refresh Token rpc method started")

	_, err := s.User.Auth().RefreshToken(in)
	if err != nil {
		s.Log.Error("Failed refresh token", "error", err)
		return nil, err

	}

	tokens, err := token.GenerateAccessToken(in.RefreshToken)
	if err != nil {
		s.Log.Error("Access token is not generated ", "error", err)
		return nil, err
	}

	return &pb.TokenResponce{AccesToken: tokens.Accestoken, RefreshToken: tokens.Refreshtoken}, nil
}

func (s *AuthService) VolideitToken(ctx context.Context, in *pb.Id) (*pb.Void, error) {
	return nil, nil
}

func (s *AuthService) GetUserProfile(ctx context.Context, in *pb.Id) (*pb.UserProfileResponse, error) {
	resp, err := s.User.Auth().GetUserProfile(in)
	if err != nil {
		s.Log.Error("Error getting user profile")
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) UpdateUserProfile(ctx context.Context, in *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error) {
	resp, err := s.User.Auth().UpdateUserProfile(in)
	if err != nil {
		s.Log.Error("Error updating user profile")
		return nil, err
	}
	return resp, nil
}

func (s *AuthService) WreateRefreshToken(in model.RefreshTokens) error {
	err := s.User.Auth().WreateRefreshToken(in)
	if err != nil {
		s.Log.Error("Error creating refresh token")
		return err
	}
	return nil
}
