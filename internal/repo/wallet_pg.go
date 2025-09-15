package repo

import (
	"context"
	"database/sql"
	"errors"
	"wallet-service/internal/models"

	"github.com/jmoiron/sqlx"
)

type WalletRepo interface {
	GetBalance(ctx context.Context, id string) (float64, error)
	UpdateBalance(ctx context.Context, id string, delta float64) error
}

type walletRepo struct {
	db *sqlx.DB
}

func NewWalletRepo(db *sqlx.DB) WalletRepo {
	return &walletRepo{db: db}
}

func (r *walletRepo) GetBalance(ctx context.Context, id string) (float64, error) {
	var w models.Wallet
	err := r.db.GetContext(ctx, &w, "SELECT id, balance FROM wallets WHERE id=$1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, errors.New("wallet not found")
		}
		return 0, err
	}
	return w.Balance, nil
}

func (r *walletRepo) UpdateBalance(ctx context.Context, id string, delta float64) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var balance float64
	err = tx.GetContext(ctx, &balance, "SELECT balance FROM wallets WHERE id=$1 FOR UPDATE", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// создаём кошелёк если не существует
			balance = 0
			_, err = tx.ExecContext(ctx, "INSERT INTO wallets (id, balance) VALUES ($1, $2)", id, delta)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	} else {
		newBalance := balance + delta
		if newBalance < 0 {
			return errors.New("insufficient funds")
		}
		_, err = tx.ExecContext(ctx, "UPDATE wallets SET balance=$1 WHERE id=$2", newBalance, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
