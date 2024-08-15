package usecases

import (
	"errors"
	"fmt"

	"github.com/Homyakadze14/RecipeSite/internal/entities"
	jwt "github.com/golang-jwt/jwt"
)

type JWTUseCase struct {
	secretKey []byte
}

var (
	ErrBadToken = errors.New("bad token")
)

func NewJWTUseCase(sk []byte) *JWTUseCase {
	return &JWTUseCase{
		secretKey: sk,
	}
}

func (u *JWTUseCase) GenerateJWT(userID int) (*entities.JWTToken, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
	})
	tokenString, err := token.SignedString(u.secretKey)
	if err != nil {
		return nil, err
	}
	return &entities.JWTToken{Token: tokenString}, nil
}

func (u *JWTUseCase) GetDataFromJWT(inToken *entities.JWTToken) (*entities.JWTData, error) {
	hashSecretGetter := func(token *jwt.Token) (interface{}, error) {
		method, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok || method.Alg() != "HS256" {
			return nil, fmt.Errorf("bad sing method")
		}
		return u.secretKey, nil
	}
	token, err := jwt.Parse(inToken.Token, hashSecretGetter)
	if err != nil || !token.Valid {
		return nil, ErrBadToken
	}

	payload, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrBadToken
	}

	return &entities.JWTData{UserID: payload["user_id"]}, nil
}
