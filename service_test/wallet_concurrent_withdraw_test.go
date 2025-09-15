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

func setupWalletServiceForWithdraw(t *testing.T) *service.WalletService {
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

func TestConcurrentWithdraw(t *testing.T) {
	svc := setupWalletServiceForWithdraw(t)

	const n = 10
	walletID := "22222222-2222-2222-2222-222222222222"

	err := svc.Deposit(context.Background(), walletID, 500)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(n)

	errorsCh := make(chan error, n)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			err := svc.Withdraw(context.Background(), walletID, 100)
			errorsCh <- err
		}()
	}

	wg.Wait()
	close(errorsCh)

	insufficientFundsCount := 0
	for err := range errorsCh {
		if err != nil {
			if err.Error() == "cannot withdraw: insufficient funds" {
				insufficientFundsCount++
			} else {
				t.Errorf("unexpected error: %v", err)
			}
		}
	}

	if insufficientFundsCount != n-5 {
		t.Errorf("expected %d insufficient funds errors, got %d", n-5, insufficientFundsCount)
	}

	balance, err := svc.GetBalance(context.Background(), walletID)
	if err != nil {
		t.Fatal(err)
	}

	expectedBalance := 0.0
	if balance != expectedBalance {
		t.Errorf("expected balance %f, got %f", expectedBalance, balance)
	}
}
