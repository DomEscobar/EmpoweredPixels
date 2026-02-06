package shop

import (
	"context"
)

// PaymentProvider defines the interface for external payment processors
type PaymentProvider interface {
	// ProcessPayment simulates or initiates a payment for a specific amount and currency
	// Returns a provider-specific transaction ID or an error
	ProcessPayment(ctx context.Context, userID int, amount int, currency string) (string, error)
}

// SimulatedPaymentProvider is a mock implementation for testing and development
type SimulatedPaymentProvider struct {
	ShouldFail bool
}

// NewSimulatedPaymentProvider creates a new simulated payment provider
func NewSimulatedPaymentProvider() *SimulatedPaymentProvider {
	return &SimulatedPaymentProvider{}
}

// ProcessPayment implements the PaymentProvider interface
func (p *SimulatedPaymentProvider) ProcessPayment(ctx context.Context, userID int, amount int, currency string) (string, error) {
	if p.ShouldFail {
		return "", context.DeadlineExceeded // Simulate a failure
	}
	// Return a dummy provider transaction ID
	return "sim_tx_12345", nil
}
