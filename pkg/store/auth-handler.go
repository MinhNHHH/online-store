package store

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/MinhNHHH/online-store/pkg/databases/schema"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Credentials struct {
	UserName string `json:"email"`
	Password string `json:"password"`
}

func (app *OnlineStore) authenticate(w http.ResponseWriter, r *http.Request) {
	var creds Credentials

	// read a json payload
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		app.SendResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// lock up the user by email address
	user, err := app.DB.GetUserByEmail(creds.UserName)
	if err != nil {
		log.Println("Error getting user by email:", err)
		app.SendResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password))
	if err != nil {
		log.Println("Error comparing password:", err)
		app.SendResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// generate tokens
	tokenPairs, err := app.generateTokenPair(user)
	if err != nil {
		log.Println("Error generating token pair:", err)
		app.SendResponse(w, http.StatusUnauthorized, errors.New("unauthorized"))
		return
	}
	// send token to user
	app.SendResponse(w, http.StatusOK, tokenPairs)
}

func (app *OnlineStore) refresh(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println("Error parsing form:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	refreshToken := r.Form.Get("refresh_token")
	claims := &Claims{}
	_, err = jwt.ParseWithClaims(refreshToken, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(app.Cfgs.JWT_SECRET), nil
	})

	if err != nil {
		log.Println("Error parsing refresh token:", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	if time.Unix(claims.ExpiresAt.Unix(), 0).Sub(time.Now()) > 30*time.Second {
		log.Println("Refresh token does not need renewed yet")
		app.SendResponse(w, http.StatusTooEarly, errors.New("refresh toekn does not need renewed yet"))
		return
	}

	// get the user id from the claims
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		log.Println("Error getting user id from claims:", err)
		app.SendResponse(w, http.StatusBadRequest, errors.New("unknow user"))
		return
	}

	user, err := app.DB.GetUser(userID)
	if err != nil {
		log.Println("Error getting user:", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	tokenPairs, err := app.generateTokenPair(user)
	if err != nil {
		log.Println("Error generating token pair:", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "__Host-refresh_token",
		Path:     "/",
		Value:    tokenPairs.RefreshToken,
		Expires:  time.Now().Add(refreshTokenExpiry),
		MaxAge:   int(refreshTokenExpiry.Seconds()),
		SameSite: http.SameSiteStrictMode,
		Domain:   "localhost",
		HttpOnly: true,
		Secure:   true,
	})

	app.SendResponse(w, http.StatusOK, tokenPairs)
}

func (app *OnlineStore) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user schema.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	userID, err := app.DB.InsertUser(user)
	if err != nil {
		log.Println("Error inserting user:", err)
		app.SendResponse(w, http.StatusBadRequest, err)
		return
	}
	app.SendResponse(w, http.StatusOK, userID)
}
