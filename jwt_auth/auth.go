package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey=[]byte(os.Getenv("SECRET_KEY"))
func SignIn(w http.ResponseWriter,r *http.Request){
    var creds Credentials
	err:=json.NewDecoder(r.Body).Decode(&creds)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	expectedPassword,ok:=users[creds.Username]
	if !ok || expectedPassword != creds.Password{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	expirationTime := time.Now().Add(time.Minute * 5)
	claims := &Claims{
       Username: creds.Username,
	   RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt:jwt.NewNumericDate(expirationTime),
	   },
	}
	
	 token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	    tokenStr,err:=token.SignedString(jwtKey)
		if err!=nil{
		 w.WriteHeader(http.StatusInternalServerError)
		   return
		}
		http.SetCookie(w,&http.Cookie{
			Name: "token",
			Value: tokenStr,
			Expires: expirationTime,

		})
	 
}