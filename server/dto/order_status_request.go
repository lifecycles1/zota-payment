package dto

type OrderStatusRequest struct {
	MerchantID      string `json:"merchantID"`
	MerchantOrderID string `json:"merchantOrderID"`
	OrderID         string `json:"orderID"`
	Timestamp       string `json:"timestamp"`
	Signature       string `json:"signature"`
}
