package usecase

import (
	"context"
)

type TransactionVerificationService interface {
	CancelByID(ctx context.Context, id string) error
	CompletedByID(ctx context.Context, id string) error
}

type VerifierUseCase struct {
	ts TransactionVerificationService
}

func NewVerifierUseCase(ts TransactionVerificationService) *VerifierUseCase {
	return &VerifierUseCase{ts: ts}
}

func (v *VerifierUseCase) CancelTransactionByID(ctx context.Context, id string) error {
	return v.ts.CancelByID(ctx, id)
}

func (v *VerifierUseCase) CompleteTransactionByID(ctx context.Context, id string) error {
	return v.ts.CompletedByID(ctx, id)
}
