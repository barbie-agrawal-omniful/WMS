package main

import (
	"context"
	"fmt"
	"log"
	"time"
	appinit "wms/init"

	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/http"
)

func main() {

	// Initialize config
	err := config.Init(time.Second * 10)
	if err != nil {
		log.Panicf("Error while initialising config, err: %v", err)
		panic(err)
	}

	ctx, err := config.TODOContext()
	if err != nil {
		log.Panicf("Error while getting context from config, err: %v", err)
		panic(err)
	}

	// Initialize services
	appinit.Initialize(ctx)

	// Run HTTP Server
	runHttpServer(ctx)

	db := appinit.GetClient()
	fmt.Printf("%T", db)
}

func runHttpServer(ctx context.Context) {
	server := http.InitializeServer(":8081", 10*time.Second, 10*time.Second, 70*time.Second)

	if err := server.StartServer("OMS Service Started"); err != nil {
		log.Fatalf("Error starting the server:", err)
	}

}
