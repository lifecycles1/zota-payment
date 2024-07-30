package dto

type CallbackNotification struct {
	Type                   string                 `json:"type"`
	Status                 OrderStatusEnum        `json:"status"`
	ErrorMessage           string                 `json:"errorMessage"`
	EndpointID             string                 `json:"endpointID"`
	ProcessorTransactionID string                 `json:"processorTransactionID"`
	OrderID                string                 `json:"orderID"`
	MerchantOrderID        string                 `json:"merchantOrderID"`
	Amount                 string                 `json:"amount"`
	Currency               string                 `json:"currency"`
	CustomerEmail          string                 `json:"customerEmail"`
	CustomParam            string                 `json:"customParam"`
	ExtraData              map[string]interface{} `json:"extraData"`
	OriginalRequest        map[string]interface{} `json:"originalRequest"`
	Signature              string                 `json:"signature"`
}
