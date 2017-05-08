package server

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/unrolled/render"
	"github.com/urfave/negroni"
)

// App struct
type App struct {
	Render *render.Render
	DB     *gorm.DB
	Logger *log.Logger
	Config *Config
	Engine *negroni.Negroni
}

// NewApp - Create new App
func NewApp() *App {
	config := LoadConfig()

	logger := log.New(os.Stdout, "[cryptogote] ", 0)

	renderer := render.New(render.Options{
		IndentJSON: true,
		Extensions: []string{".tpl"},
		Layout:     "application_layout",
	})

	db, err := gorm.Open(config.Database.Driver, config.Database.Connection)
	if err != nil {
		logger.Fatal("Failed connecting to database!")
	}

	db.AutoMigrate(&Notes{})

	engine := setupEngine()

	app := &App{
		renderer,
		db,
		logger,
		config,
		engine,
	}

	return app
}

func setupEngine() *negroni.Negroni {
	pages := NewPages()

	router := mux.NewRouter()
	router.HandleFunc("/", pages.NewMessageHandler).Methods("GET")
	router.HandleFunc("/messages", pages.PostCryptoMessageHandler).Methods("POST")
	router.HandleFunc("/messages/{key}", pages.ShowMessageHandler).Methods("GET")
	router.HandleFunc("/statistics", pages.StatisticsHandler).Methods("POST")
	router.PathPrefix("/assets/").Handler(http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets"))))

	n := negroni.Classic()
	n.UseHandler(router)

	return n
}

// Run ...
func (app *App) Run() {
	addr := strings.Join([]string{app.Config.App.Addr, ":", app.Config.App.Port}, "")

	errs := make(chan error, 2)
	go func() {
		app.Logger.Printf("Listening on http://%s", addr)
		errs <- http.ListenAndServe(
			addr,
			csrf.Protect(
				[]byte(app.Config.App.Secret),
				csrf.Secure(false),
			)(app.Engine),
		)
	}()

	app.Logger.Fatal(<-errs)
}
