package payment

import (
	"fmt"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/midtrans/midtrans-go/snap"
)

type MidtransClient struct {
	snapClient snap.Client
	coreClient coreapi.Client
}

type CreateTransactionRequest struct {
	OrderID        string
	GrossAmount    int64
	CustomerName   string
	CustomerEmail  string
	CustomerPhone  string
	ItemName       string
	ItemPrice      int64
	ItemQuantity   int32
}

type CreateTransactionResponse struct {
	Token       string
	RedirectURL string
}

type TransactionStatus struct {
	OrderID           string
	TransactionID     string
	TransactionStatus string
	PaymentType       string
	GrossAmount       string
	TransactionTime   string
	FraudStatus       string
}

func NewMidtransClient(serverKey string, isProduction bool) *MidtransClient {
	var env midtrans.EnvironmentType
	if isProduction {
		env = midtrans.Production
	} else {
		env = midtrans.Sandbox
	}

	var s snap.Client
	s.New(serverKey, env)

	var c coreapi.Client
	c.New(serverKey, env)

	return &MidtransClient{
		snapClient: s,
		coreClient: c,
	}
}

func (m *MidtransClient) CreateTransaction(req CreateTransactionRequest) (*CreateTransactionResponse, error) {
	// Build Snap request
	snapReq := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  req.OrderID,
			GrossAmt: req.GrossAmount,
		},
		CustomerDetail: &midtrans.CustomerDetails{
			FName: req.CustomerName,
			Email: req.CustomerEmail,
			Phone: req.CustomerPhone,
		},
		Items: &[]midtrans.ItemDetails{
			{
				ID:    req.OrderID,
				Name:  req.ItemName,
				Price: req.ItemPrice,
				Qty:   req.ItemQuantity,
			},
		},
		EnabledPayments: snap.AllSnapPaymentType,
	}

	// Create Snap transaction
	snapResp, err := m.snapClient.CreateTransaction(snapReq)
	if err != nil {
		return nil, fmt.Errorf("failed to create Midtrans transaction: %w", err)
	}

	return &CreateTransactionResponse{
		Token:       snapResp.Token,
		RedirectURL: snapResp.RedirectURL,
	}, nil
}

func (m *MidtransClient) GetStatus(orderID string) (*TransactionStatus, error) {
	// Use core API to check transaction status
	statusResp, err := m.coreClient.CheckTransaction(orderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction status: %w", err)
	}

	return &TransactionStatus{
		OrderID:           statusResp.OrderID,
		TransactionID:     statusResp.TransactionID,
		TransactionStatus: statusResp.TransactionStatus,
		PaymentType:       statusResp.PaymentType,
		GrossAmount:       statusResp.GrossAmount,
		TransactionTime:   statusResp.TransactionTime,
		FraudStatus:       statusResp.FraudStatus,
	}, nil
}

type WebhookNotification struct {
	TransactionTime   string `json:"transaction_time"`
	TransactionStatus string `json:"transaction_status"`
	TransactionID     string `json:"transaction_id"`
	StatusMessage     string `json:"status_message"`
	StatusCode        string `json:"status_code"`
	SignatureKey      string `json:"signature_key"`
	PaymentType       string `json:"payment_type"`
	OrderID           string `json:"order_id"`
	MerchantID        string `json:"merchant_id"`
	GrossAmount       string `json:"gross_amount"`
	FraudStatus       string `json:"fraud_status"`
	Currency          string `json:"currency"`
}

func (m *MidtransClient) HandleWebhook(notification *WebhookNotification) (*TransactionStatus, error) {
	// Verify notification (in production, verify signature)
	// For now, just return the status
	return &TransactionStatus{
		OrderID:           notification.OrderID,
		TransactionID:     notification.TransactionID,
		TransactionStatus: notification.TransactionStatus,
		PaymentType:       notification.PaymentType,
		GrossAmount:       notification.GrossAmount,
		TransactionTime:   notification.TransactionTime,
		FraudStatus:       notification.FraudStatus,
	}, nil
}

// Helper function to determine if payment is successful
func IsPaymentSuccess(status string) bool {
	return status == "capture" || status == "settlement"
}

// Helper function to determine if payment is pending
func IsPaymentPending(status string) bool {
	return status == "pending"
}

// Helper function to determine if payment failed
func IsPaymentFailed(status string) bool {
	return status == "deny" || status == "cancel" || status == "expire" || status == "failure"
}
