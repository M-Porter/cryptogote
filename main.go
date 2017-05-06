package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

var addr = "127.0.0.1:3000"

func main() {
	l := log.New(os.Stdout, "[cryptogote] ", 0)

	db, err := gorm.Open("sqlite3", "tmp/dev.db")
	if err != nil {
		l.Fatal("Failed connecting to database!")
	}

	db.AutoMigrate(&Notes{})

	renderer := render.New(render.Options{
		IndentJSON: true,
		Extensions: []string{".tpl"},
		Layout:     "application_layout",
	})

	pages := Pages{renderer, db}

	router := mux.NewRouter()
	router.HandleFunc("/", pages.NewMessageHandler).Methods("GET")
	router.HandleFunc("/messages", pages.PostCryptoMessageHandler).Methods("POST")
	router.HandleFunc("/messages/{key}", pages.ShowMessageHandler).Methods("GET")
	router.HandleFunc("/statistics", pages.StatisticsHandler).Methods("POST")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets"))))

	n := negroni.Classic()
	n.UseHandler(router)

	errs := make(chan error, 2)
	go func() {
		l.Printf("Listening on http://%s", addr)
		errs <- http.ListenAndServe(
			addr,
			csrf.Protect(
				[]byte("7f4711e775a65a807d50a92997e5f479"),
				csrf.Secure(false),
			)(n),
		)
	}()

	l.Fatal(<-errs)
}
