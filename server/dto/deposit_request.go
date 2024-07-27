package dto

type DepositRequest struct {
	MerchantOrderID           string `json:"merchantOrderID"`
	MerchantOrderDesc         string `json:"merchantOrderDesc"`
	OrderAmount               string `json:"orderAmount"`
	OrderCurrency             string `json:"orderCurrency"`
	CustomerEmail             string `json:"customerEmail"`
	CustomerFirstName         string `json:"customerFirstName"`
	CustomerLastName          string `json:"customerLastName"`
	CustomerAddress           string `json:"customerAddress"`
	CustomerCountryCode       string `json:"customerCountryCode"`
	CustomerCity              string `json:"customerCity"`
	CustomerState             string `json:"customerState,omitempty"`
	CustomerZipCode           string `json:"customerZipCode"`
	CustomerPhone             string `json:"customerPhone"`
	CustomerIP                string `json:"customerIP"`
	CustomerPersonalID        string `json:"customerPersonalID,omitempty"`
	CustomerBankCode          string `json:"customerBankCode,omitempty"`
	CustomerBankAccountNumber string `json:"customerBankAccountNumber,omitempty"`
	RedirectURL               string `json:"redirectUrl"`
	CallbackURL               string `json:"callbackUrl,omitempty"`
	CheckoutURL               string `json:"checkoutUrl"`
	CustomParam               string `json:"customParam,omitempty"`
	Language                  string `json:"language,omitempty"`
	Signature                 string `json:"signature"`
}
