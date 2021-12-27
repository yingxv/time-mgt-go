package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/NgeKaworu/time-mgt-go/src/app"
	"github.com/NgeKaworu/time-mgt-go/src/db"
	"github.com/go-redis/redis/v8"
	"github.com/julienschmidt/httprouter"
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		addr   = flag.String("l", ":8050", "绑定Host地址")
		dbinit = flag.Bool("i", false, "init database flag")
		mongo  = flag.String("m", "mongodb://localhost:27017", "mongod addr flag")
		mdb    = flag.String("db", "time-mgt", "database name")
		ucHost = flag.String("uc", "https://api.furan.xyz/user-center", "user center host")
		r      = flag.String("r", "localhost:6379", "rdb addr")
	)
	flag.Parse()

	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)

	mongoClient := db.NewMongoClient()
	err := mongoClient.Open(*mongo, *mdb, *dbinit)
	if err != nil {
		panic(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     *r,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	app := app.New(ucHost, mongoClient, rdb)
	if err != nil {
		panic(err)
	}

	router := httprouter.New()
	// tag ctrl
	router.POST("/v1/tag/create", app.AddTag)
	router.PUT("/v1/tag/update", app.SetTag)
	router.GET("/v1/tag/list", app.ListTag)
	router.DELETE("/v1/tag/:id", app.RemoveTag)
	//record ctrl
	router.POST("/v1/record/create", app.AddRecord)
	router.PUT("/v1/record/update", app.SetRecord)
	router.GET("/v1/record/list", app.ListRecord)
	router.DELETE("/v1/record/:id", app.RemoveRecord)
	router.POST("/v1/record/statistic", app.StatisticRecord)

	srv := &http.Server{Handler: app.IsLogin(router), ErrorLog: nil}
	srv.Addr = *addr

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	log.Println("server on http port", srv.Addr)

	signalChan := make(chan os.Signal, 1)
	cleanupDone := make(chan bool)
	cleanup := make(chan bool)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		for range signalChan {
			ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
			defer cancel()
			go func() {
				_ = srv.Shutdown(ctx)
				cleanup <- true
			}()
			<-cleanup
			mongoClient.Close()
			rdb.Close()
			fmt.Println("safe exit")
			cleanupDone <- true
		}
	}()
	<-cleanupDone

}
