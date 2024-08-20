package postgres

import (
	pb "auth/generated/auth"
	"auth/model"
	"database/sql"
	"fmt"
	"log/slog"
	"time"
)

type AuthRepository interface {
	Register(in *pb.RegistrRequest) (*pb.RegistrResponse, error)
	Login(in *pb.LoginRequest) (*pb.LoginResponse, error)
	RefreshToken(in *pb.RefreshTokenRequest) (*pb.Void, error)
	LogOut(in *pb.TokenRequest) error
	GetUserProfile(in *pb.Id) (*pb.UserProfileResponse, error)
	UpdateUserProfile(in *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error)
	WreateRefreshToken(token model.RefreshTokens) error
}

type AuthRepo struct {
	DB  *sql.DB
	Log *slog.Logger
}

func NewAuthRepository(DB *sql.DB, log *slog.Logger) AuthRepository {
	return &AuthRepo{
		DB:  DB,
		Log: log,
	}
}

func (a *AuthRepo) Register(in *pb.RegistrRequest) (*pb.RegistrResponse, error) {
	var resp pb.RegistrResponse

	err := a.DB.QueryRow(`
			INSERT INTO 
				users(
					email,
                    password,
					first_name,
					last_name,
                    phone_number,
                    role
				)
			 VALUES($1, 
					$2, 
					$3, 
					$4, 
					$5, 
					$6)
			Returning 
					id, 
					email,
					password, 
					first_name, 
					last_name, 
					phone_number, 
					role,
					created_at, 
					updated_at;
					`, in.Email, in.Password, in.FirstName, in.LastName, in.PhoneNumber, in.Role).Scan(
		&resp.Id,
		&resp.Email,
		&resp.Password,
		&resp.FirstName,
		&resp.LastName,
		&resp.PhoneNumber,
		&resp.Role,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil {
		fmt.Println(in.Password, err)
		a.Log.Error("Error registering user", "error", err)
		return nil, err
	}

	return &resp, nil

}

func (a *AuthRepo) Login(in *pb.LoginRequest) (*pb.LoginResponse, error) {
	var resp pb.LoginResponse

	err := a.DB.QueryRow(`
            SELECT
                    id,
                    email,
                    password,
					phone_number,
                    role,
					first_name
            FROM
                    users
            WHERE
                    email = $1 AND
					deleted_at = 0
                    `, in.Email).Scan(
		&resp.Id,
		&resp.Email,
		&resp.Password,
		&resp.PhoneNumber,
		&resp.Role,
		&resp.FullName,
	)

	if err != nil {
		a.Log.Error("Error logging in user", "error", err)
		return nil, err
	}

	return &resp, nil
}

func (a *AuthRepo) RefreshToken(in *pb.RefreshTokenRequest) (*pb.Void, error) {
	var exists string
	err := a.DB.QueryRow(`
	  		SELECT token FROM
			    refresh_tokens
			WHERE 
				token = $1`, in.RefreshToken).Scan(&exists)
	if err != nil {
		a.Log.Error("Error refreshing token", "error", err)
		return nil, err
	}

	if exists == "" {
		a.Log.Error("Error refreshing token", "error", fmt.Errorf("token does not exist"))
		return nil, sql.ErrNoRows
	}
	return &pb.Void{}, nil
}

func (a *AuthRepo) LogOut(in *pb.TokenRequest) error {

	_, err := a.DB.Exec(`
        DELETE FROM refresh_tokens
		WHERE token = $1
		`, in.RefreshToken)

	if err != nil {
		a.Log.Error("Error logging out user", "error", err)
		return err
	}
	return nil

}

func (a AuthRepo) GetUserProfile(in *pb.Id) (*pb.UserProfileResponse, error) {
	a.Log.Info("Getting")
	var resp pb.UserProfileResponse
	err := a.DB.QueryRow(`
        SELECT 
			id, 
			email, 
			password, 
			first_name, 
			last_name, 
			phone_number, 
			role, 
			created_at, 
			updated_at
        FROM users
        WHERE 
			id = $1 AND 
			deleted_at = 0`, in.Id).Scan(
		&resp.Id,
		&resp.Email,
		&resp.Password,
		&resp.FirstName,
		&resp.LastName,
		&resp.PhoneNumber,
		&resp.Role,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)
	if err != nil {
		a.Log.Error("Error getting user profile", "error", err)
		return nil, err
	}

	return &resp, nil
}

func (a AuthRepo) UpdateUserProfile(in *pb.UpdateUserProfileRequest) (*pb.UserProfileResponse, error) {
	var resp pb.UserProfileResponse
	fmt.Println(in.Id)
	err := a.DB.QueryRow(`
        UPDATE users 
        SET email=$1,
		    password=$2,
            first_name=$3,
            last_name=$4,
            phone_number=$5,
			role=$6,
            updated_at=$7
        WHERE 
            id=$8 AND 
            deleted_at = 0
        RETURNING 
            id, 
            email, 
            password, 
            first_name, 
            last_name, 
            phone_number, 
            role, 
            created_at, 
            updated_at`, in.Email, in.Password, in.FirstName, in.LastName, in.PhoneNumber, in.Role, time.Now(), in.Id).Scan(
		&resp.Id,
		&resp.Email,
		&resp.Password,
		&resp.FirstName,
		&resp.LastName,
		&resp.PhoneNumber,
		&resp.Role,
		&resp.CreatedAt,
		&resp.UpdatedAt,
	)

	if err != nil {
		a.Log.Error("Error updating user profile", "error", err)
		return nil, err
	}

	return &resp, nil
}

func (a *AuthRepo) WreateRefreshToken(token model.RefreshTokens) error {

	fmt.Println(token)
	_, err := a.DB.Exec(`
	INSERT INTO
		refresh_tokens(
				user_id,
				token,
				expires_at
			)
		VALUES(
				$1,
				$2,
				$3
			)
	`, token.User_id, token.Token, token.Exp)

	return err
}
