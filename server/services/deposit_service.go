package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"zota_payment/dto"
	"zota_payment/kafka"
	"zota_payment/utils"
)

type DepositService struct {
	BaseURL           string
	MerchantSecretKey string
	kafkaProducer     *kafka.Producer
}

func NewDepositService(baseURL, merchantSecretKey string, kafkaProducer *kafka.Producer) *DepositService {
	return &DepositService{
		BaseURL:           baseURL,
		MerchantSecretKey: merchantSecretKey,
		kafkaProducer:     kafkaProducer,
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
		return nil, fmt.Errorf("error marshalling request: %w", err)
	}

	log.Printf("body: %s", body)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		log.Printf("Error creating request: %v", err)
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error making request: %v", err)
		return nil, fmt.Errorf("error making request: %w", err)
	}

	defer resp.Body.Close()

	var depositResponse dto.DepositResponse
	if err := json.NewDecoder(resp.Body).Decode(&depositResponse); err != nil {
		log.Printf("Error decoding response: %v", err)
		return nil, fmt.Errorf("error decoding response: %w", err)
	}

	log.Printf("depositResponse: %+v", depositResponse)
	log.Printf("resp.StatusCode: %d", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return &depositResponse, fmt.Errorf("failed to create deposit request: %s", resp.Status)
	}

	return &depositResponse, nil
}

func (s *DepositService) HandleCallback(request dto.CallbackNotification) error {
	log.Printf("Received callback request: %+v", request)

	// validate the signature
	signatureString := fmt.Sprintf("%s%s%s%s%s%s%s",
		request.EndpointID, request.OrderID, request.MerchantOrderID, request.Status, request.Amount, request.CustomerEmail, s.MerchantSecretKey)
	signature := utils.GenerateSignature(signatureString)

	if signature != request.Signature {
		log.Printf("Invalid signature")
		return fmt.Errorf("unable to validate callback signature")
	}

	log.Printf("Callback signature is valid")

	// send the callback notification to the Kafka topic
	message, err := json.Marshal(request)
	if err != nil {
		log.Printf("Error marshalling callback request: %v", err)
		return fmt.Errorf("error marshalling callback request: %w", err)
	}

	err = s.kafkaProducer.SendMessage("callback_notifications", string(message))
	if err != nil {
		log.Printf("Error sending message to Kafka: %v", err)
		return fmt.Errorf("error sending message to Kafka: %w", err)
	}

	log.Printf("Callback notification sent to Kafka")

	return nil
}
