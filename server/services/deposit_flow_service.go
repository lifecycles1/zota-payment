package services

import (
	"fmt"
	"time"
	"zota_payment/dto"
)

type DepositFlowService struct {
	DepositService     *DepositService
	OrderStatusService *OrderStatusService
	MerchantID         string
}

func NewDepositFlowService(depositService *DepositService, orderStatusService *OrderStatusService, merchantID string) *DepositFlowService {
	return &DepositFlowService{
		DepositService:     depositService,
		OrderStatusService: orderStatusService,
		MerchantID:         merchantID,
	}
}

// backend flow with backend polling for final order status
func (s *DepositFlowService) ExecuteDepositFlow(endpointID string, request dto.DepositRequest) (*dto.OrderStatusResponse, error) {

	depositResponse, err := s.DepositService.CreateDepositRequest(endpointID, request)
	if err != nil {
		if depositResponse != nil {
			orderStatusResponse := dto.OrderStatusResponse{
				Code:    depositResponse.Code,
				Message: depositResponse.Message,
			}
			return &orderStatusResponse, err
		}
		return nil, err
	}

	params := dto.OrderStatusRequest{
		MerchantID:      s.MerchantID,
		MerchantOrderID: depositResponse.Data.MerchantOrderID,
		OrderID:         depositResponse.Data.OrderID,
		Timestamp:       fmt.Sprintf("%d", time.Now().Unix()),
	}

	// start polling for final order status
	finalStatusResponse, err := s.OrderStatusService.PollOrderStatus(params)
	if err != nil {
		if finalStatusResponse != nil {
			return finalStatusResponse, err
		}
		return nil, err
	}

	return finalStatusResponse, nil
}
