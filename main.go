package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
	"hacktiv8-golang-assignment-final/api/controllers"
	"hacktiv8-golang-assignment-final/socket"
	"hacktiv8-golang-assignment-final/utils"
	"log"
	"net/http"
	"os"
)

// address flag with default :8080
var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	// Init Logrus
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.DebugLevel)

	// Set Static Assets
	http.Handle("/assets/",
		http.StripPrefix("/assets/",
			http.FileServer(http.Dir("static/assets"))))

	// Set Static Routes
	// http.HandleFunc("/", handler.Home)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/views/index.html")
	})
	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/views/register.html")
	})
	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/views/chat.html")
	})

	// Init Engine Session
	var store = sessions.NewCookieStore([]byte(utils.SessionKey))

	authController := controllers.NewAuthController(store, logger)

	http.HandleFunc("/api/check_session", authController.CheckSession)
	http.HandleFunc("/api/login", authController.Login)
	http.HandleFunc("/api/logout", authController.Logout)
	http.HandleFunc("/api/register", authController.Register)

	hub := socket.NewHub()
	go hub.Run()

	// Set Web Socket Routes
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, store, w, r)
	})

	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	} else {
		fmt.Println()
	}
}
