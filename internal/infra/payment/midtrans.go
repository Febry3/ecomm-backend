package payment

import (
	"context"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/sirupsen/logrus"
)

type MidtransGateway struct {
	client    coreapi.Client
	serverKey string
	log       *logrus.Logger
}

func NewMidtransGateway(client coreapi.Client, serverKey string, log *logrus.Logger) PaymentGateway {
	return &MidtransGateway{
		client:    client,
		serverKey: serverKey,
		log:       log,
	}
}

// ChargeVA creates a Virtual Account payment request
// expiresAt is optional - if nil, uses default 5 minute expiry
func (m *MidtransGateway) ChargeVA(ctx context.Context, orderID string, amount int64, bankCode string, expiresAt *time.Time) (*VAPaymentResult, error) {
	var req *coreapi.ChargeReq

	switch bankCode {
	case "bca", "bni", "bri", "cimb", "permata":
		req = &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeBankTransfer,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: amount,
			},
			BankTransfer: &coreapi.BankTransferDetails{
				Bank: midtrans.Bank(bankCode),
			},
		}
	case "mandiri":
		// Mandiri uses echannel (Bill Payment)
		req = &coreapi.ChargeReq{
			PaymentType: coreapi.PaymentTypeEChannel,
			TransactionDetails: midtrans.TransactionDetails{
				OrderID:  orderID,
				GrossAmt: amount,
			},
			EChannel: &coreapi.EChannelDetail{
				BillInfo1: "Payment:",
				BillInfo2: "Online Purchase",
			},
		}
	default:
		return nil, fmt.Errorf("unsupported bank code: %s", bankCode)
	}

	// Set custom expiry based on expiresAt or default to 5 minutes
	if expiresAt != nil {
		duration := time.Until(*expiresAt)
		if duration < time.Minute {
			duration = time.Minute // minimum 1 minute
		}
		expiryMinutes := int(duration.Minutes())
		req.CustomExpiry = &coreapi.CustomExpiry{
			ExpiryDuration: expiryMinutes,
			Unit:           "minute",
		}
	} else {
		req.CustomExpiry = &coreapi.CustomExpiry{
			ExpiryDuration: 5,
			Unit:           "minute",
		}
	}

	resp, err := m.client.ChargeTransaction(req)
	if err != nil {
		m.log.Errorf("Midtrans ChargeTransaction error: %v", err)
		return nil, fmt.Errorf("failed to create VA payment: %w", err)
	}

	if resp.StatusCode != "201" && resp.StatusCode != "200" {
		m.log.Errorf("Midtrans ChargeTransaction failed: %s - %s", resp.StatusCode, resp.StatusMessage)
		return nil, fmt.Errorf("midtrans error: %s", resp.StatusMessage)
	}

	result := &VAPaymentResult{
		TransactionID: resp.TransactionID,
		OrderID:       resp.OrderID,
		GrossAmount:   parseAmount(resp.GrossAmount),
		Status:        resp.TransactionStatus,
		ExpiredAt:     m.calculateExpiry(resp.TransactionTime),
	}

	// Handle different VA response formats
	if len(resp.VaNumbers) > 0 {
		result.Bank = resp.VaNumbers[0].Bank
		result.VANumber = resp.VaNumbers[0].VANumber
	} else if resp.PermataVaNumber != "" {
		result.Bank = "permata"
		result.VANumber = resp.PermataVaNumber
	} else if resp.BillKey != "" {
		// Mandiri Bill Payment
		result.Bank = "mandiri"
		result.BillKey = resp.BillKey
		result.BillerCode = resp.BillerCode
	}

	m.log.Infof("VA payment created: Order=%s, Bank=%s, VA=%s", orderID, result.Bank, result.VANumber)
	return result, nil
}

// GetTransactionStatus checks the current status of a transaction
func (m *MidtransGateway) GetTransactionStatus(ctx context.Context, orderID string) (*PaymentStatusResult, error) {
	resp, err := m.client.CheckTransaction(orderID)
	if err != nil {
		m.log.Errorf("Midtrans CheckTransaction error for order %s: %v", orderID, err)
		return nil, fmt.Errorf("failed to check transaction status: %w", err)
	}

	result := &PaymentStatusResult{
		TransactionID: resp.TransactionID,
		OrderID:       resp.OrderID,
		Status:        resp.TransactionStatus,
		PaymentType:   resp.PaymentType,
		GrossAmount:   parseAmount(resp.GrossAmount),
	}

	// Parse settlement time if available
	if resp.SettlementTime != "" {
		paidAt, err := time.Parse("2006-01-02 15:04:05", resp.SettlementTime)
		if err == nil {
			result.PaidAt = &paidAt
		}
	}

	return result, nil
}

// VerifySignature validates the webhook notification signature
// Signature = SHA512(order_id + status_code + gross_amount + server_key)
func (m *MidtransGateway) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	rawSignature := orderID + statusCode + grossAmount + m.serverKey
	hash := sha512.Sum512([]byte(rawSignature))
	expectedSignature := hex.EncodeToString(hash[:])

	return signatureKey == expectedSignature
}

// CancelTransaction cancels a pending transaction
func (m *MidtransGateway) CancelTransaction(ctx context.Context, orderID string) error {
	resp, err := m.client.CancelTransaction(orderID)
	if err != nil {
		m.log.Errorf("Midtrans CancelTransaction error for order %s: %v", orderID, err)
		return fmt.Errorf("failed to cancel transaction: %w", err)
	}

	if resp.StatusCode != "200" && resp.StatusCode != "201" && resp.StatusCode != "407" {
		return fmt.Errorf("cancel failed: %s", resp.StatusMessage)
	}

	m.log.Infof("Transaction cancelled: Order=%s", orderID)
	return nil
}

// Helper functions

func (m *MidtransGateway) calculateExpiry(transactionTime string) time.Time {
	// Parse Midtrans transaction time format: "2019-10-23 16:33:49"
	parsed, err := time.Parse("2006-01-02 15:04:05", transactionTime)
	if err != nil {
		m.log.Warnf("Failed to parse transaction time: %v, using current time + 5min", err)
		return time.Now().Add(5 * time.Minute)
	}
	// Add 5 minutes for expiry
	return parsed.Add(5 * time.Minute)
}

func parseAmount(amountStr string) float64 {
	var amount float64
	fmt.Sscanf(amountStr, "%f", &amount)
	return amount
}
