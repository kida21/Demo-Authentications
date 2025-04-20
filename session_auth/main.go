package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))
var users = map[string]string{"thomas":"12345","admin":"password"} //hardcoded for testin purpose

func healthCheckHandler(w http.ResponseWriter,r *http.Request){
	session,_:=store.Get(r,"session.id")
	if (session.Values["authenticated"] != nil) && (session.Values["authenticated"] != false){
		w.Write([]byte(time.Now().GoString()))
	}else{
		http.Error(w,"forbidden",http.StatusForbidden)
		return
	}
}

func loginHandler(w http.ResponseWriter,r *http.Request){
	session,err:= store.Get(r,"session.id")
    if err!=nil{
		w.Write([]byte("error..."))
	}
	err=r.ParseForm()
	if err!=nil{
		http.Error(w,"please pass the data as url parameter",http.StatusBadRequest)
	}
	username := r.PostForm.Get("username")
	password := r.PostForm.Get("password")
  
	if orginalpassword,ok := users[username];ok{
        if password == orginalpassword{
		session.Values["authenticated"] = true
		session.Save(r,w)
	   }else{
		http.Error(w,"Invalid credentials",http.StatusUnauthorized)
	   }
	 }else{
		http.Error(w,"user not found",http.StatusNotFound)
		return
	 }
	 w.Write([]byte("login successfully......"))
}

func main() {
  r:=mux.NewRouter()
  r.HandleFunc("/login",loginHandler)
  r.HandleFunc("/healthcheck",healthCheckHandler)
  srv:=&http.Server{
	Addr: ":8080",
	Handler: r,
	WriteTimeout: 10 * time.Second,
	ReadTimeout: 15 * time.Second,
  }
  log.Println("server started listening.....")
  log.Print(srv.ListenAndServe())
}