package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/Niiazgulov/bwg.git/internal/storage"
)

// Обработчик входящих транзакций. На вход подается список объектов JSON
func PostInvoiceHandler(repo storage.Transaction, jobCh chan storage.InvoiceJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "PostInvoiceHandler: can't read r.Body", http.StatusBadRequest)
			return
		}
		requestMoney := make([]storage.Money, 0)
		err = json.Unmarshal(request, &requestMoney)
		if err != nil {
			http.Error(w, "PostInvoiceHandler: can't Unmarshal request", http.StatusBadRequest)
			return
		}
		for _, v := range requestMoney {
			v.Status = "Created"
		}
		jobCh <- storage.InvoiceJob{RequestMoney: requestMoney}
		w.WriteHeader(http.StatusCreated)
	}
}

// Обработчик исходящих транзакций. На вход подается список объектов JSON
func PostWithdrawHandler(repo storage.Transaction, jobCh2 chan storage.InvoiceJob) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		request, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "PostWithdrawHandler: can't read r.Body", http.StatusBadRequest)
			return
		}
		reqMoney := make([]storage.Money, 0)
		err = json.Unmarshal(request, &reqMoney)
		if err != nil {
			http.Error(w, "PostWithdrawHandler: can't Unmarshal request", http.StatusBadRequest)
			return
		}
		for _, v := range reqMoney {
			v.Status = "Created"
			fmt.Println(v.Amount, v.CardID)
		}
		jobCh2 <- storage.InvoiceJob{RequestMoney: reqMoney}
		w.WriteHeader(http.StatusCreated)
	}
}

// Обработчик для проверки баланса
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
