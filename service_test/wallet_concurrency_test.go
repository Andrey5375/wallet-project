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

func setupWalletService(t *testing.T) *service.WalletService {
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

func TestConcurrentDeposit(t *testing.T) {
	svc := setupWalletService(t)

	const n = 100
	walletID := "11111111-1111-1111-1111-111111111111"

	err := svc.Deposit(context.Background(), walletID, 0)
	if err != nil {
		t.Fatal(err)
	}

	var wg sync.WaitGroup
	wg.Add(n)

	sem := make(chan struct{}, 20)

	for i := 0; i < n; i++ {
		go func() {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			err := svc.Deposit(context.Background(), walletID, 100)
			if err != nil {
				t.Error(err)
			}
		}()
	}

	wg.Wait()

	balance, err := svc.GetBalance(context.Background(), walletID)
	if err != nil {
		t.Fatal(err)
	}

	expected := float64(n * 100)
	if balance != expected {
		t.Errorf("expected balance %f, got %f", expected, balance)
	}
}
