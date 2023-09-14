package storage

import (
	"os"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Структура для теста DB
type DBRepoTest struct {
	suite.Suite
	repo Transaction
}

// Установка параметров
func (s *DBRepoTest) SetupSuite() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		dsn = "postgres://postgres:180612@localhost:5432/urldb?sslmode=disable"
	}
	repo, err := NewDB(dsn)
	require.NoError(s.T(), err)
	s.repo = repo
}

// TestInvoice
func (s *DBRepoTest) TestInvoice() {
	testMoney := Money{
		Currency: "RUB",
		Amount:   100.00,
		CardID:   1234111,
	}
	testMoney2 := Money{
		Currency: "USD",
		Amount:   200.00,
		CardID:   1234116,
	}
	var ReqMoney InvoiceJob
	ReqMoney.RequestMoney = append(ReqMoney.RequestMoney, testMoney, testMoney2)
	err := s.repo.Invoice(ReqMoney)
	require.NoError(s.T(), err)
	response, err := s.repo.GetBalance()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), testMoney.Amount, response[0].Amount) // 100 на 100
	// assert.Equal(s.T(), testMoney.CardID, response[0].CardID) При повторной проверке, т.к. баланс вырастет
}

func (s *DBRepoTest) TestWithdraw() {
	testMoney := Money{
		Currency: "RUB",
		Amount:   99.00,
		CardID:   1234111,
	}
	testMoney2 := Money{
		Currency: "USD",
		Amount:   198.00,
		CardID:   1234116,
	}
	var ReqMoney InvoiceJob
	ReqMoney.RequestMoney = append(ReqMoney.RequestMoney, testMoney, testMoney2)
	err := s.repo.Withdraw(ReqMoney)
	require.NoError(s.T(), err)
	response, err := s.repo.GetBalance()
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1.00, response[0].Amount) // Одноразовая проверка
}

// Запуск теста
func TestDBRepoTest(t *testing.T) {
	suite.Run(t, new(DBRepoTest))
}
