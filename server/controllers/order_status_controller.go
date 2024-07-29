package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"zota_payment/dto"
	"zota_payment/services"
)

type OrderStatusController struct {
	OrderStatusService *services.OrderStatusService
}

func NewOrderStatusController(orderStatusService *services.OrderStatusService) *OrderStatusController {
	return &OrderStatusController{
		OrderStatusService: orderStatusService,
	}
}

// GET /api/v1/query/order-status/
func (c *OrderStatusController) GetOrderStatusHandler(w http.ResponseWriter, r *http.Request) {

	queryParams := r.URL.Query()
	merchantID := queryParams.Get("merchantID")
	orderID := queryParams.Get("orderID")
	merchantOrderID := queryParams.Get("merchantOrderID")
	timestamp := queryParams.Get("timestamp")
	signature := queryParams.Get("signature")

	if merchantID == "" || orderID == "" || merchantOrderID == "" || timestamp == "" || signature == "" {
		http.Error(w, "Missing required query parameters", http.StatusBadRequest)
		return
	}

	params := dto.OrderStatusRequest{
		MerchantID:      merchantID,
		OrderID:         orderID,
		MerchantOrderID: merchantOrderID,
		Timestamp:       timestamp,
		Signature:       signature,
	}

	response, err := c.OrderStatusService.GetOrderStatus(params)
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
