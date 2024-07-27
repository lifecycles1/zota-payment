package main

import (
	"log"
	"net/http"
	"zota_payment/controllers"
	"zota_payment/routes"
	"zota_payment/services"

	"github.com/rs/cors"
)

func main() {
	baseURL := "https://api.zotapay-stage.com"
	merchantID := "BUGBOUNTY231"
	merchantSecretKey := "866adddb-7b91-4b1b-82a2-364479e17486"

	// services
	// 1) create deposit request (responds immediately with redirectUrl to client)
	depositService := services.NewDepositService(baseURL, merchantSecretKey)
	// 2) check order status (polling until final status) (merchantID is obtained through request params since a second call is being made to the service)
	orderStatusService := services.NewOrderStatusService(baseURL, merchantSecretKey)
	// optional deposit flow with both services combined into an automated direct call. POST request with endpointID specified in url and a req.body with at least the mandatory fields (without returning redirectUrl to client - only returns response once final status is reached) (merchantID is obtained through hardcoded constant here since we are automating the call to the order status service)
	depositFlowService := services.NewDepositFlowService(depositService, orderStatusService, merchantID)

	// controllers
	depositController := controllers.NewDepositController(depositFlowService, depositService)
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
