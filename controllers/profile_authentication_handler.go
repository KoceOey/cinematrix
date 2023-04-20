package controllers

import (
	// "fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtProfileKey = []byte("Bebasapasaja123!")
var profileTokenName = "userCookie"

type ProfileClaims struct {
	ID   int    `json:"id"`
	Nama string `json:"nama"`
	jwt.StandardClaims
}

func generateProfileToken(w http.ResponseWriter, id int, nama string) {
	tokenExpiryTime := time.Now().Add(24 * time.Hour)

	// create claims with user data
	claims := &ProfileClaims{
		ID:   id,
		Nama: nama,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	// encrypt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// jwtKey := os.Getenv("JWT_TOKEN")
	jwtToken, err := token.SignedString(jwtProfileKey)
	if err != nil {
		return
	}

	// set token to cookies
	http.SetCookie(w, &http.Cookie{
		Name:     profileTokenName,
		Value:    jwtToken,
		Expires:  tokenExpiryTime,
		Secure:   false,
		HttpOnly: true,
	})
}

func resetProfileToken(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     profileTokenName,
		Value:    "",
		Expires:  time.Now(),
		Secure:   false,
		HttpOnly: true,
	})
}

func AuthenticateProfile(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		isValidToken := validateProfileTokenFromCookies(r)
		if !isValidToken {
			sendResponse(w, 400, "Token tidak valid!")
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func validateProfileTokenFromCookies(r *http.Request) bool {
	if cookie, err := r.Cookie(profileTokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &ProfileClaims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken,
			accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
				return jwtProfileKey, nil
			})
		if err == nil && parsedToken.Valid {
			return true
		}
	}
	return false
}

func getProfileTokenData(r *http.Request) (int, string) {
	if cookie, err := r.Cookie(profileTokenName); err == nil {
		jwtToken := cookie.Value
		accessClaims := &ProfileClaims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken,
			accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
				return jwtUserKey, nil
			})
		if err == nil && parsedToken.Valid {
			return accessClaims.ID, accessClaims.Nama
		}
	}
	return -1, ""
}
