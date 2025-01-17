package routes

import (
	"net/http"
	"zota_payment/controllers"

	"github.com/gorilla/mux"
)

// SetupRoutes sets up the routes for the application.
func SetupRoutes(depositController *controllers.DepositController, orderStatusController *controllers.OrderStatusController) http.Handler {
	r := mux.NewRouter()

	// single deposit request route
	// 1. can test individually
	// 1.1 with postman
	// 1.2 TEST Button provided in frontend(TestEndpoints).
	r.HandleFunc("/api/v1/deposit/request/{endpointID}/", depositController.DepositRequestHandler).Methods("POST")

	// single order status request route (
	// 1. can test individually
	// 1.1 with postman
	// 1.2 TEST Button provided in frontend(TestEndpoints). Test after running a single deposit request in frontend(TestEndpoints)
	r.HandleFunc("/api/v1/query/order-status/", orderStatusController.GetOrderStatusHandler).Methods("GET")

	// frontend deposit flow (TEST Client Flow Button provided in frontend(TestFlows)) with redirects using frontend(PaymentReturn.jsx as redirectUrl) polling until final order status is received.
	r.HandleFunc("/api/v1/deposit/client-flow/{endpointID}/", depositController.DepositRequestHandler).Methods("POST")

	// callback notification route just validates the signature for now, logs the request and returns 200 OK
	r.HandleFunc("/api/v1/payment-callback/", depositController.CallbackHandler).Methods("POST")

	return r
}
