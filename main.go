package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"

	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var addr = "127.0.0.1:3000"

// Config - loaded from config.yml
var Config = struct {
	Database struct {
		Driver     string `default:"sqlite3"`
		Connection string `default:"tmp/dev.db"`
	}
}{}

func main() {
	l := log.New(os.Stdout, "[cryptogote] ", 0)

	configor.Load(&Config, "config.yml")

	db, err := gorm.Open(Config.Database.Driver, Config.Database.Connection)
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
