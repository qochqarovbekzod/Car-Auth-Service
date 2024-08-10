package token

import (
	"auth/config"
	pb "auth/generated/auth"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.StandardClaims
}

func GenerateJwt(user *pb.UserClaims) (*pb.Tokens, error) {
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:    user.Id,
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(), // Refresh token muddati 7 gun
			IssuedAt:  time.Now().Unix(),
		},
	})
	access, err := accessToken.SignedString([]byte(config.Load().ACCESS_TOKEN))
	if err != nil {
		return nil, err
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:    user.Id,
		Email: user.Email,
		Role:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	})
	refresh, err := refreshToken.SignedString([]byte(config.Load().REFRESH_TOKEN))
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		Accestoken:   access,
		Refreshtoken: refresh,
	}, nil
}

func GenerateAccessToken(refresh string) (*pb.Tokens, error) {

	claims, err := ExstractClaims(refresh)
	if err != nil {
		return nil, err
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		ID:    claims.Id,
		Email: claims.Email,
		Role:  claims.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(30 * time.Minute).Unix(), // Refresh token muddati 7 gun
			IssuedAt:  time.Now().Unix(),
		},
	})
	access, err := accessToken.SignedString([]byte(config.Load().ACCESS_TOKEN))
	if err != nil {
		return nil, err
	}

	return &pb.Tokens{
		Accestoken:   access,
		Refreshtoken: refresh,
	}, nil
}

func ExstractClaims(refresh string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(refresh, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(config.Load().REFRESH_TOKEN), nil
	})

	if err != nil {
		fmt.Println("Error parsing token")
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}

func ExstractClaimsAccess(acces string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(acces, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		
		return []byte(config.Load().ACCESS_TOKEN), nil
	})

	if err != nil {
		fmt.Println("Error parsing token")
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return nil, fmt.Errorf("invalid claims")
	}

	return claims, nil
}
