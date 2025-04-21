package main

import (
	"encoding/json"
	"log"
	"net/http"
	 "os"
	"time"
     "github.com/golang-jwt/jwt/v5"
	
)

var users = map[string]string{
	"kebede": "password1",
	"abebe":  "password2",
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type Claims struct{
	Username string `json:"username"`
	jwt.RegisteredClaims
}

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
	jwtKey:=[]byte(os.Getenv("SECRET_KEY"))
	 token:=jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	    tokenStr,err:=token.SignedString(jwtKey)
		if err!=nil{
			log.Println(err)
		   w.WriteHeader(http.StatusInternalServerError)
		   return
		}
		http.SetCookie(w,&http.Cookie{
			Name: "token",
			Value: tokenStr,
			Expires: expirationTime,

		})
	 
}


func main() {
   http.HandleFunc("/signin",SignIn)
   log.Fatal(http.ListenAndServe(":8080",nil))
}