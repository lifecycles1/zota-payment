package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zota_payment/dto"
	"zota_payment/utils"
)

type DepositService struct {
	BaseURL           string
	MerchantSecretKey string
}

func NewDepositService(baseURL, merchantSecretKey string) *DepositService {
	return &DepositService{
		BaseURL:           baseURL,
		MerchantSecretKey: merchantSecretKey,
	}
}

func (s *DepositService) CreateDepositRequest(endpointID string, request dto.DepositRequest) (*dto.DepositResponse, error) {
	// Generate the signature
	signatureString := fmt.Sprintf("%s%s%s%s%s",
		endpointID, request.MerchantOrderID, request.OrderAmount, request.CustomerEmail, s.MerchantSecretKey)
	request.Signature = utils.GenerateSignature(signatureString)
	log.Printf("request.Signature: %s", request.Signature)

	url := fmt.Sprintf("%s/api/v1/deposit/request/%s/", s.BaseURL, endpointID)
	log.Printf("url: %s", url)

	body, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling request: %v", err)
		return nil, err
	}

	log.Printf("body: %s", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, err
	}

	defer resp.Body.Close()

	var depositResponse dto.DepositResponse
	if err := json.NewDecoder(resp.Body).Decode(&depositResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, err
	}

	log.Printf("depositResponse: %+v", depositResponse)
	log.Printf("resp.StatusCode: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return &depositResponse, fmt.Errorf("failed to create deposit request: %s", resp.Status)
	}

	return &depositResponse, nil
}
