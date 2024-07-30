package repositories

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"zota_payment/dto"
)

type CallbackRepository struct {
	db *sql.DB
}

func NewCallbackRepository(db *sql.DB) *CallbackRepository {
	return &CallbackRepository{db: db}
}

func (r *CallbackRepository) InsertMessage(message string) error {
	var callbackNotification dto.CallbackNotification

	// unmarshal the kafka message into a CallbackNotification struct
	err := json.Unmarshal([]byte(message), &callbackNotification)
	if err != nil {
		log.Printf("Error unmarshalling message: %v", err)
		return fmt.Errorf("error unmarshalling message: %w", err)
	}

	// marshal extraData and originalRequest to JSON so we can store them as JSONB in the database
	extraData, err := json.Marshal(callbackNotification.ExtraData)
	if err != nil {
		log.Printf("Error marshalling extraData: %v", err)
		return fmt.Errorf("error marshalling extraData: %w", err)
	}

	originalRequest, err := json.Marshal(callbackNotification.OriginalRequest)
	if err != nil {
		log.Printf("Error marshalling originalRequest: %v", err)
		return fmt.Errorf("error marshalling originalRequest: %w", err)
	}

	// insert the message into the callback_notifications table
	query := `
		INSERT INTO callback_notifications (type, status, error_message, endpoint_id, processor_transaction_id, order_id, merchant_order_id, amount, currency, customer_email, custom_param, extra_data, original_request, signature)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`

	_, err = r.db.Exec(query,
		callbackNotification.Type,
		callbackNotification.Status,
		callbackNotification.ErrorMessage,
		callbackNotification.EndpointID,
		callbackNotification.ProcessorTransactionID,
		callbackNotification.OrderID,
		callbackNotification.MerchantOrderID,
		callbackNotification.Amount,
		callbackNotification.Currency,
		callbackNotification.CustomerEmail,
		callbackNotification.CustomParam,
		extraData,
		originalRequest,
		callbackNotification.Signature,
	)
	if err != nil {
		log.Printf("Error inserting message into callback_notifications table: %v", err)
		return fmt.Errorf("error inserting message into callback_notifications table: %w", err)
	}

	log.Printf("Message inserted into callback_notifications table")

	return nil
}
