package controllers

import (
	"fmt"
	"net/http"
	"time"
	
	"github.com/dgrijalva/jwt-go"
)

var jwtUserKey = []byte("Bebasapasaja123!")
var userTokenName = "userCookie"

type UserClaims struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UserType string    `json:"usertype"`
	jwt.StandardClaims
}

func generateUserToken(w http.ResponseWriter, id int, name string, userType string){
	tokenExpiryTime := time.Now().Add(24 * time.Hour)

	// create claims with user data
	claims := &UserClaims{
		ID : id,
		Name : name,
		UserType: userType,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: tokenExpiryTime.Unix(),
		},
	}

	// encrypt claim to jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// jwtKey := os.Getenv("JWT_TOKEN")
	jwtToken, err := token.SignedString(jwtUserKey)
	if err != nil{
		return
	}

	// set token to cookies
	http.SetCookie(w, &http.Cookie{
		Name: userTokenName,
		Value: jwtToken,
		Expires: tokenExpiryTime,
		Secure: false,
		HttpOnly: true,
	})
}

func resetUserToken(w http.ResponseWriter){
	http.SetCookie(w, &http.Cookie{
		Name: userTokenName,
		Value: "",
		Expires: time.Now(),
		Secure: false,
		HttpOnly: true,
	})
}

func AuthenticateUser(next http.HandlerFunc, accessType string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		isValidToken := validateUserToken(r, accessType)
		if !isValidToken {
			sendResponse(w, 400, "Token tidak valid!")
		}else{
			next.ServeHTTP(w, r)
		}
	})
}

func validateUserToken(r *http.Request, accessType string) bool {
	isAccessTokenValid, id, name, userType := 
	validateUserTokenFromCookies(r)
	fmt.Print(id, name, userType, accessType, isAccessTokenValid)

	if isAccessTokenValid {
		isUserValid := userType == accessType
		if isUserValid {
			return true
		}
	}
	return false
}

func validateUserTokenFromCookies(r *http.Request) (bool, int, string, string) {
	if cookie, err := r.Cookie(userTokenName); err == nil{
		jwtToken := cookie.Value
		accessClaims := &UserClaims{}
		parsedToken, err := jwt.ParseWithClaims(jwtToken,
		accessClaims, func(accessToken *jwt.Token) (interface{}, error) {
			return jwtUserKey, nil
		})
		if err == nil && parsedToken.Valid {
			return true, accessClaims.ID, accessClaims.Name, accessClaims.UserType
		}
	}
	return false, -1, "", ""
}