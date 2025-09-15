package service

import (
	"context"
	"fmt"
	"wallet-service/internal/repo"
)

type WalletService struct {
	repo repo.WalletRepo
}

func NewWalletService(r repo.WalletRepo) *WalletService {
	return &WalletService{repo: r}
}

func (s *WalletService) GetBalance(ctx context.Context, id string) (float64, error) {
	return s.repo.GetBalance(ctx, id)
}

func (s *WalletService) Deposit(ctx context.Context, walletId string, amount float64) error {
	return s.repo.UpdateBalance(ctx, walletId, amount)
}

func (s *WalletService) Withdraw(ctx context.Context, walletId string, amount float64) error {
	err := s.repo.UpdateBalance(ctx, walletId, -amount)
	if err != nil && err.Error() == "insufficient funds" {
		return fmt.Errorf("cannot withdraw: %w", err)
	}
	return err
}
