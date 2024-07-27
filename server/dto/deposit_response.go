package dto

type DepositResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    struct {
		DepositURL      string `json:"depositUrl"`
		MerchantOrderID string `json:"merchantOrderID"`
		OrderID         string `json:"orderID"`
	} `json:"data"`
}
