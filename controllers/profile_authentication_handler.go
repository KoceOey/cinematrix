package controllers

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/dgrijalva/jwt-go"
)

var jwtProfileKey = []byte("Bebasapasaja123!")
var profileTokenName = "userCookie"

type ProfileClaims struct {
	ID       int    `json:"id"`
	Email     string `json:"email"`
	UserType string    `json:"usertype"`
	jwt.StandardClaims
}

func generateProfileToken(w http.ResponseWriter, id int, email string, userType string){
	tokenExpiryTime := time.Now().Add(24 * time.Hour)

	// create claims with user data
	claims := &ProfileClaims{
		ID : id,
		Email : email,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	// encrypt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// jwtKey := os.Getenv("JWT_TOKEN")
	jwtToken, err := token.SignedString(jwtProfileKey)
	if err != nil{
		return
	}

	// set token to cookies
	http.SetCookie(w, &http.Cookie{
		Name: profileTokenName,
		Value: jwtToken,
		Expires: tokenExpiryTime,
		Secure: false,
		HttpOnly: true,
	})
}

func resetProfileToken(w http.ResponseWriter){
	http.SetCookie(w, &http.Cookie{
		Name: profileTokenName,
		Value: "",
		Expires: time.Now(),
		Secure: false,
		HttpOnly: true,
	})
}

func AuthenticateProfile(next http.HandlerFunc, accessType string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		isValidToken := validateProfileToken(r, accessType)
		if !isValidToken {
			sendResponse(w, 400, "Token tidak valid!")
		}else{
			next.ServeHTTP(w, r)
		}
	})
}

func validateProfileToken(r *http.Request, accessType string) bool {
	isAccessTokenValid, id, name, userType := validateProfileTokenFromCookies(r)
	fmt.Print(id, name, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateProfileTokenFromCookies(r *http.Request) (bool, int, string, string) {
	if cookie, err := r.Cookie(profileTokenName); err == nil{
		jwtToken := cookie.Value
		accessClaims := &ProfileClaims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken,
		accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtProfileKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Email, accessClaims.UserType
		}
	}
	return false, -1, "", ""
}