package auth

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	"news-master/app"
	"news-master/logger"
	"news-master/repository"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Token struct {
	Value string
}

type DecodedUser struct {
	Valid bool
	Admin bool
	ID    uint
	Email string
}

func defaultDecodedUser() *DecodedUser {
	return &DecodedUser{
		Valid: false,
		Admin: false,
		ID:    0,
		Email: "",
	}
}

var (
	privateKey *rsa.PrivateKey
	errPvtKey  error
	publicKey  *rsa.PublicKey
	errPubKey  error
)

func LoadKeys() {
	logger.Log.Info("Loading Private Public Keys")
	privateKey, errPvtKey = loadPrivateKey()
	publicKey, errPubKey = loadPublicKey()
	if errPubKey != nil || errPvtKey != nil {
		panic(fmt.Sprintf("Unable Load Keys %v %v", errPubKey, errPvtKey))
	}

}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(app.Config.PrivateKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
}

func loadPublicKey() (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(app.Config.PublicKeyPath)
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
}

func (token *Token) validateAdminToken() (*DecodedUser, error) {
	user := defaultDecodedUser()
	if token.Value == app.Config.AdminToken {
		user.Admin = true
		user.Valid = true
		return user, nil
	} else {
		return nil, errors.New("invalid admin token")
	}

}

func (token *Token) validateSubscriberToken() (*DecodedUser, error) {
	user, error := ValidateJWT(token.Value)
	return user, error

}

func ValidateAdminToken(token Token) (*DecodedUser, error) {
	return token.validateAdminToken()
}

func ValidateSubscriberToken(token Token) (*DecodedUser, error) {
	return token.validateSubscriberToken()
}

func AuthMiddleware(validateToken func(Token) (*DecodedUser, error)) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := Token{Value: c.Request.Header.Get("Authorization")}
		if token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		user, err := validateToken(token)

		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if user != nil && !user.Admin {
			_, userErr := repository.GetUserByEmail(user.Email)
			if userErr != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
				c.Abort()
				return
			}
		}

		if user != nil && !user.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid User"})
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}

func ValidateJWT(tokenString string) (*DecodedUser, error) {
	user := defaultDecodedUser()
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})
	if err != nil {
		logger.Log.Error("Error parsing token:", "error", err.Error())
		return nil, err
	}
	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		user.Admin = false
		user.Valid = true
		user.Email = claims.Email
		user.ID = claims.ID
		resetLoginAttemptCounter(user.ID)
	} else {
		logger.Log.Debug("Invalid token.")
	}

	return user, nil
}

func resetLoginAttemptCounter(userId uint) {
	repository.ResetLoginCounter(userId)
}

func SubscriberToken(id uint, email string, validityInHours int) (string, error) {
	claims := CustomClaims{
		email,
		id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(validityInHours) * time.Hour)),
			Issuer:    "api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func User(c *gin.Context) *DecodedUser {
	contextUser, _ := c.Get("user")
	user := contextUser.(*DecodedUser)
	return user
}

type CustomClaims struct {
	Email string `json:"Email"`
	ID    uint   `json:"id"`
	jwt.RegisteredClaims
}
