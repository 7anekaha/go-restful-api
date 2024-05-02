package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/7anekaha/go-restful-api-week-17/services"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":3000", "HTTP network address")
	flag.Parse()

	app := &App{
		mux:             http.NewServeMux(),
		Addr:            addr,
		mongoService:    services.NewMongoService(),
		inMemoryService: services.NewInMemoryService(),
	}

	app.Initialize()
	if err := app.Run(); err != nil {
		log.Fatalln(err)
	}
}

type App struct {
	mux             *http.ServeMux
	Addr            string
	mongoService    *services.MongoService
	inMemoryService *services.InMemoryService
}

func (app *App) Initialize() {
	app.mux.HandleFunc("POST /mongodb", app.mongoService.Fetch)
	app.mux.HandleFunc("POST /in-memory", app.inMemoryService.Store)
	app.mux.HandleFunc("GET /in-memory", app.inMemoryService.Fetch)
}

func (app *App) Run() error {
	log.Printf("Server is running on %s\n", app.Addr)
	return http.ListenAndServe(app.Addr, app.mux)
}
