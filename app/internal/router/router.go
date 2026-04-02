package router

import (
	"context"
	"log"
	"net/http"
	"time"

	"app/config"
	"app/internal/api"
	"app/internal/blockchain"
	"app/internal/repository"
	"app/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func Initialize() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not loaded")
	}
	applicationConfig, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}
	contextWithTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	postgresRepository, err := repository.NewPostgresRepository(applicationConfig.PostgresDSN)
	if err != nil {
		log.Fatalf("init repository: %v", err)
	}
	defer func() {
		if closeError := postgresRepository.Close(); closeError != nil {
			log.Printf("close repository: %v", closeError)
		}
	}()
	blockchainClient, err := blockchain.NewClient(
		contextWithTimeout,
		applicationConfig.RpcURL,
		applicationConfig.ContractAddress,
		applicationConfig.PrivateKey,
	)
	if err != nil {
		log.Fatalf("init blockchain client: %v", err)
	}
	defer blockchainClient.Close()

	stateService := service.NewService(blockchainClient, postgresRepository)
	handler := api.NewHandler(stateService)
	router := gin.Default()
	routesHandler := Handler{service: stateService}
	routesHandler.initializeRoutes(router, handler)

	server := &http.Server{
		Addr:         ":" + applicationConfig.HttpPort,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	log.Printf("server listening on :%s", applicationConfig.HttpPort)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

}
