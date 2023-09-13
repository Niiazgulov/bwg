package storage

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
)

// Функция для создания нового объекта структуры DataBase.
func NewDB(dbPath string) (Transaction, error) {
	db, err := sql.Open("pgx", dbPath)
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS transsys (
			id SERIAL PRIMARY KEY,
			currency VARCHAR, 
			amount DECIMAL,
			card_id INTEGER UNIQUE)
		`)
	if err != nil {
		return nil, fmt.Errorf("unable to CREATE TABLE in DB: %w", err)
	}
	_, err = db.Exec(`ALTER TABLE transsys DROP CONSTRAINT IF EXISTS amount_nonnegative`)
	if err != nil {
		return nil, fmt.Errorf("unable to DROP CONSTRAINT in DB: %w", err)
	}
	_, err = db.Exec(`ALTER TABLE transsys ADD CONSTRAINT amount_nonnegative CHECK (amount >= 0)`)
	if err != nil {
		return nil, fmt.Errorf("unable to ADD amount_nonnegative CHECK in DB: %w", err)
	}
	return &DataBase{DB: db}, nil
}

// Метод Invoice для зачисления средств на карту (1 карта = 1 валюта)
func (d *DataBase) Invoice(m Money) error {
	query := `INSERT INTO transsys (currency, amount, card_id) VALUES ($1, $2, $3)`
	_, err := d.DB.Exec(query, m.Currency, m.Amount, m.CardID)
	if err != nil && strings.Contains(err.Error(), pgerrcode.UniqueViolation) {
		query = `UPDATE transsys SET amount = amount+$1 WHERE currency = $2 AND card_id = $3`
		_, err := d.DB.Exec(query, m.Amount, m.Currency, m.CardID)
		if err != nil {
			return fmt.Errorf("[Invoice DB] Error while UPDATE transsys: %w", err)
		}
	}
	return nil
}

// Метод Withdraw для списани средств с карты
func (d *DataBase) Withdraw(m Money) error {
	query := `UPDATE transsys SET amount = amount-$1 WHERE currency = $2 AND card_id = $3`
	_, err := d.DB.Exec(query, m.Amount, m.Currency, m.CardID)
	if err != nil {
		return ErrTransactionFailed
	}

	return nil
}

// Метод для извлечения из хранилища информации о балансе.
func (d *DataBase) GetBalance() ([]Money, error) {
	var result []Money
	rows, err := d.DB.Query(`SELECT currency, amount, card_id FROM transsys WHERE EXISTS (SELECT 1 FROM transsys WHERE card_id > 0)`)
	if err != nil {
		return []Money{}, err
	}
	for rows.Next() {
		s := Money{}
		if err := rows.Scan(&s.Currency, &s.Amount, &s.CardID); err != nil {
			return []Money{}, err
		}
		fmt.Println(s.Amount, s.CardID, s.Currency)
		result = append(result, s)
	}
	return result, nil
}

// Метод для закрытия БД.
func (d DataBase) Close() {
	d.DB.Close()
}
