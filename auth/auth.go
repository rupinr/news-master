package auth

import (
	"crypto/rsa"
	"fmt"
	"net/http"
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
	fmt.Printf("Initing Keys.....")
	privateKey, errPvtKey = loadPrivateKey()
	publicKey, errPubKey = loadPublicKey()
	if errPubKey != nil || errPvtKey != nil {
		panic(fmt.Sprintf("Unable Load Keys...%v", errPubKey))
	}

}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	privateKeyData, err := os.ReadFile(os.Getenv("PRIVATE_KEY_PATH"))
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPrivateKeyFromPEM(privateKeyData)
}

func loadPublicKey() (*rsa.PublicKey, error) {
	publicKeyData, err := os.ReadFile(os.Getenv("PUBLIC_KEY_PATH"))
	if err != nil {
		return nil, err
	}
	return jwt.ParseRSAPublicKeyFromPEM(publicKeyData)
}

func (token *Token) validateAdminToken() *DecodedUser {
	user := defaultDecodedUser()
	if token.Value == os.Getenv("ADMIN_TOKEN") {
		user.Admin = true
		user.Valid = true
	}
	return user
}

func (token *Token) validateSubscriberToken() *DecodedUser {
	user, _ := ValidateJWT(token.Value)
	return user

}

func ValidateAdminToken(token Token) *DecodedUser {
	return token.validateAdminToken()
}

func ValidateSubscriberToken(token Token) *DecodedUser {
	return token.validateSubscriberToken()
}

func AuthMiddleware(validateToken func(Token) *DecodedUser) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := Token{Value: c.Request.Header.Get("Authorization")}
		user := validateToken(token)

		if token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if !user.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Set("user", user)
		c.Next()
	}
}

func ValidateJWT(tokenString string) (*DecodedUser, error) {
	user := defaultDecodedUser()
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return user, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		user.Admin = false
		user.Valid = true
		user.Email = claims["email"].(string)
		user.ID = uint(claims["id"].(float64))

	} else {
		fmt.Println("Invalid token.")
	}

	return user, nil
}

func SubsriberToken(id int, email string, validityInHours int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Duration(validityInHours) * time.Hour).Unix(),
		"id":    id,
	})
	tokenString, err := token.SignedString(privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
