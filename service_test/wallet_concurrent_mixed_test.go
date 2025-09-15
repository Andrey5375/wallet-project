package service_test

import (
	"context"
	"sync"
	"testing"
	"wallet-service/internal/repo"
	"wallet-service/internal/service"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func setupWalletServiceMixed(t *testing.T) *service.WalletService {
	t.Helper()

	db, err := sqlx.Connect("postgres", "host=localhost port=5433 user=user password=password dbname=wallet_test sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.Exec("TRUNCATE TABLE wallets")
	if err != nil {
		t.Fatal(err)
	}

	walletRepo := repo.NewWalletRepo(db)
	return service.NewWalletService(walletRepo)
}

func TestConcurrentMixed(t *testing.T) {
	svc := setupWalletServiceMixed(t)

	const nDeposits = 50
	const nWithdraws = 30
	walletID := "33333333-3333-3333-3333-333333333333"

	err := svc.Deposit(context.Background(), walletID, 1000)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(nDeposits + nWithdraws)

	sem := make(chan struct{}, 20)

	errorsCh := make(chan error, nDeposits+nWithdraws)

	for i := 0; i < nDeposits; i++ {
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := svc.Deposit(context.Background(), walletID, 50)
			errorsCh <- err
		}()
	}

	for i := 0; i < nWithdraws; i++ {
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := svc.Withdraw(context.Background(), walletID, 30)
			errorsCh <- err
		}()
	}

	wg.Wait()
	close(errorsCh)

	for err := range errorsCh {
		if err != nil && err.Error() != "cannot withdraw: insufficient funds" {
			t.Errorf("unexpected error: %v", err)
		}
	}

	balance, err := svc.GetBalance(context.Background(), walletID)
	if err != nil {
		t.Fatal(err)
	}

	expectedBalance := 1000 + float64(nDeposits*50) - float64(nWithdraws*30)
	if balance != expectedBalance {
		t.Errorf("expected balance %f, got %f", expectedBalance, balance)
	}
}
