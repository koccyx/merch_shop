package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
    UserId string
    jwt.RegisteredClaims
}

func NewToken(userId string, secret string) (string, error) {    
    claims := UserClaims{
        UserId: userId,
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ParseToken(token string, secret string) (string, error) {    
    prsdToken, err := jwt.ParseWithClaims(token, &UserClaims{}, func(token *jwt.Token) (interface{}, error) {
		
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("parsing token error")
		}

		return []byte(secret), nil
	})
    if err != nil {
        return "", err 
    }

    if claims, ok := prsdToken.Claims.(*UserClaims); ok && prsdToken.Valid {
		return claims.UserId, nil
	} 

    return "", fmt.Errorf("parsing token error")
}