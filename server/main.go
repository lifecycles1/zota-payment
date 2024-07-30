package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"zota_payment/controllers"
	"zota_payment/kafka"
	"zota_payment/postgres"
	"zota_payment/repositories"
	"zota_payment/routes"
	"zota_payment/services"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	baseURL := "https://api.zotapay-stage.com"
	merchantSecretKey := os.Getenv("MERCHANT_SECRET_KEY")
	log.Println("merchantSecretKey: ", merchantSecretKey)
	callbackTopic := "callback_notifications"

	// initialize database
	db, err := postgres.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// create kafka topic if it doesn't exist
	kafka.CreateKafkaTopic(callbackTopic, 1, 1)

	// initialize kafka producer
	kafkaProducer, err := kafka.NewProducer()
	if err != nil {
		log.Fatalf("Failed to initialize Kafka producer: %v", err)
	}
	defer kafkaProducer.Close()

	// initialize repositories
	callbackRepo := repositories.NewCallbackRepository(db)

	// initialize kafka consumer
	go kafka.StartKafkaConsumer(callbackTopic, callbackRepo)

	// initialize services
	// 1) create deposit request (responds immediately with redirectUrl to client)
	depositService := services.NewDepositService(baseURL, merchantSecretKey, kafkaProducer)
	// 2) check order status (polling until final status)
	orderStatusService := services.NewOrderStatusService(baseURL, merchantSecretKey)

	// initialize controllers
	depositController := controllers.NewDepositController(depositService)
	orderStatusController := controllers.NewOrderStatusController(orderStatusService)

	router := routes.SetupRoutes(depositController, orderStatusController)

	// Enable CORS
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // allow all origins for now
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}).Handler(router)

	// log.Println("Server starting on port 8080...")
	// log.Fatal(http.ListenAndServe(":8080", corsHandler))

	// start the server
	server := &http.Server{
		Addr:    ":8080",
		Handler: corsHandler,
	}

	go func() {
		log.Println("Server starting on port 8080...")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	gracefulShutdown(server)
}

func gracefulShutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Failed to gracefully shutdown server: %v", err)
	}

	log.Println("Server shutdown")
}
