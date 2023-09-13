package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Niiazgulov/bwg.git/internal/storage"
)

func PostInvoiceHandler(repo storage.Transaction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputMoney storage.Money
		if err := json.NewDecoder(r.Body).Decode(&inputMoney); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err := repo.Invoice(inputMoney)
		if err != nil {
			http.Error(w, "PostInvoiceHandler: Status internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func PostWithdrawHandler(repo storage.Transaction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var inputMoney storage.Money
		if err := json.NewDecoder(r.Body).Decode(&inputMoney); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		err := repo.Withdraw(inputMoney)
		if err != nil {
			http.Error(w, "PostWithdrawHandler: Status internal server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	}
}

func GetBalanceHandler(repo storage.Transaction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		outputMoney, err := repo.GetBalance()
		if err != nil {
			log.Printf("GetBalanceHandler: unable to get balance from repo: %v", err)
			http.Error(w, "GetBalanceHandler: unable to get balance", http.StatusInternalServerError)
			return
		}
		response, err := json.Marshal(outputMoney)
		if err != nil {
			log.Println("GetBalanceHandler: Error while serializing response", err)
			http.Error(w, "Status internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
