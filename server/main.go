package main

import (
	"log"
	"net/http"
	"os"
	"zota_payment/controllers"
	"zota_payment/routes"
	"zota_payment/services"

	"github.com/joho/godotenv"
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
	// services
	// 1) create deposit request (responds immediately with redirectUrl to client)
	depositService := services.NewDepositService(baseURL, merchantSecretKey)
	// 2) check order status (polling until final status)
	orderStatusService := services.NewOrderStatusService(baseURL, merchantSecretKey)

	// controllers
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

	log.Println("Server starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", corsHandler))
}
