package services

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zota_payment/dto"
	"zota_payment/utils"
)

type OrderStatusService struct {
	BaseURL           string
	MerchantSecretKey string
}

func NewOrderStatusService(baseURL, merchantSecretKey string) *OrderStatusService {
	return &OrderStatusService{
		BaseURL:           baseURL,
		MerchantSecretKey: merchantSecretKey,
	}
}

// GET /api/v1/query/order-status/
// Makes an Order Status request to Zota
func (s *OrderStatusService) GetOrderStatus(params dto.OrderStatusRequest) (*dto.OrderStatusResponse, error) {

	endpointURL := "/api/v1/query/order-status/"

	signatureString := fmt.Sprintf("%s%s%s%s%s",
		params.MerchantID, params.MerchantOrderID, params.OrderID, params.Timestamp, s.MerchantSecretKey)
	params.Signature = utils.GenerateSignature(signatureString)

	queryString := fmt.Sprintf("merchantID=%s&merchantOrderID=%s&orderID=%s&timestamp=%s&signature=%s",
		params.MerchantID, params.MerchantOrderID, params.OrderID, params.Timestamp, params.Signature)

	url := fmt.Sprintf("%s%s?%s", s.BaseURL, endpointURL, queryString)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to make GET request to %s: %w", url, err)
	}

	defer resp.Body.Close()

	var orderStatusResponse dto.OrderStatusResponse
	if err := json.NewDecoder(resp.Body).Decode(&orderStatusResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("Order Status Response: %+v", orderStatusResponse)
	log.Printf("resp.StatusCode: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return &orderStatusResponse, fmt.Errorf("failed to get order status: %s", resp.Status)
	}

	return &orderStatusResponse, nil
}
