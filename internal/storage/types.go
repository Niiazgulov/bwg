package storage

import "database/sql"

// Основная структура для хранения информации о транзакции
type Money struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount,string"`
	CardID   int     `json:"card_id"`
	Status   string  `json:"status"`
}

// Основной интерфейс проекта, в котором описаны методы работы с хранилищем.
type Transaction interface {
	Invoice(m InvoiceJob) error
	Withdraw(m InvoiceJob) error
	GetBalance() ([]Money, error)
	Close()
}

// Структура БД
type DataBase struct {
	DB *sql.DB
}

// Структура для хранения слайса Money
type InvoiceJob struct {
	RequestMoney []Money
}
