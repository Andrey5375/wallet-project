package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"wallet-service/internal/service"
)

type Handler struct {
	svc *service.WalletService
}

func NewHandler(svc *service.WalletService) *Handler {
	return &Handler{svc: svc}
}

type OperationRequest struct {
	WalletId      string  `json:"walletId"`
	OperationType string  `json:"operationType"`
	Amount        float64 `json:"amount"`
}

func (h *Handler) HandleOperation(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req OperationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	var err error
	switch req.OperationType {
	case "DEPOSIT":
		err = h.svc.Deposit(r.Context(), req.WalletId, req.Amount)
	case "WITHDRAW":
		err = h.svc.Withdraw(r.Context(), req.WalletId, req.Amount)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_ = json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/v1/wallets/")
	balance, err := h.svc.GetBalance(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"balance": balance})
}
