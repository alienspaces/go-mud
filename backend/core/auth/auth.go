package auth

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"gitlab.com/alienspaces/go-mud/backend/core/type/configurer"
	"gitlab.com/alienspaces/go-mud/backend/core/type/logger"
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
	l := j.Logger("Init")

	// JWT signing key
	jwtSigningKey := j.Config.Get("APP_SERVER_JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		msg := "missing APP_SERVER_JWT_SIGNING_KEY, failed initialisation"
		l.Warn(msg)
		return fmt.Errorf(msg)
	}

	j.JwtSigningKey = jwtSigningKey
	j.JwtExpiryMinutes = jwtDefaultExpiryMinutes

	return nil
}

// EncodeJWT -
func (j *Auth) EncodeJWT(claims *Claims) (string, error) {
	l := j.Logger("EncodeJWT")

	l.Info("JWT token expires >%v<", claims.StandardClaims.ExpiresAt)
	l.Info("JWT token roles >%v<", claims.Roles)
	l.Info("JWT token identity >%v<", claims.Identity)

	// Standard claims
	claims.IssuedAt = time.Now().Unix()
	claims.ExpiresAt = time.Now().Add(time.Minute * time.Duration(j.JwtExpiryMinutes)).Unix()

	// Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	l.Info("Signing token with key >%s<", j.JwtSigningKey)

	// Signed
	tokenString, err := token.SignedString([]byte(j.JwtSigningKey))
	if err != nil {
		l.Warn("failed signing JWT >%v<", err)
		return "", err
	}

	return tokenString, nil
}

// DecodeJWT -
func (j *Auth) DecodeJWT(tokenString string) (*Claims, error) {
	l := j.Logger("DecodeJWT")

	l.Info("Decoding JWT token string >%s<", tokenString)

	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.JwtSigningKey), nil
	})
	if err != nil {
		l.Warn("failed parsing JWT claims >%v<", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		l.Info("JWT token expires >%v<", claims.StandardClaims.ExpiresAt)
		l.Info("JWT token roles >%v<", claims.Roles)
		l.Info("JWT token identity >%v<", claims.Identity)
		return claims, nil
	}

	return nil, err
}

// Logger -
func (j *Auth) Logger(functionName string) logger.Logger {
	return j.Log.WithPackageContext("core/auth").WithFunctionContext(functionName)
}
