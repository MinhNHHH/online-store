package store

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/golang-jwt/jwt/v4"
)

var jwtTokenExpiry = time.Minute * 15
var refreshTokenExpiry = time.Hour * 24

type TokenPairs struct {
	Token        string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type Claims struct {
	UserName string `json:"name"`
	jwt.RegisteredClaims
}

func (app *OnlineStore) getTokenFromHeaderandVerify(w http.ResponseWriter, r *http.Request) (string, *Claims, error) {
	w.Header().Add("Vary", "Authorization")

	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", nil, errors.New("no auth header")
	}

	headerParts := strings.Split(authHeader, " ")
	if len(headerParts) != 2 {
		return "", nil, errors.New("invalid off header")
	}

	if headerParts[0] != "Bearer" {
		return "", nil, errors.New("unauthorized: no Bearer")
	}

	token := headerParts[1]

	claims := &Claims{}

	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(app.Cfgs.JWT_SECRET), nil
	})

	if err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			return "", nil, errors.New("expried token")
		}
		return "", nil, err
	}

	if claims.Issuer != app.Cfgs.DOMAIN {
		return "", nil, errors.New("incorrect issuer")
	}

	return token, claims, nil
}

func (app *OnlineStore) generateTokenPair(user *schema.User) (TokenPairs, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = user.Name
	claims["sub"] = fmt.Sprint(user.ID)
	claims["aud"] = app.Cfgs.DOMAIN
	claims["iss"] = app.Cfgs.DOMAIN
	if user.IsAdmin {
		claims["admin"] = true
	} else {
		claims["admin"] = false
	}

	claims["exp"] = time.Now().Add(jwtTokenExpiry).Unix()

	signedAccessToken, err := token.SignedString([]byte(app.Cfgs.JWT_SECRET))
	if err != nil {
		return TokenPairs{}, err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshTokenClaims := refreshToken.Claims.(jwt.MapClaims)

	refreshTokenClaims["sub"] = fmt.Sprint(user.ID)

	refreshTokenClaims["exp"] = time.Now().Add(refreshTokenExpiry).Unix()

	signedRefreshToken, err := refreshToken.SignedString([]byte(app.Cfgs.JWT_SECRET))
	if err != nil {
		return TokenPairs{}, nil
	}

	var tokenPairs = TokenPairs{
		Token:        signedAccessToken,
		RefreshToken: signedRefreshToken,
	}
	return tokenPairs, nil
}
