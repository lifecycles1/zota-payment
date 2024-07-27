package dto

type OrderStatusEnum string

const (
	Created    OrderStatusEnum = "CREATED"
	Processing OrderStatusEnum = "PROCESSING"
	Pending    OrderStatusEnum = "PENDING"
	Unknown    OrderStatusEnum = "UNKNOWN"
	Filtered   OrderStatusEnum = "FILTERED" // final statuses start from here
	Approved   OrderStatusEnum = "APPROVED"
	Declined   OrderStatusEnum = "DECLINED"
	Error      OrderStatusEnum = "ERROR"
)

type OrderStatusResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Type                   string          `json:"type"`
		Status                 OrderStatusEnum `json:"status"`
		ErrorMessage           string          `json:"errorMessage"`
		ProcessorTransactionID string          `json:"processorTransactionID"`
		MerchantOrderID        string          `json:"merchantOrderID"`
		OrderID                string          `json:"orderID"`
		Amount                 string          `json:"amount"`
		Currency               string          `json:"currency"`
		CustomerEmail          string          `json:"customerEmail"`
		CustomParam            string          `json:"customParam"`
		ExtraData              interface{}     `json:"extraData"`
		Request                interface{}     `json:"request"`
	} `json:"data"`
}
