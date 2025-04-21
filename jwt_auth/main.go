package main

import (
	
	"errors"
	"fmt"
	"log"
	"net/http"
	

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
func Welcome(w http.ResponseWriter,r *http.Request){
   c,err:= r.Cookie("token")
   if err!= nil{
	if errors.Is(err,http.ErrNoCookie){
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.WriteHeader(http.StatusBadRequest)
	return
   }
   tokenStr:=c.Value
   claims :=&Claims{}

   token,err:=jwt.ParseWithClaims(tokenStr,claims,func(t *jwt.Token) (any, error) {
	  if _,ok:=t.Method.(*jwt.SigningMethodHMAC);!ok{
		w.WriteHeader(http.StatusUnauthorized)
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
	if !token.Valid{
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	w.Write([]byte(fmt.Sprintf("Welcome,%s",claims.Username)))
}

func main() {
   http.HandleFunc("/signin",SignIn)
   http.HandleFunc("/welcome",Welcome)
   http.HandleFunc("/refresh",Refresh)
   log.Fatal(http.ListenAndServe(":8080",nil))
}