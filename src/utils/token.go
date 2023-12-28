package utils

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

var hmacSampleSecret = []byte("your_secret_key") // Replace with your secret key

// CustomClaims includes the standard JWT claims and additional data
type CustomClaims struct {
    jwt.StandardClaims
    UserID string `json:"userId"`
}

func VerifyToken(tokenString string) (*CustomClaims, error) {
    // Initialize a new instance of `CustomClaims`
    claims := &CustomClaims{}

    // Parse the token
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        // Validate the alg is what you expect:
        if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil,errors.New("TUnexpected signing method")
        }

        return hmacSampleSecret, nil
    })

    if err != nil {
        // Handle specific token parsing error (e.g., expired)
        if ve, ok := err.(*jwt.ValidationError); ok {
            if ve.Errors&jwt.ValidationErrorExpired != 0 {
                // Token is expired
                return nil, errors.New("Token is expired")
            } else {
                // Other validation error
                return nil, errors.New("Token is invalid")
            }
        }
        return nil, err
    }

    if !token.Valid {
        return nil, errors.New("Token is invalid")
    }

    return claims, nil
}