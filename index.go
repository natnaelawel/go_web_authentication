package main

import (
	_ "github.com/go-sql-driver/mysql"
	"html/template"
	"net/http"
	"github.com/natnaelawel/authentication/handlers"
)
//var db *sql.DB

var Temp *template.Template
//var Store = sessions.NewCookieStore([]byte("secrete password"))


const(
	ServerPort = ":8080"
)

func main(){

	mux := http.NewServeMux()
	assetFileServer := http.FileServer(http.Dir("assets"))
	tempFileServer := http.FileServer(http.Dir("template"))
	mux.Handle("/template/", http.StripPrefix("/template/", tempFileServer))
	mux.Handle("/assets/", http.StripPrefix("/assets/", assetFileServer))
	mux.HandleFunc("/signup", handlers.SignupHandler)
	mux.HandleFunc("/signin", handlers.LoginHandler)
	mux.HandleFunc("/logout", handlers.LogoutHandler)

	_ = http.ListenAndServe(ServerPort, mux)
}


