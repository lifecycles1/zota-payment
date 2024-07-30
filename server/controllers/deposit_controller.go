package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"zota_payment/dto"
	"zota_payment/services"

	"github.com/gorilla/mux"
)

type DepositController struct {
	depositService *services.DepositService
}

func NewDepositController(depositService *services.DepositService) *DepositController {
	return &DepositController{
		depositService: depositService,
	}
}

func (c *DepositController) DepositRequestHandler(w http.ResponseWriter, r *http.Request) {
	var request dto.DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	log.Printf("Received deposit request: %v", request)

	vars := mux.Vars(r)
	log.Printf("Vars: %v", vars)
	endpointID := vars["endpointID"]
	log.Printf("Endpoint ID: %v", endpointID)
	if endpointID == "" {
		http.Error(w, "Missing endpoint ID", http.StatusBadRequest)
		return
	}

	response, err := c.depositService.CreateDepositRequest(endpointID, request)
	if err != nil {
		if response != nil {
			// return the response code and message from the service's response
			code, convErr := strconv.Atoi(response.Code)
			if convErr != nil {
				log.Printf("Error converting response code: %v", convErr)
				http.Error(w, "Error converting response code", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(code)
			json.NewEncoder(w).Encode(map[string]string{"message": response.Message})
		} else {
			// return a generic error response
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
	}
}

func (c *DepositController) CallbackHandler(w http.ResponseWriter, r *http.Request) {
	var callbackRequest dto.CallbackNotification

	if err := json.NewDecoder(r.Body).Decode(&callbackRequest); err != nil {
		log.Printf("Error decoding request: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := c.depositService.HandleCallback(callbackRequest)
	if err != nil {
		log.Printf("Error handling callback: %v", err)
		http.Error(w, "Error handling callback", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
