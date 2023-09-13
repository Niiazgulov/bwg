package storage

import "database/sql"

// Основная структура для хранения информации о транзакции
type Money struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount,string"`
	CardID   int     `json:"card_id"`
}

// Основной интерфейс проекта, в котором описаны методы работы с хранилищем.
type Transaction interface {
	Invoice(m Money) error
	Withdraw(m Money) error
	GetBalance() ([]Money, error)
	Close()
}

type DataBase struct {
	DB *sql.DB
}
