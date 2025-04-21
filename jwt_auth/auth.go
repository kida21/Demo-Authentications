package main

import (
	"encoding/json"
	"errors"
	"fmt"
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

func Refresh(w http.ResponseWriter,r *http.Request){
	c,err:=r.Cookie("token")
	if err!=nil{
		if errors.Is(err,http.ErrNoCookie){
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tokenString:=c.Value
	claims:=&Claims{}

	token,err:=jwt.ParseWithClaims(tokenString,claims,func(t *jwt.Token) (any, error) {
		 if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{
			return nil,fmt.Errorf("unexpected signing method:%v",t.Header["alg"])
		 }
		 return jwtKey,nil
	})
	if err!=nil{
		if errors.Is(err,jwt.ErrSignatureInvalid){
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if time.Until(claims.ExpiresAt.Time) > 30*time.Second{
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("not expired"))
		return
	}
	expirationTime := time.Now().Add(5*time.Second)
	claims.ExpiresAt = jwt.NewNumericDate(expirationTime)
	token=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
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