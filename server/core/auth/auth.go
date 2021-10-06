package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.com/alienspaces/go-boilerplate/server/core/type/configurer"
	"gitlab.com/alienspaces/go-boilerplate/server/core/type/logger"
)

const (
	// AuthTypeJWT -
	AuthTypeJWT string = "jwt"

	// Internal defaults
	jwtDefaultExpiryMinutes int32 = 180
)

// Auth -
type Auth struct {
	Config configurer.Configurer
	Log    logger.Logger
	// JWT -
	JwtSigningKey    string
	JwtExpiryMinutes int32
}

// Claims -
type Claims struct {
	Roles    []string               `json:"roles"`
	Identity map[string]interface{} `json:"identity"`
	jwt.StandardClaims
}

// NewAuth -
func NewAuth(c configurer.Configurer, l logger.Logger) (*Auth, error) {

	j := Auth{
		Log:    l,
		Config: c,
	}

	err := j.Init()
	if err != nil {
		return nil, err
	}

	return &j, nil
}

// Init -
func (j *Auth) Init() error {

	// JWT signing key
	jwtSigningKey := j.Config.Get("APP_SERVER_JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		msg := "APP_SERVER_JWT_SIGNING_KEY not defined, cannot sign JWT"
		j.Log.Warn(msg)
		return fmt.Errorf(msg)
	}

	j.JwtSigningKey = jwtSigningKey
	j.JwtExpiryMinutes = jwtDefaultExpiryMinutes

	return nil
}

// EncodeJWT -
func (j *Auth) EncodeJWT(claims *Claims) (string, error) {

	j.Log.Info("Encoding JWT")

	// Standard claims
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(j.JwtExpiryMinutes)).Unix()

	// Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	j.Log.Info("JWT token >%v< signing key >%s<", token, j.JwtSigningKey)

	// Signed
	tokenString, err := token.SignedString([]byte(j.JwtSigningKey))
	if err != nil {
		j.Log.Warn("Failed signing JWT >%v<", err)
		return "", err
	}

	return tokenString, nil
}

// DecodeJWT -
func (j *Auth) DecodeJWT(tokenString string) (*Claims, error) {

	j.Log.Info("Decoding JWT >%s<", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.JwtSigningKey), nil
	})
	if err != nil {
		j.Log.Warn("Failed parsing JWT claims >%v<", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		j.Log.Info("Expires >%v<", claims.StandardClaims.ExpiresAt)
		j.Log.Info("Roles >%v<", claims.Roles)
		j.Log.Info("Identity >%v<", claims.Identity)
		return claims, nil
	}

	return nil, err
}
