package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/7anekaha/go-restful-api-week-17/services"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	var addr string
	flag.StringVar(&addr, "addr", ":3000", "HTTP network address")
	flag.Parse()

	app := &App{
		mux:             http.NewServeMux(),
		Addr:            addr,
		mongoService:    services.NewMongoService(context.Background()),
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

func init() {
	log.Println("Populate mongodb in order to test the API")

	var dbName string = os.Getenv("DB")
	var collectionName string = os.Getenv("COLLECTION")
	var uri string = os.Getenv("MONGO_URI")

	ctx := context.Background()
	mongoClient, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatalln(err)
	}

	// Check the connection
	if err := mongoClient.Ping(ctx, nil); err != nil {
		log.Fatalln(err)
	}

	createdAtStr1 := "2017-01-27T08:19:14.135Z"
	createdAtStr2 := "2017-01-27T13:22:10.421Z"
	createdAtStr3 := "2017-01-28T01:22:14.398Z"
	createdAt1, _ := time.Parse(time.RFC3339, createdAtStr1)
	createdAt2, _ := time.Parse(time.RFC3339, createdAtStr2)
	createdAt3, _ := time.Parse(time.RFC3339, createdAtStr3)

	records := []interface{}{
		services.MongoRecord{
			Key:       "TAKwGc6Jr4i8Z487",
			CreatedAt: createdAt1,
			Count:     []int{500, 400, 450, 550, 300, 150, 350},
		},
		services.MongoRecord{
			Key:       "NAeQ8eX7e5TEg70H",
			CreatedAt: createdAt2,
			Count:     []int{540, 400, 450, 550, 300, 160, 350},
		},
		services.MongoRecord{
			Key:       "cCddT2RPqWmUI4Nf",
			CreatedAt: createdAt3,
			Count:     []int{120, 400, 450, 660, 500, 770, 250},
		},
	}

	res, err := mongoClient.Database(dbName).Collection(collectionName).InsertMany(ctx, records)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(res.InsertedIDs)

}
