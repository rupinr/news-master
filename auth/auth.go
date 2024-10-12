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

var (
	privateKey *rsa.PrivateKey
	errPvtKey  error
	publicKey  *rsa.PublicKey
	errPubKey  error
)

func InitKeys() {
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

func (token *Token) validAdminToken() bool {
	return token.Value == os.Getenv("ADMIN_TOKEN")
}

func (token *Token) validSubscriberToken() bool {
	_, err := ValidateJWT(token.Value)
	if err == nil {
		return true
	} else {
		return false
	}
}

func ValidateAdminToken(token Token) bool {
	return token.validAdminToken()
}

func ValidateSubscriberToken(token Token) bool {
	return token.validSubscriberToken()
}

func AuthMiddleware(validateToken func(Token) bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		token := Token{Value: c.Request.Header.Get("Authorization")}
		if token.Value == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing token"})
			c.Abort()
			return
		}
		if !validateToken(token) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
		c.Next()
	}
}

func ValidateJWT(tokenString string) (*jwt.Token, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println("Token is valid.")
		fmt.Println("Username:", claims["email"])
		fmt.Println("Expiration:", claims["exp"])
	} else {
		fmt.Println("Invalid token.")
	}

	return token, nil
}

func SubsriberToken(email string, validityInHours int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
		"email": email,
		"exp":   time.Now().Add(time.Duration(validityInHours) * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString(&privateKey)
	if err != nil {
		return "", err
	}

	return tokenString, err
}
